package core

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"rexy/config"
	"strings"
)

type Handler struct {
	endpoints map[string]*config.EndpointForwardConfig
}

func NewHandler(c *config.Config) *Handler {
	e := make(map[string]*config.EndpointForwardConfig)
	for _, endpoint := range c.Endpoints {
		e[endpoint.Context] = &endpoint.Forward
	}

	return &Handler{
		endpoints: e,
	}
}

func (h *Handler) Handler(w http.ResponseWriter, req *http.Request) {

	var c *config.EndpointForwardConfig
	var pref string
	for k, v := range h.endpoints {
		if strings.HasPrefix(req.RequestURI, k) {
			c = v
			pref = k
		}
	}

	if c == nil {
		http.Error(w, fmt.Sprintf("REXY | No matching handler for path: %s", req.RequestURI), http.StatusNotFound)
		return
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("REXY | Error reading input body: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	req.Body = ioutil.NopCloser(bytes.NewReader(body))

	// create a new url from the raw RequestURI sent by the client
	var dest string
	if c.Port == 0 {
		dest = c.Host
	} else {
		dest = fmt.Sprintf("%s:%d", c.Host, c.Port)
	}
	var reqContext string
	if c.Rewrite {
		reqContext = c.Context + req.RequestURI[len(pref):]
	} else {
		reqContext = c.Context + req.RequestURI
	}
	url := fmt.Sprintf("%s://%s%s", c.Protocol, dest, reqContext)

	proxyReq, err := http.NewRequest(req.Method, url, bytes.NewReader(body))

	// We may want to filter some headers, otherwise we could just use a shallow copy
	// proxyReq.Header = req.Header
	proxyReq.Header = make(http.Header)
	for h, val := range req.Header {
		proxyReq.Header[h] = val
	}

	for q, all := range req.URL.Query() {
		for _, v := range all {
			proxyReq.URL.Query().Add(q, v)
		}
	}

	httpClient := http.Client{}

	resp, err := httpClient.Do(proxyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			http.Error(w, fmt.Sprintf("REXY | Error writing response body: %s", err.Error()), http.StatusInternalServerError)
		}
	}(resp.Body)

	for h, val := range resp.Header {
		for _, v := range val {
			w.Header().Add(h, v)
		}
	}

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("REXY | Error writing response body: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(resp.StatusCode)
}
