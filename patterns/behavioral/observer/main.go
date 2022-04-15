package main

import "fmt"

type publisher interface {
	register(Observer observer)   //adding new subscriber
	deRegister(Observer observer) //deleting subscriber
	notify()                      //notifying subscriber
}

type observer interface { //observer interface
	update(string)
	getID() string
}

type item struct { //the subject of interest from observers side
	observerList []observer
	name         string
	in           bool
}

func newItem(name string) *item { //creating new instance of item
	return &item{
		name: name,
	}
}

func (i item) updateAvailability() { // updating the state of object of interest
	fmt.Printf("Item %q is now available\n", i.name)
	i.in = true
	i.notify() //calling notifying func for subscribed observers
}

func (i *item) register(o observer) { //adding new observer to the list
	i.observerList = append(i.observerList, o)
}

func (i *item) deRegister(o observer) { //deleting observer from the list
	i.observerList = removeFromSlice(i.observerList, o)
}

func (i *item) notify() { //notifying function that calls update method from observer side
	for _, observer := range i.observerList {
		observer.update(i.name)
	}
}

func removeFromSlice(observerList []observer, observerToRemove observer) []observer { // makes the list shorter by deleting the observer with signed name
	observerListLength := len(observerList)
	for i, observer := range observerList {
		if observerToRemove.getID() == observer.getID() {
			observerList[observerListLength-1], observerList[i] = observerList[i], observerList[observerListLength-1]
			return observerList[:observerListLength-1]
		}
	}
	return observerList
}

type customer struct { //the object that is interested in notification
	id string
}

func (c *customer) update(itemName string) { //the func that sends notification to the ineterested object
	fmt.Printf("Sending notification to customer %s for item %q\n", c.id, itemName)
}

func (c *customer) getID() string { // func for getting interested object`s id
	return c.id
}

func main() {

	magazineItem := newItem("The economist") // creating new instance of object of interest

	observerOne := &customer{id: "Alice"} //defining "subscribers"
	observerTwo := &customer{id: "Tom"}
	observerThree := &customer{id: "Frank"}

	magazineItem.register(observerOne) // adding subscribers to the "dispatch" list
	magazineItem.register(observerTwo)
	magazineItem.register(observerThree)

	magazineItem.updateAvailability() //creating the changing state ivent

	//now checking the new list of the observers

	magazineItem.deRegister(observerThree) // deleting 1 subscriber

	magazineItem.updateAvailability() //creating the changing state ivent

}
