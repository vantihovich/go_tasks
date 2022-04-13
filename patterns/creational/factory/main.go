package main

import "fmt"

type commonTransport interface {
	setName(name string)
	getName() string
	setMaxWeight(maxWeight string)
	getMaxWeight() string
	setType(transpType string)
	getType() string
}

type transport struct {
	name       string
	maxWeight  string
	transpType string
}

func (t *transport) setName(name string) {
	t.name = name
}

func (t *transport) getName() string {
	return t.name
}

func (t *transport) setMaxWeight(maxWeight string) {
	t.maxWeight = maxWeight
}

func (t *transport) getMaxWeight() string {
	return t.maxWeight
}

func (t *transport) setType(transpType string) {
	t.transpType = transpType
}

func (t *transport) getType() string {
	return t.transpType
}

type vehicle struct {
	transport
}

func newVehicle() commonTransport {
	return &vehicle{
		transport: transport{
			name:       "TranspKings",
			maxWeight:  "20",
			transpType: "truck",
		},
	}
}

type ship struct {
	transport
}

func newShip() commonTransport {
	return &ship{
		transport: transport{
			name:       "Very Big Ship",
			maxWeight:  "400000",
			transpType: "ship",
		},
	}
}

func getTransport(transpType string) (commonTransport, error) {
	if transpType == "ship" {
		return newShip(), nil
	} else if transpType == "truck" {
		return newVehicle(), nil
	} else {
		return nil, fmt.Errorf("wrong transport type passed")
	}
}

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
