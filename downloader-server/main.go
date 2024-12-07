package main

import (
	"downloader-server/src"
	"log"
	shared "shared_mods"
)

func main() {
	err := shared.LoadConfig() // Gets config.json properties.
	if err != nil {
		log.Fatalf("%s\n", err)
		return
	}

	src.InitServer()
}
