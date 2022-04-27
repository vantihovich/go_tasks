package main

import (
	"errors"
	"fmt"
)

type ReaderInterface interface {
	Read()
}

type Object struct {
	kind string
	name string
}

func (o *Object) Read() {
	fmt.Printf(" kind:%s, Name:%s", o.kind, o.name)
}

type Reader struct {
	User string
}

type ProxyObject struct {
	Object
	Reader
}

func (po *ProxyObject) Read() {
	err := errors.New("username not provided")

	if po.Reader.User == "Alice" {
		fmt.Printf("Reader:%s has an access to :", po.Reader)
		po.Object.Read()
	} else if po.Reader.User == "" {
		fmt.Println(err)
	}
}

var Ob = Object{
	kind: "article",
	name: "Save the world",
}

var re = Reader{
	User: "Alice",
}

var ProxyOb = ProxyObject{
	Ob,
	re,
}

func main() {
	ProxyOb.Read()
}
