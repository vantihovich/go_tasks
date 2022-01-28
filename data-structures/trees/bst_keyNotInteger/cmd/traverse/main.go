package main

import (
	"fmt"

	"github.com/vantihovich/go_tasks/tree/master/data-structures/trees/bst_keyNotInteger/model"
)

func main() {
	var t model.Tree


	org1 := model.Organisation{OrgName: "OrganisationOne", EmployeeNumber: 25, Resident: true}
	org2 := model.Organisation{OrgName: "OrgTwo", EmployeeNumber: 15, Resident: false}
	org3 := model.Organisation{OrgName: "OrgOne", EmployeeNumber: 14, Resident: true}
	org4 := model.Organisation{OrgName: "organisation 435", EmployeeNumber: 100, Resident: true}
	org5 := model.Organisation{OrgName: "Horror Movies", EmployeeNumber: 100, Resident: true}
	org6 := model.Organisation{OrgName: "LalaBand", EmployeeNumber: 20, Resident: true}
	org7 := model.Organisation{OrgName: "Mighty Teletubbies", EmployeeNumber: 4, Resident: false}
	org8 := model.Organisation{OrgName: "Horns & Hooves", EmployeeNumber: 5, Resident: false}
	// Insert takes struct as argument with 3 fields: 1-name of the organisation, number of employees, resident - true/false
	t.Insert(&org1)
	t.Insert(&org2)
	t.Insert(&org3)
	t.Insert(&org4)
	t.Insert(&org5)
	t.Insert(&org6)
	t.Insert(&org7)
	t.Insert(&org8)


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


func printPreOrder(n *model.Organisation) {

	if n == nil {
		return
	}

	fmt.Printf("OrgName:%q, EmployeeNum:%d, Resident-%t;\n", n.OrgName, n.EmployeeNumber, n.Resident)
	printPreOrder(n.Left)
	printPreOrder(n.Right)
}


func printPostOrder(n *model.Organisation) {

	if n == nil {
		return
	}

	printPostOrder(n.Left)
	printPostOrder(n.Right)
	fmt.Printf("OrgName:%q, EmployeeNum:%d, Resident-%t;\n", n.OrgName, n.EmployeeNumber, n.Resident)
}


func printInOrder(n *model.Organisation) {

	if n == nil {
		return
	}

	printInOrder(n.Left)
	fmt.Printf("OrgName:%q, EmployeeNum:%d, Resident-%t;\n", n.OrgName, n.EmployeeNumber, n.Resident)
	printInOrder(n.Right)
}
