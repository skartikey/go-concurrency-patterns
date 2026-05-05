package main

import (
	"fmt"
	"sync"
	"time"
)

func fanOut(jobs <-chan int, numWorkers int) <-chan string {
	out := make(chan string)
	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Go(func() {
			fmt.Printf("[worker %d] started\n", w)
			for j := range jobs {
				time.Sleep(50 * time.Millisecond)
				out <- fmt.Sprintf("worker %d processed job %d -> %d", w, j, j*j)
			}
			fmt.Printf("[worker %d] no more jobs, exiting\n", w)
		})
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	const numJobs = 12
	const numWorkers = 4

	jobs := make(chan int)
	go func() {
		for i := 1; i <= numJobs; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	fmt.Printf("[main] dispatching %d jobs to %d workers\n", numJobs, numWorkers)
	for result := range fanOut(jobs, numWorkers) {
		fmt.Println("[result]", result)
	}
	fmt.Println("[main] done")
}
