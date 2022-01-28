package model

import (
	"fmt"

	"github.com/pkg/errors"
)

type DeqList struct {
	Root *Element
}

type Element struct {
	El   int
	IndF *Element
	IndR *Element
}

func (d *DeqList) Add(n int, fwd bool) {
	//Checking if the list with the root element exists if not creating a new one
	if d.Root == nil {
		d.Root = &Element{El: n}
		return
	}
	//Checking the direction of adding elements - forward (true\false)
	if fwd {
		d.Root.AddElemF(n)
		return
	}
	d.Root.AddElemB(n)
}

func (e *Element) AddElemF(n int) {
	//insert to the front
	new := &Element{El: n} //new element of the list created
	last := e.FindLast()   //found the last element from the list
	new.IndR = last        //added link to the previous element for new one
	new.IndR.IndF = new    //added link to the new element for the previous element
}

func (e *Element) AddElemB(n int) {
	//insert to the beginning
	new := &Element{El: n} //new element of the list created
	first := e.FindFirst() //found the first element from the list
	new.IndF = first       //added link to the previous element for new one
	new.IndF.IndR = new    //added link to the new element for the previous element
}

func (e *Element) FindLast() *Element {
	if e.IndF != nil {
		return e.IndF.FindLast()
	}
	return e
}

func (e *Element) FindFirst() *Element {
	if e.IndR != nil {
		return e.IndR.FindFirst()
	}
	return e
}

func (d *DeqList) PrintQueue() {

	if d.Root == nil {
		fmt.Println("There are no objects in the queue")
		return
	}
	var count int
	n := d.Root

	first := n.FindFirst()
	if first.IndF == nil {
		fmt.Printf("%d ", first.El)
		count += 1
	}

	for i := n.FindFirst(); i.IndF != nil; i = i.IndF {
		if i.IndR == nil {
			fmt.Printf("%d ", i.El)
			count += 1
		}
		fmt.Printf("%d ", i.IndF.El)
		count += 1
	}
	fmt.Println("The amount of objects in the queue is:", count)
}

var ErrEmp = errors.New("The queue is empty")
var ErrNoEl = errors.New("No such element in the queue")

func (d *DeqList) Search(n int) (*Element, error) {
	var result *Element
	var err error
	//Check if the queue is empty
	if d.Root == nil {
		return nil, ErrEmp
	}
	//Check if the root is equal to the searched int
	if d.Root.El == n {
		result = d.Root
		return result, nil
	} else {
		result, err = d.Root.SearchEl(n)
	}
	return result, err
}

func (e *Element) SearchEl(n int) (*Element, error) {
	var result *Element
	//Cycle checking starts from the 1st element, if there is no link to the next el - the cycle stops and the err is returned
	for i := e.FindFirst(); i.IndF != nil; i = i.IndF {
		//Checking for equality the first element itself
		if i.IndR == nil {
			if i.El == n {
				result = i
				return result, nil
			}
		}
		//Checking for equality the link to the next element
		if i.IndF.El == n {
			result = i.IndF
			return result, nil
		}
	}
	return result, ErrNoEl
}

func RemoveLast(e **Element) {
	var d = *e
	f := d.FindLast()

	if f.El == d.El {
		if f.IndR == nil {
			//Condition-  the last el is the only el
			*e = nil
			return
		} else {
			//Condition-  the last el is the root , but not the only el
			d.IndR.IndF = nil
			*e = d.IndR
			return
		}
	}
	//Condition-  the last el is not the root
	f.IndR.IndF = nil
	f = nil
}

func RemoveFirst(e **Element) {
	var d = *e
	f := d.FindFirst()

	if f.El == d.El {
		if f.IndF == nil {
			//Condition-  the first el is the only el
			*e = nil
			return
		} else {
			//Condition-  the first el is the root , but not the only el
			d.IndF.IndR = nil
			*e = d.IndF
			return
		}
	}

	//Condition-  the 1st el is not the root
	f.IndF.IndR = nil
	f = nil
}

func (d *DeqList) RemoveExact(r int) {
	searchRes, err := d.Search(r)
	if err != nil {
		fmt.Println("There error occurred:", err)
		return
	}
	RemoveExactFunc(&d.Root, searchRes)
}

func RemoveExactFunc(e **Element, r *Element) {
	var f = *e
	// s, err := f.Search(r)

	// if err != nil {
	// 	fmt.Println("There error occurred:", err)
	// 	return
	// }

	if r == f {
		if f.IndF == nil && f.IndR == nil {
			//Condition-  the el is the only el
			*e = nil
			return
		} else if f.IndF != nil && f.IndR == nil {
			//Condition-  the el is the root and it is first
			f.IndF.IndR = nil
			*e = f.IndF
			return
		} else if f.IndR != nil && f.IndF == nil {
			//Condition-  the el is the root and it is last
			f.IndR.IndF = nil
			*e = f.IndR
			return
		} else if f.IndF != nil && f.IndR != nil {
			//Condition-  the el is the root , but but there are els from the left and right
			f.IndF.IndR = r.IndR
			f.IndR.IndF = r.IndF
			*e = f.IndF
			return
		}

	}
	if r.IndF != nil && r.IndR == nil {
		//Condition -  element in the queue not root and is the first one
		r.IndF.IndR = nil
		r = nil
	} else if r.IndF == nil && r.IndR != nil {
		//Condition -  element in the queue not root and is the last one
		r.IndR.IndF = nil
		r = nil
	} else if r.IndF != nil && r.IndR != nil {
		//Condition -  element in the queue not root and has elements from the right and left
		r.IndF.IndR = r.IndR
		r.IndR.IndF = r.IndF
		r = f.IndF
	}

}
