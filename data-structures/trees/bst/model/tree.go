package model

import (
	"errors"
	"fmt"
)

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

func (t *Tree) Find(data int) (*Node, error) {
	if t.Root != nil {
		r, err := t.Root.Find(data)
		return r, err
	} else {
		return nil, errors.New("the tree is empty")
	}
}

func (n *Node) Find(data int) (*Node, error) {
	var result *Node
	errNoKey := errors.New("there is no such key in the tree")
	if n != nil {
		if data != n.Key {
			if n.Left == nil && n.Right == nil {
				//fmt.Println("Key not found - no more trees")
				return nil, errNoKey
			} else if data < n.Key {
				if n.Left != nil {
					//fmt.Println("Left tree not nil and data less than CURRENT NODE key")
					//fmt.Println("searching in LEFT tree")
					l, err := n.Left.Find(data)
					return l, err
				} else {
					//fmt.Println("Key not found - no more trees from the left")
					return nil, errNoKey
				}
			} else if data > n.Key {
				if n.Right != nil {
					//fmt.Println("Right tree not nil and data is greater than CURRENT NODE key")
					//fmt.Println("searching in RIGHT tree")
					r, err := n.Right.Find(data)
					return r, err
				} else {
					//fmt.Println("Key not found - no more trees from the right")
					return nil, errNoKey
				}
			}
		} else {
			result = n
			//Here - insert the additional search for duplicates, smth like:
			// d,err :=n.Find(result.Key)
			// return d,err
			// with implemented counter of duplicates
		}
	} else {
		return nil, errors.New("the node is empty")
	}
	return result, nil
}

func (t *Tree) Remove(key int) error {
	// checks if the tree is not empty and countains the specified key
	node, err := t.Find(key)
	if node == nil && err != nil {
		fmt.Println("Result of search:", err)
		return err
	} else if node != nil && err == nil {
		return removeNode(&t.Root, key)
	}
	return errors.New("unexpected error")
}

func removeNode(node **Node, key int) error {
	n := *node

	//recursion to get the specified node
	if *node == nil {
		return errors.New("empty BST")
	} else if n.Key > key {
		removeNode(&n.Left, key)
	} else if n.Key < key {
		removeNode(&n.Right, key)
	}
	//deleting process
	if n.Key == key {
		//check whether the node has children
		if n.Left == nil && n.Right == nil {
			//fmt.Println("deleting the leaf", *node)
			*node = nil
		} else if n.Left != nil && n.Right == nil {
			//fmt.Println("check the condition - only left child")
			//fmt.Println("Assigning left child to the node")
			*node = n.Left
		} else if n.Left == nil && n.Right != nil {
			//fmt.Println("check the condition - only rigt child")
			//fmt.Println("Assigning right child to the node")
			*node = n.Right
		} else if n.Left != nil && n.Right != nil {
			//fmt.Println("check the condition - has both children")
			//fmt.Println("Searching for the biggest child from the left tree to assign to deleted node")
			b := searchBiggest(n.Left)
			//removing the found biggest leaf from the left tree
			removeNode(&n.Left, b.Key)
			//Assigning the new key(the biggest from left tree) to the node with the deleted key
			n.Key = b.Key
		}

	}
	return errors.New("unexpected error occurred")
}

func searchBiggest(node *Node) *Node {
	if node.Right != nil {
		node = searchBiggest(node.Right)
	}
	return node
}
