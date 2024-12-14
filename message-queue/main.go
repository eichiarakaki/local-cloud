package main

import (
	"bufio"
	"fmt"
	"log"
	"message-queue/queue"
	"net"
	"os"
	shared "shared_mods"
	"time"
)

// Global Variable
var mainQueue = queue.NewQueue()

func main() {
	// Loading config
	err := shared.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}
	// Getting socket from the loaded config
	socket := shared.DownloaderServerSocket

	// Initializing connection to the socket.
	conn, err := net.Dial("tcp", socket)
	if err != nil {
		log.Println("Couldn't connect to the downloader server:", err)
	}
	defer conn.Close()

	// Making a goroutine for handling infinite responses.
	responseChannel := make(chan string)
	go func() {
		for {
			// Reads response
			response, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				log.Println("Error when reading from the server:", err)
				return
			}
			responseChannel <- response
		}
	}()

	// A special goroutine to know the Downlaoder Server status
	var response string
	go func() {
		for {
			response = <-responseChannel
			log.Println("INFO: New response from the Downloader Server:", response)
		}
	}()

	// Handling the response from the Downloader Server
	go func() {
		for {
			if response == "free\n" && !mainQueue.IsEmpty() { // If Downloader Server is free and the queue isn't empty
				url := mainQueue.Dequeue()
				if url != "" {
					fmt.Fprintf(conn, "lock %s\n", url) // Sends the next URL to the downloader server
				}

				log.Println("INFO: Sent to the downloader server.")
			}

			if response == "busy\n" {
				log.Println("INFO: Downloader Server is busy.")
			}

			time.Sleep(time.Second * 2)
		}

	}()

	scanner := bufio.NewScanner(os.Stdin)
	var url string
	for {
		// Checks if there's input from the user
		if scanner.Scan() {
			url = scanner.Text()
			mainQueue.Enqueue(url)
			mainQueue.PrintQueue()
		} else {
			break
		}
	}
}
