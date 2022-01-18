package main

import (
	"fmt"
	"time"
)

func main() {
	a := make(chan int)           //1st channel
	alarmClock := make(chan bool) //stopper channel

	go CheckEven(a, alarmClock) //starting 1st goroutine that checks if received value is even

	go func() { //goroutine for timer
		time.Sleep(4 * time.Second) //timer
		alarmClock <- true          // sending the value to the stopper channel
		close(a)
		defer close(alarmClock)
	}()

	for i := 1; ; i += 1 { //endless cycle
		select {
		case _, ok := <-a: // check if the 1st channel is open, if not - stopping the main
			if !ok {
				return
			}
		default: // default sending the values to the 1st channel every 0.5 of the second
			a <- i
			time.Sleep(500 * time.Millisecond)
		}
	}

}

func CheckEven(ch chan int, alarm chan bool) {
	b := make(chan int)    //the second channel
	go PrintEven(b, alarm) //starting 2nd goroutine for printing even numbers

	for {
		select { //listening to both channels(channel with value, stopper channel)
		case <-alarm: // if alarm channel received value - close the 2nd channel and stop the routine
			close(b)
			return
		case value := <-ch: // listening to value channel
			if value%2 == 0 { //check if the value is even
				b <- value //sending the value to the 2nd channel
			}
		}
	}
}

func PrintEven(ch chan int, alarm chan bool) {
	for {
		select { //listening to both channels(channel with value, stopper channel)
		case <-alarm: // if alarm channel received value - stop the routine
			return
		case value := <-ch: // listening to value channel
			fmt.Println("check", value)
		}
	}
}
