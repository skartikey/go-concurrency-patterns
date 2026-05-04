package main

import "fmt"

func gen(done <-chan struct{}, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			fmt.Printf("[stage gen]    sending %d\n", n)
			select {
			case out <- n:
			case <-done:
				fmt.Println("[stage gen]    cancelled")
				return
			}
		}
	}()
	return out
}

func square(done <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			result := n * n
			fmt.Printf("[stage square] %d -> %d\n", n, result)
			select {
			case out <- result:
			case <-done:
				fmt.Println("[stage square] cancelled")
				return
			}
		}
	}()
	return out
}

func main() {
	done := make(chan struct{})
	defer close(done)

	fmt.Println("[main] running pipeline: gen -> square -> consumer")
	for v := range square(done, gen(done, 1, 2, 3, 4, 5)) {
		fmt.Printf("[consumer]     got %d\n", v)
	}
	fmt.Println("[main] pipeline drained")
}
