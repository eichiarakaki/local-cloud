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
	defer func(ln net.Listener) {
		err := ln.Close()
		if err != nil {
			log.Fatalln("[ERROR] When closing the listener.", err)
		}
	}(ln)

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
	_, err := conn.Write([]byte(ServerStatus))
	if err != nil {
		log.Println("[ERROR] When writing to connection:", err)
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("[ERROR] When closing the connection:", err)
		}
	}(conn)

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
		} else {
			log.Printf("[INFO] Commmand: %s. URL: %s\n", dataDecom.Command, dataDecom.URL)
		}

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
				_, err := conn.Write([]byte(ServerStatus))
				if err != nil {
					log.Println("[ERROR] When sending current Server Status:", err)
					return
				}

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
