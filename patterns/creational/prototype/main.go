package main

import "fmt"

type (
	node interface {
		print(string)
		clone() node
	}

	file struct {
		name string
	}

	folder struct {
		children []node
		name     string
	}
)

func (f *file) print(indentation string) {
	fmt.Println(indentation + f.name)
}

func (f *file) clone() node {
	return &file{name: f.name + "_clone"}
}

func (f *folder) print(indentation string) {
	fmt.Println(indentation + f.name)
	for _, i := range f.children {
		i.print(indentation + indentation)
	}
}

func (f *folder) clone() node {
	cloneFolder := &folder{name: f.name + "_clone"}

	var tempChildren []node
	for _, i := range f.children {
		copy := i.clone()
		tempChildren = append(tempChildren, copy)
	}

	cloneFolder.children = tempChildren
	return cloneFolder
}

func main() {
	file1 := &file{name: "File1"}
	file2 := &file{name: "File2"}
	file3 := &file{name: "File3"}

	folder1 := &folder{
		children: []node{file1},
		name:     "Folder1",
	}

	folder2 := &folder{
		children: []node{folder1, file2, file3},
		name:     "Folder2",
	}

	fmt.Println("\n The hierarchy of Folder2")
	folder2.print("   ")

	cloneFolder := folder2.clone()
	fmt.Println("\n The hierarchy of cloned Folder2")
	cloneFolder.print("   ")
}
