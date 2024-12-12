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
	responseChannel := make(chan string, 1)
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

	go func() {
		for {
			fmt.Fprintf(conn, "%s\n", "serverStatus") // Requesting server's status

			mainQueue.PrintQueue()
			response := <-responseChannel
			fmt.Println(response)

			if response != "busy\n" {
				url := mainQueue.Dequeue()
				if url != "" {
					fmt.Fprintf(conn, "%s\n", url) // Sends the next URL to the downloader server
				}
			}
			time.Sleep(time.Millisecond * 200)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	var url string
	for {
		// Checks if there's input from the user
		if scanner.Scan() {
			url = scanner.Text()
			mainQueue.Enqueue(url)
			fmt.Printf("%s Enqueued\n", url)
		} else {
			break
		}
	}
}
