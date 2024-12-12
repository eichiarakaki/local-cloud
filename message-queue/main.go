package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	shared "shared_mods"
)

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
	ch := make(chan string, 1)
	go func() {
		for {
			// Reads response
			response, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				log.Println("Error when reading from the server:", err)
				return
			}
			ch <- response
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		// Checks if there's input from the user
		if scanner.Scan() {
			url := scanner.Text()
			fmt.Fprintf(conn, "%s\n", url)
		} else {
			break
		}

		response, ok := <-ch
		if !ok {
			log.Println("Server connection closed. Exiting.")
			break
		}
		ResponseHandler(response)
	}
}
