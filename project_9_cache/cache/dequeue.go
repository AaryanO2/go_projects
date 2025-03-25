package cache

type Node struct {
	Left  *Node
	Right *Node
	Entry string
}

type DeQueue struct {
	Head *Node
	Tail *Node
}

func NewQueue() DeQueue {
	head := &Node{}
	tail := &Node{}
	head.Right = tail
	tail.Left = head
	return DeQueue{Head: head, Tail: tail}
}
