package main

// Parallel pipeline: split a stage across multiple goroutines (fan-out)
// then merge their outputs back into a single channel (fan-in). Use this
// when one stage is the bottleneck and the work parallelizes cleanly.
//
// merge() owns the merged output channel and closes it once every input
// channel has drained. The shared done channel propagates cancellation to
// every goroutine in every stage, so a downstream early-exit unwinds the
// whole pipeline without leaks.

import (
	"fmt"
	"sync"
)

func gen(done <-chan struct{}, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}()
	return out
}

func square(done <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <-done:
				return
			}
		}
	}()
	return out
}

func merge(done <-chan struct{}, cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	done := make(chan struct{})
	defer close(done)

	in := gen(done, 1, 2, 3, 4, 5, 6, 7, 8)

	// Two parallel square workers consuming the same input channel,
	// merged into one output stream.
	c1 := square(done, in)
	c2 := square(done, in)

	sum := 0
	for n := range merge(done, c1, c2) {
		sum += n
	}
	fmt.Println("sum of squares:", sum)
}
