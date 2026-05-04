package main

import (
	"fmt"
	"math/rand"
	"time"
)

func slowOperation() <-chan string {
	out := make(chan string)
	go func() {
		work := time.Duration(100+rand.Intn(400)) * time.Millisecond
		fmt.Printf("[op] simulating work that will take %v\n", work)
		time.Sleep(work)
		out <- "operation result"
	}()
	return out
}

func main() {
	const budget = 250 * time.Millisecond
	fmt.Printf("[main] giving operation a %v budget\n", budget)

	select {
	case result := <-slowOperation():
		fmt.Println("[main] success:", result)
	case <-time.After(budget):
		fmt.Println("[main] timeout: operation took longer than budget")
	}
}
