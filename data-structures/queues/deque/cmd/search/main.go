package main

import (
	"fmt"

	"github.com/vantihovich/go_tasks/tree/master/data-structures/queues/deque/model"
)

func main() {
	var t model.DeqList

	//adding an element to the list , 2nd variable is about Forward adding, missed for the 1st element, if true - adds to the right, else adds to the left,
	t.Add(10, true)
	t.Add(11, true)
	t.Add(12, true)
	t.Add(9, false)
	t.Add(8, false)
	t.Add(6, false)

	//search provides result in 2 fields: 1)- the returned element(if found),2) the returned error
	fmt.Println("Search result for 6:")
	fmt.Println(t.Root.Search(6))
	fmt.Println("======")

	fmt.Println("Search result for 10:")
	fmt.Println(t.Root.Search(10))
	fmt.Println("======")

	fmt.Println("Search result for 11:")
	fmt.Println(t.Root.Search(11))
	fmt.Println("======")

	fmt.Println("Search result for 12:")
	fmt.Println(t.Root.Search(12))
	fmt.Println("======")

	fmt.Println("Search result for 24:")
	fmt.Println(t.Root.Search(24))
	fmt.Println("======")

}
