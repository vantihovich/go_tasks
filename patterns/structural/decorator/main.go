package main

import "fmt"

type pizza interface {
	getPrice() int
}

type pepperoni struct{}

func (p *pepperoni) getPrice() int {
	return 15
}

type tomatoTopping struct {
	pizza pizza
}

func (t *tomatoTopping) getPrice() int {
	pizzaPrice := t.pizza.getPrice()
	return pizzaPrice + 5
}

type cheeseTopping struct {
	pizza pizza
}

func (c *cheeseTopping) getPrice() int {
	pizzaPrice := c.pizza.getPrice()
	return pizzaPrice + 7
}

func main() {
	pizza := &pepperoni{}

	pizzaWithCheese := &cheeseTopping{
		pizza: pizza,
	}

	pizzaWithCheeseAndTomato := &tomatoTopping{
		pizza: pizzaWithCheese,
	}

	fmt.Printf("Price with cheese and tomato is %d\n", pizzaWithCheeseAndTomato.getPrice())
}
