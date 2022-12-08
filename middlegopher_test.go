package middlegopher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {

	var input = make(chan int)
	var one = func(input chan int, output chan int) {
		defer close(output)
		for num := range input {
			output <- num + 1
		}
	}
	var two = func(input chan int, output chan int) {
		defer close(output)
		for num := range input {
			output <- num + 2
		}
	}
	var three = func(input chan int, output chan int) {
		defer close(output)
		for num := range input {
			output <- num + 3
		}
	}

	var mg = New(input, one, two, three)

	var done = make(chan struct{})
	go func() {
		defer close(done)
		var expected = []int{7, 8, 9}
		var i int
		for num := range mg.Output() {
			assert.Equal(t, expected[i], num)
			i++
		}
	}()

	mg.Run()

	input <- 1
	input <- 2
	input <- 3

	close(input)
	<-done
}
