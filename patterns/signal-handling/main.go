package main

import (
	"context"
	"fmt"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func worker(ctx context.Context, id int) {
	fmt.Printf("[worker %d] starting\n", id)
	defer fmt.Printf("[worker %d] exited\n", id)
	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			fmt.Printf("[worker %d] shutdown received after %d ticks, draining current work\n", id, i)
			// Real cleanup goes here: release locks, flush buffers, close
			// DB connections, return in-flight items to a queue, etc.
			time.Sleep(80 * time.Millisecond)
			return
		default:
			fmt.Printf("[worker %d] tick %d\n", id, i)
			time.Sleep(120 * time.Millisecond)
		}
	}
}

func main() {
	// signal.NotifyContext returns a context cancelled when SIGINT or
	// SIGTERM arrives. In production this is the only mechanism you need.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop() // releases the signal handler

	// For this demo we also bound the run with a timeout so the example
	// terminates without manual Ctrl+C. Real code would just use the line
	// above; whichever cancellation source fires first wins.
	ctx, cancel := context.WithTimeout(ctx, 1500*time.Millisecond)
	defer cancel()

	fmt.Println("[main] starting workers (press Ctrl+C, or wait 1.5s for the demo timeout)")

	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Go(func() { worker(ctx, i) })
	}

	start := time.Now()
	<-ctx.Done()
	if time.Since(start) < 1400*time.Millisecond {
		fmt.Println("[main] shutdown triggered by signal, waiting for workers")
	} else {
		fmt.Println("[main] shutdown triggered by demo timeout, waiting for workers")
	}

	wg.Wait()
	fmt.Println("[main] all workers exited cleanly, goodbye")
}
