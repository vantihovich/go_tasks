package model

import "fmt"

type Node struct {
	Key    int
	Height int
	Left   *Node
	Right  *Node
}

func NewNode(key int) *Node {
	return &Node{
		Key:    key,
		Height: 1,
	}
}

func (n *Node) String() string {
	if n == nil {
		return "nil"
	}

	return fmt.Sprintf("%d", n.Key)
}

func (n *Node) fixHeight() {
	fmt.Printf("\nfixing height for node [%s]\n", n)

	hl := getNodeHeight(n.Left)
	hr := getNodeHeight(n.Right)

	if hl > hr {
		n.Height = hl + 1
	} else {
		n.Height = hr + 1
	}

	fmt.Printf("\n  node height is %d now\n", n.Height)
}

func insertNode(n *Node, key int) *Node {
	if n == nil {
		return NewNode(key)
	}

	if n.Key > key {
		n.Left = insertNode(n.Left, key)
	} else if n.Key < key {
		n.Right = insertNode(n.Right, key)
	}

	n.Height = 1 + getNodeHeight(n.Left) + getNodeHeight(n.Right)

	return balance(n)
}

func removeNode(n *Node, key int) *Node {
	if key < n.Key {
		n.Left = removeNode(n.Left, key)
		return balance(n)
	}

	if key > n.Key {
		n.Right = removeNode(n.Right, key)
		return balance(n)
	}

	q, r := n.Left, n.Right
	if r == nil {
		return q
	}

	min := findMin(r)
	min.Right = removeMin(r)
	min.Left = q

	return balance(min)
}

func findMin(n *Node) *Node {
	if n.Left != nil {
		return findMin(n.Left)
	}

	return n
}

func removeMin(root *Node) *Node {
	if root.Left == nil {
		return root.Right
	}

	root.Left = removeMin(root.Left)
	return balance(root)
}
