package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
)

func main() {
	if _, err := os.Stat("config.json"); err != nil {
		handleErrors("Unable to find config file", err)
	}
	configFile, err := os.Open("config.json")
	handleErrors("Unable to open config file", err)

	defer configFile.Close()

	var apps Config
	decoder := json.NewDecoder(configFile)
	if err := decoder.Decode(&apps); err != nil {
		handleErrors("Unable to decode JSON", err)
	}

	for _, app := range apps.Apps {
		connChannel := make(chan net.Conn)
		for _, port := range app.Ports {
			server := Server{Port: port, Targets: app.Targets}
			log.Println("Listening on port ", port)
			go server.Listen(connChannel)
		}

		go proxy(connChannel, app.Targets)

	}

	select {}

}
