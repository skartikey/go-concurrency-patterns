package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, tokens <-chan time.Time, jobs <-chan int, wg *sync.WaitGroup, start time.Time) {
	defer wg.Done()
	fmt.Printf("[worker %d] online\n", id)
	for j := range jobs {
		<-tokens
		elapsed := time.Since(start).Round(10 * time.Millisecond)
		fmt.Printf("[worker %d] processed job %2d at +%v\n", id, j, elapsed)
		time.Sleep(50 * time.Millisecond) // simulate work
	}
	fmt.Printf("[worker %d] queue closed, shutting down\n", id)
}

func main() {
	const ratePerSecond = 5
	const numWorkers = 3
	const numJobs = 12

	interval := time.Second / time.Duration(ratePerSecond)
	ticker := time.NewTicker(interval)
	defer ticker.Stop() // releases the ticker's goroutine; forgetting this leaks

	jobs := make(chan int)
	var wg sync.WaitGroup
	start := time.Now()

	fmt.Printf("[main] %d jobs, %d workers, rate=%d ops/sec (interval=%v)\n",
		numJobs, numWorkers, ratePerSecond, interval)

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, ticker.C, jobs, &wg, start)
	}

	go func() {
		for i := 1; i <= numJobs; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	wg.Wait()

	actual := time.Since(start).Round(10 * time.Millisecond)
	min := (time.Duration(numJobs) * interval).Round(10 * time.Millisecond)
	fmt.Printf("[main] all done in %v (theoretical min %v at this rate)\n", actual, min)
}
