package main

// A counting semaphore caps the number of goroutines that may run a given
// section of code concurrently. The Go idiom is a buffered channel of
// empty structs: a goroutine sends to acquire a slot and receives to
// release it. The channel's buffer size IS the concurrency limit.
//
// Use this when you want to bound the load on a resource (open file
// descriptors, outbound HTTP connections, expensive CPU work) without the
// long-lived worker overhead of a worker pool. Compared to worker-pool:
//
//   - worker-pool reuses N goroutines across many jobs (good for steady
//     streams of work where startup cost matters).
//   - semaphore lets each task be its own goroutine and only gates the
//     critical section (good for ad-hoc spawned work where you just need
//     to cap simultaneous resource use).

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	const (
		maxConcurrent = 3
		totalTasks    = 10
	)

	sem := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup

	for i := 1; i <= totalTasks; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			sem <- struct{}{}        // acquire: blocks if maxConcurrent are already running
			defer func() { <-sem }() // release: defer keeps it correct even on panic

			fmt.Printf("task %2d: running\n", id)
			time.Sleep(300 * time.Millisecond)
			fmt.Printf("task %2d: done\n", id)
		}(i)
	}

	wg.Wait()
	fmt.Println("all tasks complete")
}
