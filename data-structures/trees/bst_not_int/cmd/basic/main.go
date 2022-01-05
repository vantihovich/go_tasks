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
	t.Insert("organisation 435", 100, true)
	t.Insert("Horns and hooves", 5, false)

	fmt.Printf("%+v", t.Root.OrgName)
}
