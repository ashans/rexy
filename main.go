package main

import (
	"fmt"
	"log"
	"net/http"
	"rexy/config"
	"rexy/core"
)

const ConfigFilePath = "config.yaml"

func main() {

	fmt.Println(`
         __________ _______________  ________.___.         
         \______   \\_   _____/\   \/  /\__  |   |         
  ______  |       _/ |    __)_  \     /  /   |   |  ______ 
 /_____/  |    |   \ |        \ /     \  \____   | /_____/ 
          |____|_  //_______  //___/\  \ / ______|         
                 \/         \/       \_/ \/`)

	log.Printf("Loading config file : %s", ConfigFilePath)
	c := config.LoadConfigFromFile(ConfigFilePath)
	log.Printf("Loaded config file completed")

	handler := core.NewHandler(c).Handler

	http.HandleFunc("/", handler)
	log.Printf("Initiating server on port : %d", c.Server.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", c.Server.Port), nil)
	if err != nil {
		log.Fatalf("Error initializing http server, cause :%s", err.Error())
	}
	log.Printf("Server initialized on port : %d", c.Server.Port)
}
