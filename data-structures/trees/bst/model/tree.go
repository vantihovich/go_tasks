package model

type Tree struct {
	Root *Node
}

func (t *Tree) Insert(data int) {
	if t.Root == nil {
		t.Root = &Node{Key: data}
		return
	}
	t.Root.Insert(data)
}

type Node struct {
	Key   int
	Left  *Node
	Right *Node
}

func (n *Node) Insert(data int) {
	if data < n.Key {
		// insert into the left tree
		if n.Left == nil {
			n.Left = &Node{Key: data}
		} else {
			n.Left.Insert(data)
		}
		return
	}
	// insert into the right tree
	if n.Right == nil {
		n.Right = &Node{Key: data}
	} else {
		n.Right.Insert(data)
	}
}
