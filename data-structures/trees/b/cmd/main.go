package main

import "github.com/apaliavy/godel-golang/data-structures/trees/b/model"

func main() {
	t := model.NewBTree(3)

	t.Insert(1)
	t.Insert(2)
	t.Insert(3)
	t.Insert(4)
	t.Insert(5)
	t.Insert(6)
	t.Insert(7)
	t.Insert(8)
	t.Insert(9)
	t.Insert(10)

	t.Traverse()
}
