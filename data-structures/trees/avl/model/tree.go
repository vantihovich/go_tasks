package model

import "fmt"

type AVLTree struct {
	Root *Node
}

func (t *AVLTree) Insert(key int) {
	fmt.Printf("\n------ inserting [%d] ------\n", key)

	t.Root = insertNode(t.Root, key)

	fmt.Printf("\nnew root is [%s]\n", t.Root)
}

func (t *AVLTree) Remove(key int) {
	if t.Root == nil {
		return
	}

	t.Root = removeNode(t.Root, key)
}

func balance(n *Node) *Node {
	fmt.Printf("\nbalancing node %d\n", n.Key)

	n.fixHeight()

	nbf := balanceFactor(n)

	fmt.Printf("\nbalance factor for node [%s] is %d\n", n, nbf)

	// left rotation
	if nbf == 2 {
		if balanceFactor(n.Right) < 0 { // big left rotation
			fmt.Printf("\n[!] doing big left rotation\n")
			n.Right = rotRight(n.Right)
		} else {
			fmt.Printf("\n[!] doing small left rotation\n")
		}
		return rotLeft(n)
	}

	if nbf == -2 {
		// right rotation
		if balanceFactor(n.Left) > 0 { // big right rotation
			fmt.Printf("\n[!] doing big right rotation\n")
			n.Left = rotLeft(n.Left)
		} else {
			fmt.Printf("\n[!] doing small right rotation\n")
		}
		return rotRight(n)
	}

	return n
}

func balanceFactor(n *Node) int {
	if n == nil {
		return 0
	}

	return getNodeHeight(n.Right) - getNodeHeight(n.Left)
}

func getNodeHeight(n *Node) int {
	if n == nil {
		return 0
	}

	return n.Height
}

func rotLeft(root *Node) *Node {
	newRoot := root.Right
	root.Right = newRoot.Left
	newRoot.Left = root

	root.fixHeight()
	newRoot.fixHeight()

	return newRoot
}

func rotRight(root *Node) *Node {
	newRoot := root.Left
	root.Left = newRoot.Right
	newRoot.Right = root

	root.fixHeight()
	newRoot.fixHeight()

	return newRoot
}
