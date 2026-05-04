package main

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
	WorkerID int
	JobID    int
	Output   int
}

func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("[worker %d] online\n", id)
	for j := range jobs {
		fmt.Printf("[worker %d] picked up job %d\n", id, j.ID)
		time.Sleep(150 * time.Millisecond)
		results <- Result{WorkerID: id, JobID: j.ID, Output: j.Input * j.Input}
	}
	fmt.Printf("[worker %d] queue closed, shutting down\n", id)
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
			fmt.Printf("[producer] enqueue job %d\n", i)
			jobs <- Job{ID: i, Input: i}
		}
		close(jobs)
		fmt.Println("[producer] all jobs enqueued")
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		fmt.Printf("[result]   worker %d finished job %d -> %d\n", r.WorkerID, r.JobID, r.Output)
	}
	fmt.Println("[main] done")
}
