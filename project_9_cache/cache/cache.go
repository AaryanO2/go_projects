package cache

import "fmt"

type Cache struct {
	Queue    DeQueue
	Hash     map[string]*Node
	Capacity int
}

func NewCache(Capacity int) Cache {
	return Cache{Queue: NewQueue(), Hash: map[string]*Node{}, Capacity: Capacity}
}

func (c *Cache) Check(s string) {
	node, ok := c.Hash[s]
	if !ok {
		if c.Capacity == 0 {
			c.RemoveLast()
		}

		node = &Node{Entry: s}
		c.Capacity--
		c.Hash[s] = node

	} else {
		c.RemoveFromQueue(node)
	}
	c.Add(node)
}

func (c *Cache) Add(node *Node) {
	node.Right = c.Queue.Head.Right
	node.Left = c.Queue.Head
	c.Queue.Head.Right.Left = node
	c.Queue.Head.Right = node
}

func (c *Cache) RemoveFromQueue(node *Node) *Node {
	node.Left.Right = node.Right
	node.Right.Left = node.Left
	return node
}

func (c *Cache) RemoveLast() {
	c.Capacity++
	delete(c.Hash, c.Queue.Tail.Left.Entry)
	c.RemoveFromQueue(c.Queue.Tail.Left)
}

func (c *Cache) Display() {
	node := c.Queue.Head.Right
	for node != c.Queue.Tail {
		fmt.Print(node.Entry)
		node = node.Right
	}
	fmt.Print("\n")
}
