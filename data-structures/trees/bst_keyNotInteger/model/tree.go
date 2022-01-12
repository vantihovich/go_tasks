package model

import (
	"errors"
	"fmt"
)

type Tree struct {
	Root *Node
}

func (t *Tree) Insert(data Node) {
	if t.Root == nil {
		t.Root = &Node{OrgName: data.OrgName, EmployeeNumber: data.EmployeeNumber, Resident: data.Resident}
		return
	}
	t.Root.Insert(data)
}

type Node struct {
	OrgName        string
	EmployeeNumber int
	Resident       bool
	Left           *Node
	Right          *Node
}

func (n *Node) Insert(data Node) {
	if stringsCompare(data.OrgName, n.OrgName) == 3 {
		// insert into the left tree
		if n.Left == nil {
			n.Left = &Node{OrgName: data.OrgName, EmployeeNumber: data.EmployeeNumber, Resident: data.Resident}
		} else {
			n.Left.Insert(data)
		}
		return
	}
	// insert into the right tree
	if n.Right == nil {
		n.Right = &Node{OrgName: data.OrgName, EmployeeNumber: data.EmployeeNumber, Resident: data.Resident}
	} else {
		n.Right.Insert(data)
	}
}

func (t *Tree) Find(data string) (*Node, error) {
	if t.Root != nil {
		r, err := t.Root.Find(data)
		return r, err
	}
	fmt.Println("the tree is empty")
	return nil, errors.New("the tree is empty")
}

var errNoKey = errors.New("there is no such key in the tree")

func (n *Node) Find(data string) (*Node, error) {
	var result *Node

	if data != n.OrgName {
		if n.Left == nil && n.Right == nil {
			fmt.Println("there is no such key in the tree")
			return nil, errNoKey
		} else if stringsCompare(data, n.OrgName) == 3 {
			if n.Left != nil {
				l, err := n.Left.Find(data)
				return l, err
			} else {
				fmt.Println("there is no such key in the tree")
				return nil, errNoKey
			}

		} else if stringsCompare(data, n.OrgName) == 2 {
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
	fmt.Printf("the key is found in the tree - OrgName:%q, EmployeeNum:%d, Resident-%t;\n", n.OrgName, n.EmployeeNumber, n.Resident)

	return result, nil
}

func (t *Tree) Remove(key string) error {
	// checks if the tree is not empty and contains the specified key
	node, err := t.Find(key)
	if node == nil && err != nil {
		fmt.Println("an error occurred:", err)
		return err
	} else if node != nil && err == nil {
		return removeNode(&t.Root, key)
	}
	return errors.New("unexpected error")
}

func removeNode(node **Node, key string) error {
	n := *node

	//recursion to get the specified node
	if *node == nil {
		return errors.New("empty BST")
	} else if stringsCompare(key, n.OrgName) == 3 {
		removeNode(&n.Left, key)
	} else if stringsCompare(key, n.OrgName) == 2 || stringsCompare(key, n.OrgName) == 1 {
		removeNode(&n.Right, key)
	}
	//deleting process
	if n.OrgName == key {
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
			removeNode(&n.Left, b.OrgName)
			//Assigning the new key(the biggest from left tree) to the node with the deleted key
			n.OrgName = b.OrgName
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

func stringsCompare(a, b string) int {
	// The result might be 1, 2 or 3
	// 1- string a is equal to string b
	// 2- string a is bigger than string b
	// 3- string a is less than string b
	var result int
	if a == b {
		result = 1
	} else if a > b {
		result = 2
	} else if a < b {
		result = 3
	}
	return result
}
