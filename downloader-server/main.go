package main

import (
	"downloader-server/src"
	"log"
)

func main() {
	err := src.LoadConfig() // Gets config.json properties.
	if err != nil {
		log.Fatalf("Error when loading the config.json: %s\n", err)
		return
	}
	socket := src.DownloaderServerSocket
	src.InitServer(socket)
}
