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

	fmt.Printf("%+v", t)
}
