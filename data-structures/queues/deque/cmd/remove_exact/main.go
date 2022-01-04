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

	fmt.Println("TRAVERSE :")
	t.PrintQueue()
	fmt.Println("======")

	//removing func takes 2 variables: 1)-the queue to delete from,2) - int to be deleted
	fmt.Println("An attempt to remove ")
	model.RemoveExact(&t.Root, 6)

	fmt.Println("TRAVERSE after deleting :")
	t.PrintQueue()
	fmt.Println("======")

}
