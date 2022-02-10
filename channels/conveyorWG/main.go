package main

import (
	"fmt"
	"math"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	a := make(chan int)
	b := make(chan int)

	wg.Add(1)
	go Square(a, b)
	go Print(&wg, b)

	for i := 0; i < 5; i++ {
		a <- i
	}
	close(a)

	wg.Wait()
}

func Square(in, out chan int) {
	for n := range in {
		out <- int(math.Pow(float64(n), 2))
	}

	close(out) //Chose to close the channel from here, otherwise (when closing from main) getting fatal error: deadlock, or need to catch with bool channel when the routine finishes- to close channel from main
}

func Print(wg *sync.WaitGroup, in chan int) {
	defer wg.Done()
	for n := range in {
		fmt.Println(n)
	}
}
