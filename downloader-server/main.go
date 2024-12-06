package main

import (
	"downloader-server/src"
	"log"
	shared "shared_mods"
)

func main() {
	err := shared.LoadConfig() // Gets config.json properties.
	if err != nil {
		log.Fatalf("Error when loading the config.json: %s\n", err)
		return
	}

	socket := shared.DownloaderServerSocket
	src.InitServer(socket)
}
