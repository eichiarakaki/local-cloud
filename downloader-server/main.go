package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	socket := "localhost:3002"
	ln, err := net.Listen("tcp", socket)
	if err != nil {
		log.Println("Error when initializing the server.", err)
		return
	}
	defer ln.Close()

	log.Println("Server listening on", socket)

	for {
		// Accept Connections
		conn, err := ln.Accept()
		if err != nil {
      fmt.Println("Error when accepting connection:", err)
      continue
		}
    // Handle new connection
    go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
  defer conn.Close()

  // Read the data from the client
  reader := bufio.NewReader(conn)
  for {
    data, err := reader.ReadString('\n')
    if err != nil {
      fmt.Println("Error when reading the data:",err)
      return
    }
    // Process the received data
    url := strings.TrimSpace(data)
    fmt.Println(url)

    // Response to the client
    response := fmt.Sprintf("URL %s processed\n", url)
    conn.Write([]byte(response))
  }
}
