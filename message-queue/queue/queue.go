package queue

import (
	"fmt"
	"strings"
)

/*
A Queue Data Structure is a fundamental concept in computer science used for
storing and managing data in a specific order. It follows the principle
of "First in, First out" (FIFO), where the first element added to the queue
is the first one to be removed. Queues are commonly used in various algorithms and
applications for their simplicity and efficiency in managing data flow.
*/
type Node struct {
	NN    *Node  // Next Node
	Value string // URL
}

type QueueList struct {
	List *Node
}

func NewQueue() *QueueList {
	return &QueueList{}
}

func (q *QueueList) Enqueue(URL string) {
  // Do nothing if the url is empty
  if strings.TrimSpace(URL) == "" {
    return
  }

	newNode := &Node{Value: URL}
	if q.List == nil {
		q.List = newNode
		return
	}

	auxNode := q.List
	for auxNode.NN != nil {
		auxNode = auxNode.NN
	}
	auxNode.NN = newNode
}

// Returns the 'First Out' from the queue and actualizes it.
func (q *QueueList) Dequeue() string {
	if q.List != nil {
		firstNodeValue := q.List.Value // if q.List is NOT nil, it means q.List.Value has a value.
		q.List = q.List.NN             // Links the List to the NEXT node

		return firstNodeValue
	}
	return ""
}

func (q *QueueList) PrintQueue() {
	auxNode := q.List
	for auxNode != nil {
		fmt.Printf("VALUE: %s, Next Node: %p\n", auxNode.Value, auxNode.NN)
		auxNode = auxNode.NN
	}
}

func (q *QueueList) IsEmpty() bool {
	return q.List == nil
}

// func (q *QueueList) SafeRequest(conn net.Conn) {
// 	// Getting the value to send to the server
// 	var firstNodeValue string
// 	if q.List != nil {
// 		firstNodeValue = q.List.Value // if q.List is NOT nil, it means q.List.Value has a value.
// 	}

// 	// Sending the value to the server
// 	toSend := fmt.Sprintf("test %s", firstNodeValue)
// 	conn.Write([]byte(toSend))
// 	// Reading the server's response
// 	response, err := bufio.NewReader(conn).ReadString('\n')
// 	if err != nil {
// 		log.Println("Error when reading from the server:", err)
// 		return
// 	}

// }
