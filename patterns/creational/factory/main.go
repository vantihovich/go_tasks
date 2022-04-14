package main

import (
	"fmt"

	
	"https:/github.com/vantihovich/go_tasks/tree/master/patterns/creational/factory/composer"
	
)

func main() {
	ship, _ := getTransport("ship")
	truck, _ := getTransport("truck")

	fmt.Println("check")

	printTransport(ship)
	printTransport(truck)
}

func printTransport(ct commonTransport) {
	fmt.Printf("Type: %s; Name: %s ; MaxWeight: %s", ct.getType(), ct.getName(), ct.getMaxWeight())
	fmt.Println()
}
