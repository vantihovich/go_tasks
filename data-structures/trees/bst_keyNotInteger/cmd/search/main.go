package main

import (
	"fmt"

	"github.com/vantihovich/go_tasks/tree/master/data-structures/trees/bst_keyNotInteger/model"
)

func main() {
	var t model.Tree

	org1 := model.Node{OrgName: "OrganisationOne", EmployeeNumber: 25, Resident: true}
	org2 := model.Node{OrgName: "OrgTwo", EmployeeNumber: 15, Resident: false}
	org3 := model.Node{OrgName: "OrgOne", EmployeeNumber: 14, Resident: true}
	org4 := model.Node{OrgName: "organisation 435", EmployeeNumber: 100, Resident: true}
	org5 := model.Node{OrgName: "Horror Movies", EmployeeNumber: 100, Resident: true}
	org6 := model.Node{OrgName: "LalaBand", EmployeeNumber: 20, Resident: true}
	org7 := model.Node{OrgName: "Mighty Teletubbies", EmployeeNumber: 4, Resident: false}
	org8 := model.Node{OrgName: "Horns & Hooves", EmployeeNumber: 5, Resident: false}
	// Insert takes struct as argument with 3 fields: 1-name of the organisation, number of employees, resident - true/false
	t.Insert(org1)
	t.Insert(org2)
	t.Insert(org3)
	t.Insert(org4)
	t.Insert(org5)
	t.Insert(org6)
	t.Insert(org7)
	t.Insert(org8)

	fmt.Println("PREFIX_TRAVERSE example:")
	printPreOrder(t.Root)
	fmt.Println("======")

	fmt.Println("An attempt to search OrgTwo:")
	t.Find("OrgTwo")
	fmt.Println("======")

	fmt.Println("An attempt to search AAA:")
	t.Find("AAA")
	fmt.Println("======")

	fmt.Println("An attempt to search Mighty Teletubbies:")
	t.Find("Mighty Teletubbies")
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
