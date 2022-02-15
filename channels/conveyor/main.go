package main

import (
	"fmt"
	"math"
)

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go Generate(ch1)
	go Square(ch1, ch2)

	for num := range ch2 {
		fmt.Println(num)
	}
}

func Generate(out chan int) {
	for i := 0; i < 10; i++ {
		out <- i
	}
	close(out)
}

func Square(in, out chan int) {
	for n := range in {
		out <- int(math.Pow(float64(n), 2))
	}
	close(out)
}
