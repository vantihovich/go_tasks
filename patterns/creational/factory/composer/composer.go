package composer

import (
	"fmt"

	"https://github.com/vantihovich/go_tasks/tree/master/patterns/creational/factory/types"
)

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

//"factory" starts here
func getTransport(transpType string) (commonTransport, error) {
	if transpType == "ship" {
		return newShip(), nil
	} else if transpType == "truck" {
		return newVehicle(), nil
	} else {
		return nil, fmt.Errorf("wrong transport type passed")
	}
}
