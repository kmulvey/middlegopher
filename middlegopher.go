package middlegopher

// Middleware is the function that you will define and pass to New().
// The two arguments are 1: input chan (read), 2: output chan (write).
// Do whatever work is necessary with the data on the input chan and
// write it to the output chan or just drop it.
// IMPORTANT: dont forget to close the output channel, otherwise the whole
// pipeline will deadlock!
type Middleware[T any] func(chan T, chan T)

// MiddleGopher
type MiddleGopher[T any] struct {
	InputChan       chan T
	OutputChan      chan T
	ChanChain       []chan T
	MiddlewareFuncs []Middleware[T]
}

// New creates a new middleware pipeline. C is the input channel i.e. the start of the pipeline.
// middlewareFuncs are all the functions you want to run in the pipeline, order matters.
// The MiddleGopher is returned.
func New[T any](c chan T, middlewareFuncs ...Middleware[T]) MiddleGopher[T] {
	var cm = MiddleGopher[T]{
		InputChan:       c,
		OutputChan:      make(chan T),
		MiddlewareFuncs: middlewareFuncs,
	}

	cm.ChanChain = make([]chan T, len(middlewareFuncs))
	for i := range cm.ChanChain {
		cm.ChanChain[i] = make(chan T)
	}

	return cm
}

// Output returns the output channel. Data on this channel has gone through
// all your middleware.
func (cm *MiddleGopher[T]) Output() chan T {
	return cm.OutputChan
}

/*
func (cm *MiddleGopher[T]) GetOutput() (T, bool) {
	var val, more = <-cm.OutputChan
	return val, more
}
*/

// Run spins up all the middleware funcs as goroutines and thus starts the pipeline.
func (cm *MiddleGopher[T]) Run() {
	var previousChan chan T
	for i, fn := range cm.MiddlewareFuncs {
		if i == 0 {
			go fn(cm.InputChan, cm.ChanChain[i])
		} else if i == len(cm.MiddlewareFuncs)-1 {
			go fn(previousChan, cm.OutputChan)
		} else {
			go fn(previousChan, cm.ChanChain[i])
		}
		previousChan = cm.ChanChain[i]
	}
}
