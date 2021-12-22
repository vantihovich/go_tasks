package main

import (
	"fmt"

	"github.com/apaliavy/godel-golang/data-structures/trees/rb/model"
)

func main() {
	t := model.NewRBTree()

	values := []int{1, 2, 3, 4, 5, 6}

	for _, v := range values {
		t.Insert(v)
	}

	printInOrder(t.Root)
}

func printInOrder(n *model.Node) {
	if n == nil {
		return
	}

	printInOrder(n.Left)
	fmt.Printf("%s\n", n)
	printInOrder(n.Right)
}
