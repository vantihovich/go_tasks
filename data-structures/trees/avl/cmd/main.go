package main

import (
	"fmt"

	"github.com/apaliavy/godel-golang/data-structures/trees/avl/model"
)

func main() {
	var t model.AVLTree

	t.Insert(10)
	t.Insert(9)
	t.Insert(8)
	t.Insert(7)
	t.Insert(6)
	t.Insert(5)

	t.Remove(6)

	infixTraverse(t.Root)
}

func infixTraverse(n *model.Node) {
	if n == nil {
		return
	}

	infixTraverse(n.Left)
	fmt.Printf("%s ", n)
	infixTraverse(n.Right)
}
