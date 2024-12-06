package main

import (
	"bufio"
	"fmt"
	"net"
  "log"
	"os"
)

func main() {
  socket := "localhost:3002" // MUST MATCH WITH THE DOWNLOADER SERVER
  conn, err := net.Dial("tcp", socket)
  if err != nil {
    log.Println("Couldn't connect to the downloader server:",err)
  }
  defer conn.Close()

  for {
    scanner := bufio.NewScanner(os.Stdin)
    fmt.Printf(">>> ")
    if scanner.Scan() {
      url := scanner.Text()
      fmt.Fprintf(conn,"%s\n",url)

      // Reads response
      response, err := bufio.NewReader(conn).ReadString('\n')
      if err != nil {
        log.Println("Error when reading from the server:", err)
        return
      }

      // Prints the response
      log.Println("Response from the server:", response)
    }
  }
}
