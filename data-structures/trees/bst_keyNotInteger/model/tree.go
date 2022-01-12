package model

import (
	"errors"
	"fmt"
)

type Tree struct {
	Root *Organisation
}

func (t *Tree) Insert(data *Organisation) {
	if t.Root == nil {
		t.Root = &Organisation{OrgName: data.OrgName, EmployeeNumber: data.EmployeeNumber, Resident: data.Resident}
		return
	}
	t.Root.Insert(data)
}

type Organisation struct {
	OrgName        string
	EmployeeNumber int
	Resident       bool
	Left           *Organisation
	Right          *Organisation
}

func (o *Organisation) Insert(data *Organisation) {
	if !compareOrganisations(data, o) {
		// insert into the left tree
		if o.Left == nil {
			o.Left = &Organisation{OrgName: data.OrgName, EmployeeNumber: data.EmployeeNumber, Resident: data.Resident}
		} else {
			o.Left.Insert(data)
		}
		return
	}
	// insert into the right tree
	if o.Right == nil {
		o.Right = &Organisation{OrgName: data.OrgName, EmployeeNumber: data.EmployeeNumber, Resident: data.Resident}
	} else {
		o.Right.Insert(data)
	}
}

func (t *Tree) Find(data string) (*Organisation, error) {
	if t.Root != nil {
		r, err := t.Root.Find(data)
		return r, err
	}
	fmt.Println("the tree is empty")
	return nil, errors.New("the tree is empty")
}

var errNoKey = errors.New("there is no such key in the tree")

func (o *Organisation) Find(data string) (*Organisation, error) {
	var result *Organisation
	assumedOrganisation := &Organisation{OrgName: data}

	if data != o.OrgName {
		if o.Left == nil && o.Right == nil {
			fmt.Println("there is no such key in the tree")
			return nil, errNoKey
		} else if !compareOrganisations(assumedOrganisation, o) {
			if o.Left != nil {
				l, err := o.Left.Find(data)
				return l, err
			} else {
				fmt.Println("there is no such key in the tree")
				return nil, errNoKey
			}
		} else if compareOrganisations(assumedOrganisation, o) {
			if o.Right != nil {
				r, err := o.Right.Find(data)
				return r, err
			} else {
				fmt.Println("there is no such key in the tree")
				return nil, errNoKey
			}
		}
	}

	result = o
	fmt.Printf("the key is found in the tree - OrgName:%q, EmployeeNum:%d, Resident-%t;\n", o.OrgName, o.EmployeeNumber, o.Resident)

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

func removeNode(node **Organisation, key string) error {
	n := *node
	assumedOrganisation := &Organisation{OrgName: key}

	//recursion to get the specified node
	if *node == nil {
		return errors.New("empty BST")
		//} else if stringsCompare(key, n.OrgName) == 3 {
	} else if !compareOrganisations(assumedOrganisation, n) {
		removeNode(&n.Left, key)
		//} else if stringsCompare(key, n.OrgName) == 2 {
	} else if compareOrganisations(assumedOrganisation, n) {
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

func searchBiggest(node *Organisation) *Organisation {
	if node.Right != nil {
		node = searchBiggest(node.Right)
	}
	return node
}

func compareOrganisations(org1, org2 *Organisation) bool {
	return org1.OrgName > org2.OrgName
}
