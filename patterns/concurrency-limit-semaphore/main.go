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
		wg.Go(func() {
			fmt.Printf("[task %2d] waiting for slot\n", i)
			sem <- struct{}{}
			fmt.Printf("[task %2d] acquired slot, working\n", i)

			time.Sleep(250 * time.Millisecond)

			fmt.Printf("[task %2d] releasing slot\n", i)
			<-sem
		})
	}

	wg.Wait()
	fmt.Println("[main] all tasks complete")
}
