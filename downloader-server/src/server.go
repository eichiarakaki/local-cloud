package src

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

var MU sync.Mutex
var ServerStatus bool = false

func InitServer(socket string) {
	ln, err := net.Listen("tcp", socket)
	if err != nil {
		log.Println("Error when initializing the server.", err)
		return
	}
	defer ln.Close()

	log.Println("Server listening on", socket)

	for {
		// Accepts only one connection
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error when connecting:", err)
			continue
		}

		// Handling connection
		go handleConnection(conn)

	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Read the data from the client
	reader := bufio.NewReader(conn)

	for {
		log.Println("Server's free, waiting for input.")

		data, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error when reading the data:", err)
			return
		}
		if !ServerStatus { // Processing the data when server isn't busy.
			url := strings.TrimSpace(data)
			response := fmt.Sprintf("Downloading %s\n", url)
			conn.Write([]byte(response))
			StartDownload(url)
		} else {
			// Response to the client
			response := fmt.Sprintf("Server is busy.\n")
			conn.Write([]byte(response))
		}
	}
}

func StartDownload(url string) {
	MU.Lock()
	ServerStatus = true // Changes status to busy server
	MU.Unlock()

	go func() {
		Download(url)
		MU.Lock()
		ServerStatus = false // Changes Status to free server
		MU.Unlock()
	}()
}
