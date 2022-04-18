package main

import "fmt"

type pizza struct {
	doughType   string
	sauseType   string
	sausageType string
	cheeseType  string
}

type pizzaBuilder interface {
	setDoughType()
	setSauseType()
	setSausageType()
	setCheeseType()
	getPizza() pizza
}

func getPizzaMaker(pizzaType string) pizzaBuilder {
	if pizzaType == "pepperoni" {
		return &pepperoniBuilder{}
	}

	if pizzaType == "vesuvius" {
		return &vesuviusBuilder{}
	}

	return nil
}

type pepperoniBuilder struct {
	doughType   string
	sauseType   string
	sausageType string
	cheeseType  string
}

func newPepperoniBuilder() *pepperoniBuilder {
	return &pepperoniBuilder{}
}

func (pb *pepperoniBuilder) setDoughType() {
	pb.doughType = "thin"
}

func (pb *pepperoniBuilder) setSauseType() {
	pb.sauseType = "tomato"
}

func (pb *pepperoniBuilder) setSausageType() {
	pb.sausageType = "pepperoni"
}

func (pb *pepperoniBuilder) setCheeseType() {
	pb.cheeseType = "mozarella"
}

func (pb *pepperoniBuilder) getPizza() pizza {
	return pizza{
		doughType:   pb.doughType,
		sauseType:   pb.sauseType,
		sausageType: pb.sausageType,
		cheeseType:  pb.cheeseType,
	}
}

type vesuviusBuilder struct {
	doughType   string
	sauseType   string
	sausageType string
	cheeseType  string
}

func newVesuviusBuilder() *vesuviusBuilder {
	return &vesuviusBuilder{}
}

func (vb *vesuviusBuilder) setDoughType() {
	vb.doughType = "thick"
}

func (vb *vesuviusBuilder) setSauseType() {
	vb.sauseType = "big tasty"
}

func (vb *vesuviusBuilder) setSausageType() {
	vb.sausageType = "ham"
}

func (vb *vesuviusBuilder) setCheeseType() {
	vb.cheeseType = "parmesan"
}

func (vb *vesuviusBuilder) getPizza() pizza {
	return pizza{
		doughType:   vb.doughType,
		sauseType:   vb.sauseType,
		sausageType: vb.sausageType,
		cheeseType:  vb.cheeseType,
	}
}

type director struct {
	builder pizzaBuilder
}

func newDirector(p pizzaBuilder) *director {
	return &director{
		builder: p,
	}
}

func (d *director) setBuilder(p pizzaBuilder) {
	d.builder = p
}

func (d *director) makePizza() pizza {
	d.builder.setDoughType()
	d.builder.setSauseType()
	d.builder.setSausageType()
	d.builder.setCheeseType()
	return d.builder.getPizza()
}

func main() {
	pepperoni := getPizzaMaker("pepperoni")
	vesuvius := getPizzaMaker("vesuvius")

	director := newDirector(pepperoni)
	pepperoniPizza := director.makePizza()

	fmt.Printf("pepperoni doughType: %s\n", pepperoniPizza.doughType)
	fmt.Printf("pepperoni sauseType: %s\n", pepperoniPizza.sauseType)
	fmt.Printf("pepperoni sausageType: %s\n", pepperoniPizza.sausageType)
	fmt.Printf("pepperoni cheeseType: %s\n", pepperoniPizza.cheeseType)

	fmt.Println("-----------")

	director.setBuilder(vesuvius)
	vesuviusPizza := director.makePizza()

	fmt.Printf("vesuvius doughType: %s\n", vesuviusPizza.doughType)
	fmt.Printf("vesuvius sauseType: %s\n", vesuviusPizza.sauseType)
	fmt.Printf("vesuvius sausageType: %s\n", vesuviusPizza.sausageType)
	fmt.Printf("vesuvius cheeseType: %s\n", vesuviusPizza.cheeseType)
}
