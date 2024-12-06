package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	socket := "localhost:3002" // MUST MATCH WITH THE DOWNLOADER SERVER
	conn, err := net.Dial("tcp", socket)
	if err != nil {
		log.Println("Couldn't connect to the downloader server:", err)
	}
	defer conn.Close()

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
			log.Println("Response from the server:", <-ch)
		}
	}()

	for {
		scanner := bufio.NewScanner(os.Stdin)

		if scanner.Scan() {
			url := scanner.Text()
			fmt.Fprintf(conn, "%s\n", url)
		}

	}
}
