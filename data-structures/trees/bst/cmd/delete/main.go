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
	t.Insert(12)
	t.Insert(3)
	t.Insert(18)
	t.Insert(11)
	t.Insert(13)
	t.Insert(14)

	fmt.Println("PREFIX_TRAVERSE example:")
	printPreOrder(t.Root)
	fmt.Println("======")

	fmt.Println("An attempt to delete 15:")
	t.Remove(15)
	fmt.Println("======")

	fmt.Println("PREFIX_TRAVERSE example after removal:")
	printPreOrder(t.Root)
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
