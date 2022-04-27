package main

import (
	"fmt"
	"reflect"
)

type (
	// Event defines an indication of a point-in-time occurrence.
	Event struct {
		// Data in this case is a simple name of the magazine that has new edition
		Data string
	}
	// Observer defines a standard interface for instances that wish to list for
	// the occurrence of a specific event.
	Observer interface {
		// update method allows an event to be "published" to interface implementations.
		update(Event)
	}
	Publisher interface {
		addObserver(Observer)
		removeObserver(Observer)
		// notify publishes new events to listeners.
		notify(Event)
	}
)

type (
	eventObserver struct {
		name string
	}

	eventNotifier struct {
		// Using a map with an empty struct allows us to keep the observers
		// unique while still keeping memory usage relatively low.
		observers map[Observer]struct{}
	}
)

func newEvent(data string) Event {
	return Event{
		Data: data,
	}
}

func (o *eventObserver) update(e Event) {
	fmt.Printf("Sending notification to customer %s for item %q\n", o.name, e.Data)
}

func (o *eventNotifier) addObserver(l Observer) {
	o.observers[l] = struct{}{}
}

func (o *eventNotifier) removeObserver(l Observer) {
	// need to use range and reflect.DeepEqual because standard "delete(map, key)"
	// can not find the provided in l key
	for m := range o.observers {
		if reflect.DeepEqual(m, l) {
			delete(o.observers, m)
			fmt.Println("Succesfully deleted user from the dispatch list")
			return
		}
	}
	fmt.Println("Provided for deletion observer name not found in dispatch list")
}

func (p eventNotifier) notify(e Event) {
	fmt.Printf("Item %q is now available\n", e.Data)
	for o := range p.observers {
		o.update(e)
	}
}

func main() {
	//Initialize new event
	magazineEvent := newEvent("The Economist")

	// Initialize a new Notifier.
	n := eventNotifier{
		observers: map[Observer]struct{}{},
	}

	// Register some observers.
	n.addObserver(&eventObserver{name: "Alice"})
	n.addObserver(&eventObserver{name: "Tom"})
	n.addObserver(&eventObserver{name: "Kate"})

	//creating the changing state event
	n.notify(magazineEvent)

	//Remove one concrete observer to check the deletion(no such data in the dispatch list)
	n.removeObserver(&eventObserver{name: "Tommy"})

	//creating the changing state event
	n.notify(magazineEvent)

	//Remove one concrete observer to check the deletion
	n.removeObserver(&eventObserver{name: "Tom"})

	//creating the changing state event
	n.notify(magazineEvent)
}
