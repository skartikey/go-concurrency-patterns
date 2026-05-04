package main

import (
	"fmt"
	"time"
)

func worker(id int, quit <-chan struct{}) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		fmt.Printf("[worker %d] starting\n", id)
		i := 0
		for {
			select {
			case <-quit:
				fmt.Printf("[worker %d] quit received, exiting after %d ticks\n", id, i)
				return
			case out <- fmt.Sprintf("worker %d tick %d", id, i):
				i++
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()
	return out
}

func main() {
	quit := make(chan struct{})

	w1 := worker(1, quit)
	w2 := worker(2, quit)

	deadline := time.After(400 * time.Millisecond)

	for {
		select {
		case msg := <-w1:
			fmt.Println("[main]", msg)
		case msg := <-w2:
			fmt.Println("[main]", msg)
		case <-deadline:
			fmt.Println("[main] signalling shutdown")
			close(quit)
			time.Sleep(100 * time.Millisecond)
			fmt.Println("[main] done")
			return
		}
	}
}
