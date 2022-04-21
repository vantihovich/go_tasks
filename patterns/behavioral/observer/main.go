package main

import "fmt"

type publisher interface {
	addObserver(Observer observer)
	removeObserver(Observer observer)
	notify() //notifying subscriber
}

type observer interface {
	update(string)
	getID() string
}

//the subject of interest from observers side
type item struct {
	observerList []observer
	name         string
	in           bool
}

//creating new instance of item
func newItem(name string) *item {
	return &item{
		name: name,
	}
}

// updating the state of object of interest
func (i item) updateAvailability() {
	fmt.Printf("Item %q is now available\n", i.name)
	i.in = true
	i.notify() //calling notifying func for subscribed observers
}

//adding new observer to the list
func (i *item) addObserver(o observer) {
	i.observerList = append(i.observerList, o)
}

//deleting observer from the list
func (i *item) removeObserver(o observer) {
	i.observerList = removeFromSlice(i.observerList, o)
}

//notifying function that calls update method from observer side
func (i *item) notify() {
	for _, t := range i.observerList {
		t.update(i.name)
	}
}

// makes the list shorter by deleting the observer with signed name
func removeFromSlice(observerList []observer, observerToRemove observer) []observer {
	observerListLength := len(observerList)
	for i, observer := range observerList {
		if observerToRemove.getID() == observer.getID() {
			observerList[observerListLength-1], observerList[i] = observerList[i], observerList[observerListLength-1]
			return observerList[:observerListLength-1]
		}
	}
	return observerList
}

//the object that is interested in notification
type subscriber struct {
	id string
}

//the func that sends notification to the ineterested object
func (c *subscriber) update(itemName string) {
	fmt.Printf("Sending notification to customer %s for item %q\n", c.id, itemName)
}

// func for getting interested object`s id
func (c *subscriber) getID() string {
	return c.id
}

func main() {

	magazineItem := newItem("The economist") // creating new instance of object of interest

	observerOne := &subscriber{id: "Alice"} //defining "subscribers"
	observerTwo := &subscriber{id: "Tom"}
	observerThree := &subscriber{id: "Frank"}

	magazineItem.addObserver(observerOne) // adding subscribers to the "dispatch" list
	magazineItem.addObserver(observerTwo)
	magazineItem.addObserver(observerThree)

	magazineItem.updateAvailability() //creating the changing state ivent

	//now checking the new list of the observers

	magazineItem.removeObserver(observerThree) // deleting 1 subscriber

	magazineItem.updateAvailability() //creating the changing state ivent
}
