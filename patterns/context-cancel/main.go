package main

import (
	"context"
	"fmt"
	"time"
)

func worker(ctx context.Context, id int) {
	fmt.Printf("[worker %d] starting\n", id)
	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			fmt.Printf("[worker %d] cancelled (%v) after %d iterations\n", id, ctx.Err(), i)
			return
		default:
			fmt.Printf("[worker %d] working iteration %d\n", id, i)
			time.Sleep(150 * time.Millisecond)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	for i := 1; i <= 3; i++ {
		go worker(ctx, i)
	}

	fmt.Println("[main] workers running for 600ms then cancelling")
	time.Sleep(600 * time.Millisecond)
	cancel()
	fmt.Println("[main] cancel() called, waiting for workers to observe")
	time.Sleep(200 * time.Millisecond)
	fmt.Println("[main] done")
}
