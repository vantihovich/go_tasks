package main

import (
	"fmt"

	"github.com/apaliavy/godel-golang/data-structures/trees/bst/model"
)

func main() {
	var t model.Tree

	t.Insert(10)
	t.Insert(1)
	t.Insert(15)
	t.Insert(16)
	t.Insert(8)
	t.Insert(10)
	t.Insert(3)

	fmt.Println("INFIX_TRAVERSE example:")
	printInOrder(t.Root)
	fmt.Println("======")

	fmt.Println("PREFIX_TRAVERSE example:")
	printPreOrder(t.Root)
	fmt.Println("======")

	fmt.Println("POSTFIX_TRAVERSE example:")
	printPostOrder(t.Root)
	fmt.Println("======")
}

func printPreOrder(n *model.Node) {
	if n == nil {
		return
	}

	fmt.Printf("%d ", n.Key)
	printPreOrder(n.Left)
	printPreOrder(n.Right)
}

func printPostOrder(n *model.Node) {
	if n == nil {
		return
	}

	printPostOrder(n.Left)
	printPostOrder(n.Right)
	fmt.Printf("%d ", n.Key)
}

func printInOrder(n *model.Node) {
	if n == nil {
		return
	}

	printInOrder(n.Left)
	fmt.Printf("%d ", n.Key)
	printInOrder(n.Right)
}
