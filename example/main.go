package main

import (
	"fmt"

	"github.com/kmulvey/middlegopher"
)

func main() {
	var input = make(chan int)
	var one = func(input chan int, output chan int) {
		defer close(output)
		for num := range input {
			fmt.Printf("one %d \n", num)
			output <- num
		}
	}
	var two = func(input chan int, output chan int) {
		defer close(output)
		for num := range input {
			fmt.Printf("two %d \n", num)
			output <- num
		}
	}

	var cm = middlegopher.New(input, one, two)

	var done = make(chan struct{})
	go func() {
		defer close(done)
		for num := range cm.Output() {
			fmt.Printf("end %d \n", num)
		}
	}()

	cm.Run()

	input <- 1
	input <- 2
	input <- 3

	close(input)
	<-done
}
