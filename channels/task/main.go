package main

import (
	"fmt"
	"time"
)

func main() {
	a := make(chan int) //1st channel
	b := make(chan int) //2nd channel

	var i int //variable for counter

	go CheckEven(a, b) //starting 1st goroutine that checks if received value is even
	go PrintEven(b)    //starting 2nd goroutine that prints even numbers

	go func() { //goroutine for timer
		time.Sleep(4 * time.Second) //timer
		close(a)
		close(b)
	}()

	for { //endless cycle
		select {
		case _, ok := <-a: // check if the 1st channel is open, if not - stopping the main
			if !ok {
				return
			}
		default: // default sending the values to the 1st channel every 0.5 of the second
			i += 1
			a <- i
			time.Sleep(500 * time.Millisecond)
		}
	}

}

func CheckEven(ch, b chan int) {
	for value := range ch {
		if value%2 == 0 { //check if the value is even
			b <- value //sending the value to the 2nd channel
		}
	}

}

func PrintEven(ch chan int) {
	for value := range ch {
		fmt.Println("check", value)
	}
}
