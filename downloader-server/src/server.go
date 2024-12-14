package src

import (
	"bufio"
	"fmt"
	"log"
	"net"
	shared "shared_mods"
	"strings"
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

	// Sending cycle: Listening for inputs -> Server Status (Free) -> Busy -> (Finishes downlading) Free -> To the beginning.
	for {
		data, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error when reading the data:", err)
			return
		}

		dataDecom, err := RDataWrapper(data)
		if err != nil {
			log.Printf("ERROR: %s\n", err)
		}
		log.Printf("INFO:\nCommmand: %s\nURL: %s\n", dataDecom.Command, dataDecom.URL)

		// Commands
		if strings.HasPrefix(dataDecom.Command, "test") {
			_, err := URLFilter(dataDecom.URL)
			if err != nil {
				conn.Write([]byte("test inurl"))
			} else {
				conn.Write([]byte("test vaurl"))
			}
			continue
		} else if strings.HasPrefix(data, "lock") && ServerStatus == Free {
			mu.Lock()
			ServerStatus = Busy
			mu.Unlock()

			log.Printf("INFO: Got a lock signal.\n")

			url, err := URLFilter(dataDecom.URL) // Filters invalid URLs
			// If url is invalid
			if err != nil {
				mu.Lock()
				ServerStatus = InvalidURL
				mu.Unlock()

				log.Println("INFO: URL Filtered.")
				conn.Write([]byte(ServerStatus))

				mu.Lock()
				ServerStatus = Free
				mu.Unlock()

				continue
			}

			fmt.Printf("Start downloading.\n")
			StartDownload(url, conn)
		}
	}
}
