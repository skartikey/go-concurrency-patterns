package main

import (
	"fmt"
	"time"
)

func produce(items []string, out chan<- string) {
	for _, item := range items {
		fmt.Printf("[producer] sending %q\n", item)
		out <- item
		time.Sleep(80 * time.Millisecond)
	}
	fmt.Println("[producer] no more items, closing channel")
	close(out)
}

func main() {
	ch := make(chan string)
	items := []string{"apple", "banana", "cherry", "date"}

	go produce(items, ch)

	fmt.Println("[consumer] ranging until channel closes")
	for v := range ch {
		fmt.Printf("[consumer] got %q\n", v)
	}
	fmt.Println("[consumer] channel closed, exiting")
}
