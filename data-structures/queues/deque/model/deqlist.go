package model

import (
	//"errors"
	"fmt"
)

type DeqList struct {
	root *Element
}

type Element struct {
	El   int
	IndF int
	IndR int
}

func (d *DeqList) AddRoot(add int, fwd bool) {
	if d == nil {
		d.root = &Element{El: add}
		fmt.Println("the list:", d)
		return
	}
	d.root.AddElement(add, fwd)
}

func (e *Element) AddElement(n int, d bool) {
	//insert to the right
	if d {
		e = &Element{El: n, IndF: e.IndF + 1}
		fmt.Println("the el after adding to the right:", d)
		return
	}
	//insert to the left
	e = &Element{El: n, IndR: e.IndR + 1}
	fmt.Println("the el after adding to the left:", d)
}
