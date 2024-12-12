package src

import (
	"bufio"
	"fmt"
	"log"
	"net"
	shared "shared_mods"
	"sync"
)

var mu sync.Mutex

type Status string

// These variables are going to be send to the message queue
const (
	Free       Status = "free\n"
	Busy       Status = "busy\n"
	InvalidURL Status = "inurl\n"
)

var ServerStatus Status = Free

func InitServer() {
	socket := shared.DownloaderServerSocket

	ln, err := net.Listen("tcp", socket)
	if err != nil {
		log.Fatalln("Error when initializing the server.", err)
		return
	}
	defer ln.Close()

	log.Println("Server listening on", socket)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error when connecting:", err)
			continue
		}

		// Handling connection
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	conn.Write([]byte(ServerStatus))

	defer conn.Close()

	// Read the data from the client
	reader := bufio.NewReader(conn)

	for {
		if ServerStatus != Busy && ServerStatus == InvalidURL {
			ServerStatus = Free
		}

		switch ServerStatus {
		case Free:
			log.Println("Server's free, waiting for the message queue.")
		case InvalidURL:
			log.Println("Server's free, input a valid URL.")
		default:
			log.Println("Server's working, wait for it.")
		}

		data, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error when reading the data:", err)
			return
		}
		// Commands
		if data == "serverStatus\n" {
			conn.Write([]byte(ServerStatus))
			continue
		}

		if ServerStatus == Free { // Processing the data when server isn't busy.
			url, err := URLFilter(data) // Filters invalid URLs
			if err != nil {             // If url is invalid
				mu.Lock()
				ServerStatus = InvalidURL
				mu.Unlock()
				conn.Write([]byte(ServerStatus))
				continue
			}
			conn.Write([]byte(ServerStatus))
			StartDownload(url)
		} else {
			conn.Write([]byte(ServerStatus))
		}
	}
}
