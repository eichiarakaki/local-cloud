package main

import "downloader-server/src"

func main() {
	socket := "localhost:3002"
	src.InitServer(socket)
}
