package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Job struct {
	ID      int
	Payload string
}

type Result struct {
	JobID    int
	WorkerID int
	OK       bool
	Duration time.Duration
}

func worker(id int, jobs <-chan Job, results chan<- Result) {
	fmt.Printf("[worker %d] online\n", id)
	for j := range jobs {
		start := time.Now()
		fmt.Printf("[worker %d] picked up job %d (%s)\n", id, j.ID, j.Payload)
		time.Sleep(time.Duration(100+rand.Intn(400)) * time.Millisecond)
		ok := rand.Intn(10) > 0
		if ok {
			fmt.Printf("[worker %d] job %d OK\n", id, j.ID)
		} else {
			fmt.Printf("[worker %d] job %d FAILED\n", id, j.ID)
		}
		results <- Result{JobID: j.ID, WorkerID: id, OK: ok, Duration: time.Since(start)}
	}
	fmt.Printf("[worker %d] queue closed, shutting down\n", id)
}

func main() {
	const numWorkers = 4
	const numJobs = 15

	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)

	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Go(func() { worker(w, jobs, results) })
	}

	go func() {
		for i := 1; i <= numJobs; i++ {
			fmt.Printf("[producer] enqueue job %d\n", i)
			jobs <- Job{ID: i, Payload: fmt.Sprintf("payload-%d", i)}
			time.Sleep(40 * time.Millisecond)
		}
		close(jobs)
		fmt.Println("[producer] all jobs enqueued, queue closed")
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	ok, fail := 0, 0
	var totalDur time.Duration
	for r := range results {
		if r.OK {
			ok++
		} else {
			fail++
		}
		totalDur += r.Duration
	}
	avg := totalDur / time.Duration(ok+fail)
	fmt.Printf("[main] processed %d jobs (%d ok, %d failed), avg %v per job\n",
		ok+fail, ok, fail, avg.Round(time.Millisecond))
}
