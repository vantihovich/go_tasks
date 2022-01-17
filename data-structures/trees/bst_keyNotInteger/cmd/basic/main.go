package main

import (
	"fmt"

	"github.com/vantihovich/go_tasks/tree/master/data-structures/trees/bst_keyNotInteger/model"
)

func main() {
	var t model.Tree

	org1 := model.Node{OrgName: "OrganisationOne", EmployeeNumber: 25, Resident: true}
	org2 := model.Node{OrgName: "OrgTwo", EmployeeNumber: 15, Resident: false}
	org3 := model.Node{OrgName: "organisation 435", EmployeeNumber: 100, Resident: true}
	org4 := model.Node{OrgName: "Horns & Hooves", EmployeeNumber: 4, Resident: true}
	// Insert takes struct as argument with 3 fields: 1-name of the organisation, number of employees, resident - true/false
	t.Insert(org1)
	t.Insert(org2)
	t.Insert(org3)
	t.Insert(org4)

	fmt.Printf("%+v", t.Root.OrgName)
}
