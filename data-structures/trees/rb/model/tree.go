package model

import "fmt"

const (
	RED   = 0
	BLACK = 1
)

type Node struct {
	Left   *Node
	Right  *Node
	Parent *Node
	Color  uint

	Value int
}

func (n *Node) String() string {
	if n.Color == RED {
		return fmt.Sprintf("Node [%d, RED]", n.Value)
	}
	return fmt.Sprintf("Node [%d, BLACK]", n.Value)
}

type RBTree struct {
	Root *Node
	NIL  *Node
}

func NewRBTree() *RBTree {
	root := &Node{Color: BLACK}
	return &RBTree{
		Root: root,
		NIL:  root,
	}
}

func (t *RBTree) Insert(val int) *Node {
	z := &Node{
		Left:  t.NIL,
		Right: t.NIL,
		Color: RED,
		Value: val,
	}

	x := t.Root
	y := t.NIL

	for x != t.NIL {
		y = x

		if z.Value < x.Value {
			x = x.Left
		} else {
			x = x.Right
		}
	}

	z.Parent = y
	if y == t.NIL {
		t.Root = z
	} else if z.Value < y.Value {
		y.Left = z
	} else {
		y.Right = z
	}

	t.insertFixup(z)
	return z
}

func (t *RBTree) insertFixup(z *Node) {
	for z.Parent.Color == RED {
		if z.Parent == z.Parent.Parent.Left {
			y := z.Parent.Parent.Right
			if y.Color == RED {
				// Case 1:
				// Parent and uncle are both RED, the grandparent must be BLACK
				z.Parent.Color = BLACK
				y.Color = BLACK
				z.Parent.Parent.Color = RED
				z = z.Parent.Parent
			} else {
				if z == z.Parent.Right {
					// Case 2:
					// Parent is RED and uncle is BLACK and the current node is right child
					z = z.Parent
					t.rotLeft(z)
				}
				// Case 3:
				z.Parent.Color = BLACK
				z.Parent.Parent.Color = RED
				t.rotRight(z.Parent.Parent)
			}
		} else {
			// same as then clause with "right" and "left" exchanged
			y := z.Parent.Parent.Left
			if y.Color == RED {
				z.Parent.Color = BLACK
				y.Color = BLACK
				z.Parent.Parent.Color = RED
				z = z.Parent.Parent
			} else {
				if z == z.Parent.Left {
					z = z.Parent
					t.rotRight(z)
				}
				z.Parent.Color = BLACK
				z.Parent.Parent.Color = RED
				t.rotLeft(z.Parent.Parent)
			}
		}
	}
	t.Root.Color = BLACK
}

// rotLeft makes a left tree rotation over X node
//          |                                  |
//          X                                  Y
//         / \         left rotate            / \
//        α  Y       ------------->         X   γ
//           / \                            / \
//          β  γ                         α  β
//
func (t *RBTree) rotLeft(x *Node) {
	if x.Right == t.NIL {
		return
	}

	y := x.Right
	x.Right = y.Left
	if y.Left != t.NIL {
		y.Left.Parent = x
	}
	y.Parent = x.Parent

	if x.Parent == t.NIL {
		t.Root = y
	} else if x == x.Parent.Left {
		x.Parent.Left = y
	} else {
		x.Parent.Right = y
	}

	y.Left = x
	x.Parent = y
}

// rotRight makes a right tree rotation over X node
//          |                                  |
//          X                                  Y
//         / \         right rotate           / \
//        Y   γ      ------------->         α  X
//       / \                                    / \
//      α  β                                 β  γ
//
func (t *RBTree) rotRight(x *Node) {
	if x.Left == t.NIL {
		return
	}

	y := x.Left
	x.Left = y.Right
	if y.Right != t.NIL {
		y.Right.Parent = x
	}
	y.Parent = x.Parent

	if x.Parent == t.NIL {
		t.Root = y
	} else if x == x.Parent.Left {
		x.Parent.Left = y
	} else {
		x.Parent.Right = y
	}

	y.Right = x
	x.Parent = y
}
