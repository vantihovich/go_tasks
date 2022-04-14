package main

import (
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
	fmt.Printf("kind:%s, Name:%s", o.kind, o.name)
}

type Reader struct {
	User string
}

type ProxyObject struct {
	Object
	Reader
}

func (po *ProxyObject) Read() {
	if po.Reader.User == "Alice" {
		fmt.Println("User has an access")
		fmt.Printf("Object:%s, Reader:%s", po.Object, po.Reader)
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
