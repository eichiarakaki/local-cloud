package queue

type Node struct {
	NextNode *Node
	value    string // URL
}

type QueueList struct {
	List []*Node
}

func (q *QueueList) Enqueue(URL string) {

}

func (q *QueueList) Dequeue(URL string) string {

	return ""
}
