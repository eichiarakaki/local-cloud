package main

import (
	"fmt"
	"strings"
)

type Node struct {
	NN    *Node  // Next Node
	Value string // URL
}

type List struct {
	Queue *Node
}

func NewQueue() *List {
	return &List{}
}

func (q *List) Enqueue(URL string) {
	// Do nothing if the url is empty
	if strings.TrimSpace(URL) == "" {
		return
	}

	newNode := &Node{Value: URL}
	if q.Queue == nil {
		q.Queue = newNode
		return
	}

	auxNode := q.Queue
	for auxNode.NN != nil {
		auxNode = auxNode.NN
	}
	auxNode.NN = newNode
}

// Dequeue Returns the 'First Out' from the queue and actualizes it.
func (q *List) Dequeue() string {
	if q.Queue != nil {
		firstNodeValue := q.Queue.Value // if q.List is NOT nil, it means q.List.Value has a value.
		q.Queue = q.Queue.NN            // Links the List to the NEXT node

		return firstNodeValue
	}
	return ""
}

func (q *List) PrintQueue() {
	auxNode := q.Queue
	for auxNode != nil {
		fmt.Printf("VALUE: %s, Next Node: %p\n", auxNode.Value, auxNode.NN)
		auxNode = auxNode.NN
	}
}

func (q *List) IsEmpty() bool {
	return q.Queue == nil
}

// Position Finds the position inside the queue
func (q *List) Position(url string) uint16 {
	if q.IsEmpty() {
		return 1
	}

	// If the Queue is not empty, searches its position
	auxNode := q.Queue
	var c uint16 = 1
	for auxNode != nil {
		// Trim spaces and newlines before comparison
		if strings.TrimSpace(auxNode.Value) == strings.TrimSpace(url) {
			return c
		}
		c++
		auxNode = auxNode.NN
	}

	// Returns 0 if not found.
	return 0
}

// Length Returns the length of the queue.
func (q *List) Length() uint16 {
	auxNode := q.Queue
	var c uint16 = 0
	for auxNode != nil {
		c++
		auxNode = auxNode.NN
	}

	if Response == "busy\n" {
		c++
	}

	return c
}
