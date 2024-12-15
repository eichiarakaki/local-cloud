package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"message-queue/queue"
	"message-queue/utils"
	"net"
	shared "shared_mods"
	"time"
)

// Global Variable
var mainQueue = queue.NewQueue()
var response string

func main() {
	// Loading config
	err := shared.LoadConfig("../config.json")
	if err != nil {
		log.Fatalln(err)
	}
	// Getting socket from the loaded config
	downloaderServerSocket := shared.DownloaderServerSocket

	// Initializing connection to the Downloader server.
	conn, err := net.Dial("tcp", downloaderServerSocket)
	if err != nil {
		log.Println("[ERROR] Couldn't connect to the downloader server:", err)
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("[ERROR] Couldn't close the connection:", err)
		}
	}(conn)

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

	// A Go-routine to know the Downloader Server status
	go func() {
		for {
			response = <-responseChannel
			log.Println("[INFO] New response from the Downloader Server:", response)
		}
	}()

	// Handling the response from the Downloader Server
	var prevResponse string
	go func() {
		for {
			if response == "free\n" && !mainQueue.IsEmpty() { // If Downloader Server is free and the queue isn't empty
				url := mainQueue.Dequeue()
				if url != "" {
					_, err = fmt.Fprintf(conn, "lock %s\n", url) // Sends the next URL to the downloader server
					if err != nil {
						log.Println("[ERROR]", err)
					}
				}

				log.Println("[INFO] Sent to the downloader server.")
			}

			if response == "busy\n" && prevResponse != "busy\n" {
				log.Println("[INFO] Downloader Server is busy.")
			}

			prevResponse = response
			time.Sleep(time.Second * 2)
		}

	}()

	// Initializing Server to handle requests from the backend
	messageQueueSocket := shared.MessageQueueSocket
	ln, err := net.Listen("tcp", messageQueueSocket)
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

	log.Println("[INFO] Message Queue listening on", messageQueueSocket)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("[ERROR] Couldn't accept the connection.", err)
			continue
		}
		// Handling connection
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {

	// Reading from the TCP
	reader := bufio.NewReader(conn)

	for {
		data, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return
			} else {
				log.Println("[ERROR] Couldn't read from the connection:", err)
				return
			}
		}
		log.Printf("[INFO] Got %s from the backend\n", data)

		// Sending response to the backend
		dataPosition := mainQueue.Position(data)
		msg := fmt.Sprintf("%s was enqueued successfully.\n", data)
		toBackend := utils.MQBackendWrapper(response, dataPosition, msg)

		jsonData, err := json.Marshal(toBackend)
		if err != nil {
			log.Println("[ERROR]", err)
			return
		}

		jsonData = append(jsonData, '\n') // It must have a '\n' at the end because this response will be obtained with a function that reads by lines.
		_, err = conn.Write(jsonData)
		if err != nil {
			log.Println("[ERROR] Couldn't write to the connection:", err)
			return
		}

		mainQueue.Enqueue(data)
		mainQueue.PrintQueue()
	}
}
