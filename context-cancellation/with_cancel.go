package main

// context.WithCancel is the modern replacement for the quit-channel pattern.
// A single ctx.Done() channel can signal any number of child goroutines to stop,
// and the cancellation propagates through every context derived from it.

import (
	"context"
	"fmt"
	"time"
)

func worker(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("worker %d: stopping (%v)\n", id, ctx.Err())
			return
		default:
			fmt.Printf("worker %d: working...\n", id)
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	for i := 1; i <= 3; i++ {
		go worker(ctx, i)
	}

	time.Sleep(1 * time.Second)
	fmt.Println("main: cancelling workers")
	cancel()

	// Give workers a moment to observe ctx.Done() and exit cleanly.
	time.Sleep(300 * time.Millisecond)
	fmt.Println("main: done")
}
