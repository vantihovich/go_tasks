package main

import (
	"fmt"
	"math"
)

func main() {
	a := make(chan int)
	b := make(chan int)
	finishA := make(chan bool)
	finishB := make(chan bool)

	go Square(a, b, finishA)
	go Print(b, finishB)

	for i := 0; i < 10; i++ {
		a <- i
	}
	close(a)
	<-finishA //Catching the moment, when 2nd channel can be closed
	close(b)
	close(finishA)
	<-finishB //Waiting the moment, when 2nd goroutine finishes, not to cut the conveyor`s last results
	close(finishB)
}

func Square(in, out chan int, fin chan bool) {
	for n := range in {
		out <- int(math.Pow(float64(n), 2))
	}
	fin <- true //Need this pointer for the main to know when to close the 2nd channel, otherwise getting fatal error deadlock
}

func Print(in chan int, fin chan bool) {
	for n := range in {
		fmt.Println("result", n)
	}
	fin <- true //Need this pointer for main to know that the routine has finished its work
}
