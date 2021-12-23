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

	fmt.Println("PREFIX_TRAVERSE example:")
	printPreOrder(t.Root)
	fmt.Println("======")

	fmt.Println("An attempt to search 10:")
	t.Find(10)
	fmt.Println("======")

	fmt.Println("An attempt to search 3:")
	t.Find(3)
	fmt.Println("======")

	fmt.Println("An attempt to search 16:")
	t.Find(16)
	fmt.Println("======")

	fmt.Println("An attempt to search 2:")
	t.Find(2)
	fmt.Println("======")

	fmt.Println("An attempt to search 13:")
	t.Find(13)
	fmt.Println("======")

	fmt.Println("An attempt to search 1:")
	t.Find(1)
	fmt.Println("======")

	fmt.Println("An attempt to search 20:")
	t.Find(20)
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
