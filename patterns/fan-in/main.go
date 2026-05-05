package main

import (
	"fmt"
	"sync"
	"time"
)

func source(name string, n int, interval time.Duration) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for i := range n {
			out <- fmt.Sprintf("%s msg %d", name, i)
			time.Sleep(interval)
		}
	}()
	return out
}

func fanIn(inputs ...<-chan string) <-chan string {
	out := make(chan string)
	var wg sync.WaitGroup
	for _, in := range inputs {
		wg.Go(func() {
			for msg := range in {
				out <- msg
			}
		})
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	s1 := source("alpha", 4, 70*time.Millisecond)
	s2 := source("beta ", 4, 110*time.Millisecond)
	s3 := source("gamma", 4, 50*time.Millisecond)

	fmt.Println("[main] consuming merged stream from 3 sources")
	for msg := range fanIn(s1, s2, s3) {
		fmt.Println("[merged]", msg)
	}
	fmt.Println("[main] all sources drained")
}
