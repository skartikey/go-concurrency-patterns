package main

// A pipeline is a series of stages connected by channels, where each stage
// is a goroutine that reads from an inbound channel, transforms values,
// and writes to an outbound channel. Each stage owns its outbound channel
// and closes it when no more values will be sent, so downstream `range`
// loops terminate naturally.
//
// A shared `done` channel lets a downstream consumer signal upstream
// stages to stop early (on error, short-circuit, or cancellation). Stages
// use select on every send so they can react to done while blocked.
//
// This example: gen -> square -> consumer.

import "fmt"

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

func main() {
	done := make(chan struct{})
	defer close(done)

	for n := range square(done, gen(done, 1, 2, 3, 4, 5)) {
		fmt.Println(n)
	}
}
