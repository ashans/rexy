package config

import (
	"io/ioutil"
	"log"
)
import "gopkg.in/yaml.v2"

type Config struct {
	Server    ServerConfig     `yaml:"server"`
	Endpoints []EndpointConfig `yaml:"endpoints"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type EndpointConfig struct {
	Context string                `yaml:"context"`
	Forward EndpointForwardConfig `yaml:"forward"`
}

type EndpointForwardConfig struct {
	Protocol string `yaml:"protocol"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Rewrite  bool   `yaml:"rewrite"`
	Context  string `yaml:"context"`
}

func LoadConfigFromFile(fileName string) *Config {
	config := &Config{}

	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Error loading config file : %s, cause :%s", fileName, err.Error())
	}

	err = yaml.Unmarshal(file, config)
	if err != nil {
		log.Fatalf("Error parsing config file : %s, cause :%s", fileName, err.Error())
	}

	return config
}
