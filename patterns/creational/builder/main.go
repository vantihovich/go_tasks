package main

import (
	"fmt"
)

type IngredientType string
type PizzaIngredients map[IngredientType]*PizzaIngredient

const (
	Dough   IngredientType = "DOUGH"
	Sause   IngredientType = "SAUSE"
	Sausage IngredientType = "SAUSAGE"
	Cheese  IngredientType = "CHEESE"
)

type PizzaIngredient struct {
	Name string
	Type IngredientType
}

func (pi *PizzaIngredient) String() string {
	return fmt.Sprintf("[Ingredient name:%s, Ingredient type:%s]", pi.Name, pi.Type)
}

type Pizza struct {
	ingredients PizzaIngredients
}

func (p *Pizza) Show(i IngredientType) *PizzaIngredient {
	return p.ingredients[i]
}

type PizzaBuilder interface {
	withDough() PizzaBuilder
	withSause() PizzaBuilder
	withSausage() PizzaBuilder
	withCheese() PizzaBuilder
	build() *Pizza
}

type simplePizzaCook struct {
	ingredients PizzaIngredients
}

func newSimplePizzaCook() PizzaBuilder {
	return &simplePizzaCook{
		ingredients: make(PizzaIngredients),
	}
}

func (c *simplePizzaCook) withDough() PizzaBuilder {
	c.ingredients[Dough] = &PizzaIngredient{"Thin dough", Dough}
	return c
}

func (c *simplePizzaCook) withSause() PizzaBuilder {
	c.ingredients[Sause] = &PizzaIngredient{"Tomato sause", Sause}
	return c
}

func (c *simplePizzaCook) withSausage() PizzaBuilder {
	c.ingredients[Sausage] = &PizzaIngredient{"Pepperoni sausage", Sausage}
	return c
}

func (c *simplePizzaCook) withCheese() PizzaBuilder {
	c.ingredients[Cheese] = &PizzaIngredient{"Mozarella cheese", Cheese}
	return c
}

func (c *simplePizzaCook) build() *Pizza {
	return &Pizza{
		ingredients: c.ingredients,
	}
}

func main() {
	pizza := newSimplePizzaCook().withDough().withSause().withSausage().build()

	fmt.Println("Pizza contains:", pizza.Show(Dough))
	fmt.Println("Pizza contains:", pizza.Show(Sause))
	fmt.Println("Pizza contains:", pizza.Show(Sausage))
	fmt.Println("Pizza contains:", pizza.Show(Cheese))

	fmt.Println("Method `withCheese` was not called, so works as expected))")
}
