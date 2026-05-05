package main

import (
	"fmt"
	"time"
)

func fibonacci(n int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		a, b := 0, 1
		for i := range n {
			fmt.Printf("[generator] producing fib(%d) = %d\n", i, a)
			out <- a
			a, b = b, a+b
			time.Sleep(50 * time.Millisecond)
		}
		fmt.Println("[generator] done, closing channel")
	}()
	return out
}

func main() {
	fmt.Println("[main] starting consumer over fibonacci(8)")
	for v := range fibonacci(8) {
		fmt.Printf("[consumer] received %d\n", v)
	}
	fmt.Println("[main] done")
}
