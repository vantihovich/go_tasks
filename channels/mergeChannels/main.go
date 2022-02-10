package main

import (
	"fmt"
	"time"
)

func SendA(out chan int) {
	for i := 0; i < 5; i++ {
		//fmt.Println("sendA", i)
		out <- i
		time.Sleep(time.Second)
	}
}

func SendB(out chan int) {
	for i := 5; i < 10; i++ {
		//fmt.Println("sendB", i)
		out <- i
		time.Sleep(time.Second)
	}
}

func SendC(out chan int) {
	for i := 10; i < 15; i++ {
		//fmt.Println("sendC", i)
		out <- i
		time.Sleep(1 * time.Second)
	}
}

func joinChannels( /*wg *sync.WaitGroup,*/ chs ...<-chan int) <-chan int {
	//defer wg.Done()

	common := make(chan int)

	for _, res := range chs {
		fmt.Println("RCV by join", <-res)
		res := res
		go func() {
			for {
				fmt.Println("send by join", <-res)
				common <- <-res
			}
		}()

	}
	return common
}

func main() {
	//var wg sync.WaitGroup

	a := make(chan int)
	b := make(chan int)
	c := make(chan int)

	// finA := make(chan bool)
	// finB := make(chan bool)
	// finC := make(chan bool)

	go SendA(a)
	go SendB(b)
	go SendC(c)

	//wg.Add(1)
	for num := range joinChannels( /*&wg,*/ a, b, c) {
		fmt.Println("========>>>>>>>>>>result", num)
	}

	//wg.Wait()

	//close(res)
}
