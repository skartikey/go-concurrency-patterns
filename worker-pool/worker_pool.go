package main

// A worker pool runs a fixed number of long-lived goroutines that consume
// jobs from a shared channel. Use it to bound concurrency for tasks like
// parallel HTTP requests, CPU-bound work, or DB queries, rather than
// spawning one goroutine per job (which can exhaust resources under load).
//
// Coordination:
//   1. Producer sends jobs and closes the jobs channel when done.
//   2. Workers range over jobs, exiting their loop when the channel closes.
//   3. A WaitGroup gates closing the results channel so the consumer's
//      range loop terminates cleanly once every worker has finished.

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID    int
	Input int
}

type Result struct {
	JobID  int
	Output int
}

func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		fmt.Printf("worker %d: processing job %d\n", id, j.ID)
		time.Sleep(200 * time.Millisecond) // simulate work
		results <- Result{JobID: j.ID, Output: j.Input * j.Input}
	}
}

func main() {
	const numWorkers = 3
	const numJobs = 8

	jobs := make(chan Job)
	results := make(chan Result)

	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	go func() {
		for i := 1; i <= numJobs; i++ {
			jobs <- Job{ID: i, Input: i}
		}
		close(jobs)
	}()

	// Wait for all workers to drain jobs, then close results so the
	// consumer's range loop below terminates.
	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		fmt.Printf("result: job %d squared = %d\n", r.JobID, r.Output)
	}
}
