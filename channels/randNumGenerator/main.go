package main

import (
	"fmt"
	"math/rand"
)

func randNumsGenerator(n int) <-chan int {
	ch := make(chan int)
	go func() {
		for i := 0; i < n; i++ {
			ch <- rand.Intn(100)
		}
		close(ch)
	}()
	return ch
}

func main() {
	for num := range randNumsGenerator(10) {
		fmt.Println("final:", num)
	}

}
