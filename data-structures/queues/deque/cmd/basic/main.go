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

	fmt.Printf("%+v", &t)
}
