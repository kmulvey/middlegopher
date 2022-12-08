package main

import (
	"fmt"

	"github.com/kmulvey/middlegopher"
)

func main() {
	var input = make(chan int)

	// This middleware demonstrates how to make extra data available to the middleware func.
	// Examples: app config, database handles, etc.
	var one = func(incrementBy int) func(input chan int, output chan int) {
		return func(input chan int, output chan int) {
			defer close(output)
			for num := range input {
				fmt.Printf("one %d \n", num+incrementBy)
				output <- num
			}
		}
	}
	var two = func(input chan int, output chan int) {
		defer close(output)
		for num := range input {
			fmt.Printf("two %d \n", num)
			output <- num
		}
	}

	var cm = middlegopher.New(input, one(10), two)

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

/*
type dummy struct {
	Name string
}

func printer[T any](data T) {
	var d, ok = data.(dummy)
	fmt.Println(d.Name, ok)
}
*/
