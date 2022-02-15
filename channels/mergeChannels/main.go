package main

import (
	"fmt"
	"sync"
)

func SendTo1stChannel(out chan int) {
	for i := 0; i < 5; i++ {
		out <- i
	}
	close(out)
}

func SendTo2ndChannel(out chan int) {
	for i := 10; i < 15; i++ {
		out <- i
	}
	close(out)
}

func SendTo3rdChannel(out chan int) {
	for i := 20; i < 25; i++ {
		out <- i
	}
	close(out)
}

func joinChannels(chs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	commonChannel := make(chan int) //starting resulting common channel

	for _, ranged := range chs { //cycle that looks through incoming channels
		wg.Add(1) //adding wait group counter for each iteration
		ranged := ranged

		go func() { //routine is started for each iteration
			for n := range ranged { //cycle that looks frough data in the channel
				commonChannel <- n //sending data to one common channel
			}
			wg.Done() //decreasing wait group counter
		}()
	}

	go func() { //starting routine for waiting, when all routines from previous cycle are done
		wg.Wait()            //waiting for all routines to finish
		close(commonChannel) //closing common channel
	}()
	return commonChannel
}

func main() {
	ch1 := make(chan int) // starting 3 channels
	ch2 := make(chan int)
	ch3 := make(chan int)

	go SendTo1stChannel(ch1) //starting routine for sending ints to the 1st channel
	go SendTo2ndChannel(ch2) //starting routine for sending ints to the 2nd channel
	go SendTo3rdChannel(ch3) //starting routine for sending ints to the 3rd channel

	for s := range joinChannels(ch1, ch2, ch3) { //ranging func that puts all 3 channels together
		fmt.Println("result", s)
	}
}
