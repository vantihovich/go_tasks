package main

import (
	"fmt"

	"github.com/vantihovich/go_tasks/tree/master/data-structures/trees/bst_not_int/model"
)

func main() {
	var t model.Tree
	// Insert takes 3 arguments: 1-name of the organisation, number of employees, resident - true/false
	t.Insert("OrganisationOne", 25, true)
	t.Insert("OrgTwo", 15, false)
	t.Insert("OrgOne", 14, true)
	t.Insert("organisation 435", 100, true)
	t.Insert("Horror Movies", 200, false)
	t.Insert("LalaBand", 20, true)
	t.Insert("Mighty Teletubbies", 4, false)
	t.Insert("Horns and hooves", 5, false)

	fmt.Println("INFIX_TRAVERSE example:")
	printInOrder(t.Root)
	fmt.Println("======")

	fmt.Println("PREFIX_TRAVERSE example:")
	printPreOrder(t.Root)
	fmt.Println("======")

	fmt.Println("POSTFIX_TRAVERSE example:")
	printPostOrder(t.Root)
	fmt.Println("======")
}

func printPreOrder(n *model.Node) {
	if n == nil {
		return
	}

	fmt.Printf("OrgName:%q, EmployeeNum:%d, Resident-%t;\n", n.OrgName, n.EmployeeNumber, n.Resident)
	printPreOrder(n.Left)
	printPreOrder(n.Right)
}

func printPostOrder(n *model.Node) {
	if n == nil {
		return
	}

	printPostOrder(n.Left)
	printPostOrder(n.Right)
	fmt.Printf("OrgName:%q, EmployeeNum:%d, Resident-%t;\n", n.OrgName, n.EmployeeNumber, n.Resident)
}

func printInOrder(n *model.Node) {
	if n == nil {
		return
	}

	printInOrder(n.Left)
	fmt.Printf("OrgName:%q, EmployeeNum:%d, Resident-%t;\n", n.OrgName, n.EmployeeNumber, n.Resident)
	printInOrder(n.Right)
}
