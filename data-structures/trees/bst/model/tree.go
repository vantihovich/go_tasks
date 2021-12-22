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
	}
	fmt.Println("the tree is empty")
	return nil, errors.New("the tree is empty")
}

var errNoKey = errors.New("there is no such key in the tree")

func (n *Node) Find(data int) (*Node, error) {
	var result *Node

	if data != n.Key {
		if n.Left == nil && n.Right == nil {
			fmt.Println("there is no such key in the tree")
			return nil, errNoKey
		} else if data < n.Key {
			if n.Left != nil {
				l, err := n.Left.Find(data)
				return l, err
			} else {
				fmt.Println("there is no such key in the tree")
				return nil, errNoKey
			}
		} else if data > n.Key {
			if n.Right != nil {
				r, err := n.Right.Find(data)
				return r, err
			} else {
				fmt.Println("there is no such key in the tree")
				return nil, errNoKey
			}
		}
	}

	result = n
	fmt.Println("the key is found in the tree:", result.Key)

	return result, nil
}

func (t *Tree) Remove(key int) error {
	// checks if the tree is not empty and countains the specified key
	node, err := t.Find(key)
	if node == nil && err != nil {
		fmt.Println("an error occurred:", err)
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
			//deleting the leaf
			*node = nil
		} else if n.Left != nil && n.Right == nil {
			//check the condition - only left child
			//Assigning left child to the node
			*node = n.Left
		} else if n.Left == nil && n.Right != nil {
			//check the condition - only rigt child
			//Assigning right child to the node
			*node = n.Right
		} else if n.Left != nil && n.Right != nil {
			//check the condition - has both children
			//Searching for the biggest child from the left tree to assign to deleted node
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
