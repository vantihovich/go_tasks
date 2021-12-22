package model

import (
	"fmt"
)

// BTree implementation (very basic with just insertion logic)
type BTree struct {
	Root *Node

	degree int
}

// Node of the BTree
type Node struct {
	leaf   bool
	keys   []int
	childs []*Node

	n int
}

// allocateNode returns a new node with allocated slices of keys and kids
func allocateNode(degree int) *Node {
	return &Node{
		keys:   make([]int, 2*degree-1),
		childs: make([]*Node, 2*degree),
	}
}

// NewBTree returns a BTree with min degree passed
// It also creates a root of the tree
func NewBTree(degree int) *BTree {
	x := allocateNode(degree)
	x.leaf = true

	return &BTree{Root: x, degree: degree}
}

// Insert key to the BTree
func (t *BTree) Insert(k int) {
	fmt.Printf("\n------trying to insert value `%d`------\n", k)

	// if root node is not full - simply insert value into it
	if t.Root.n != 2*t.degree-1 {
		t.insertNonFull(t.Root, k)
		return
	}

	// b-tree grows from the bottom to the top.
	// here we allocate a new root and make existing root it's child
	s := allocateNode(t.degree)
	s.leaf = false
	s.childs[0] = t.Root

	// split root node
	t.splitChild(s, 0)

	i := 0
	if s.keys[0] < k {
		i++
	}

	t.insertNonFull(s.childs[i], k)

	t.Root = s
}

// Traverse displays tree (very naive approach) in CLI
func (t *BTree) Traverse() {
	if t.Root == nil {
		fmt.Printf("\nTree is empty\n")
		return
	}

	traverse(t.Root)
}

func traverse(n *Node) {
	if n == nil {
		fmt.Printf(" ---[empty]--- ")
		return
	}

	fmt.Printf(" --- Node %+v --- ", n.keys[0:n.n])
	if !n.leaf {
		fmt.Println("")
		for _, nn := range n.childs {
			traverse(nn)
		}
	}
}

func (t *BTree) insertNonFull(x *Node, k int) {
	fmt.Printf("\ncalled t.insertNonFull for node %+v and value %d\n", x, k)

	// initialise i with right element
	i := x.n - 1

	if x.leaf {
		// find the place for the new key
		for i >= 0 && x.keys[i] > k {
			x.keys[i+1] = x.keys[i]
			i--
		}

		x.keys[i+1] = k
		x.n = x.n + 1
		return
	}

	// find position for the new node
	for i >= 0 && x.keys[i] > k {
		i = i - 1
	}

	if x.childs[i+1].n == 2*t.degree-1 {
		// when child not is full
		t.splitChild(x, i+1)
		// after splitting key in the middle of child node goes up
		// child node gets splitted in two
		if x.keys[i+1] < k {
			i++
		}
	}
	t.insertNonFull(x.childs[i+1], k)
}

func (t *BTree) splitChild(x *Node, i int) {
	fmt.Printf("\n[!] called t.SplitChild for node %+v at index %d\n", x, i)
	y := x.childs[i]

	// allocated a new node
	z := allocateNode(t.degree)
	z.leaf = y.leaf
	z.n = t.degree - 1

	// copy Y keys to Z
	for j := 0; j < t.degree-1; j++ {
		z.keys[j] = y.keys[j+t.degree]
	}

	if !y.leaf {
		for j := 0; j < t.degree; j++ {
			z.childs[j] = y.childs[j+t.degree]
		}
	}
	y.n = t.degree - 1

	// insert child node to the child node
	for j := x.n; j >= i+1; j-- {
		x.childs[j+1] = x.childs[j]
	}

	x.childs[i+1] = z

	// move keys in the node
	for j := x.n - 1; j >= i; j-- {
		x.keys[j+1] = x.keys[j]
	}
	x.keys[i] = y.keys[t.degree-1]

	x.n = x.n + 1

	fmt.Printf("\nALLOCATED NODE X %+v\n", x)
	fmt.Printf("\nALLOCATED NODE Z %+v\n", z)
	fmt.Printf("\nALLOCATED NODE Y %+v\n", y)
}
