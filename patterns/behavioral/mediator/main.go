package main

import "fmt"

type (
	truck interface {
		takeRamp()
		freeRamp()
		allowTakeRamp()
	}

	mediator interface {
		canTakeRamp(truck) bool
		notiftyOnFreeRamp()
	}

	rampDispatcher struct {
		isRampClear bool
		truckQueue  []truck
	}

	groceryTruck struct {
		mediator mediator
	}

	meatTruck struct {
		mediator mediator
	}
)

func newRampDispatcher() *rampDispatcher {
	return &rampDispatcher{
		isRampClear: true,
	}
}

func (d *rampDispatcher) canTakeRamp(t truck) bool {
	if d.isRampClear {
		d.isRampClear = false
		return true
	}
	d.truckQueue = append(d.truckQueue, t)
	return false
}

func (d *rampDispatcher) notiftyOnFreeRamp() {
	if !d.isRampClear {
		d.isRampClear = true
	}
	if len(d.truckQueue) > 0 {
		firstTruckInQueue := d.truckQueue[0]
		d.truckQueue = d.truckQueue[1:]
		firstTruckInQueue.allowTakeRamp()
	}
}

func (p *groceryTruck) takeRamp() {
	if !p.mediator.canTakeRamp(p) {
		fmt.Println("Grocery truck: taking ramp BLOCKED, waiting")
		return
	}
	fmt.Println("Grocery truck: took the ramp")
}

func (p *groceryTruck) freeRamp() {
	fmt.Println("Grocery truck:leaving the ramp")
	p.mediator.notiftyOnFreeRamp()
}

func (p *groceryTruck) allowTakeRamp() {
	fmt.Println("Grocery truck: taking ramp ALLOWED, taking the ramp")
	p.takeRamp()
}

func (c *meatTruck) takeRamp() {
	if !c.mediator.canTakeRamp(c) {
		fmt.Println("Meat truck: taking ramp BLOCKED, waiting")
		return
	}
	fmt.Println("Meat truck: took the ramp")
}

func (c *meatTruck) freeRamp() {
	fmt.Println("Meat truck:leaving the ramp")
	c.mediator.notiftyOnFreeRamp()
}

func (c *meatTruck) allowTakeRamp() {
	fmt.Println("Meat truck: taking ramp ALLOWED, taking the ramp")
	c.takeRamp()
}

func main() {
	dispatcher := newRampDispatcher()

	groceryTruck := &groceryTruck{
		mediator: dispatcher,
	}

	meatTruck := &meatTruck{
		mediator: dispatcher,
	}

	groceryTruck.takeRamp()
	meatTruck.takeRamp()
	groceryTruck.freeRamp()
}
