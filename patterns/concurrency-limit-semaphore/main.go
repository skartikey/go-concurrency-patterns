package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	const maxConcurrent = 3
	const totalTasks = 10

	sem := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup

	fmt.Printf("[main] launching %d tasks, capped at %d concurrent\n", totalTasks, maxConcurrent)

	for i := 1; i <= totalTasks; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			fmt.Printf("[task %2d] waiting for slot\n", id)
			sem <- struct{}{}
			fmt.Printf("[task %2d] acquired slot, working\n", id)

			time.Sleep(250 * time.Millisecond)

			fmt.Printf("[task %2d] releasing slot\n", id)
			<-sem
		}(i)
	}

	wg.Wait()
	fmt.Println("[main] all tasks complete")
}
