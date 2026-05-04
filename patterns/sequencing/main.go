package main

import (
	"fmt"
	"time"
)

type Step struct {
	From   string
	Output string
	Done   chan struct{}
}

func runner(name string) <-chan Step {
	out := make(chan Step)
	go func() {
		defer close(out)
		for i := 1; i <= 3; i++ {
			done := make(chan struct{})
			fmt.Printf("[%s] producing step %d\n", name, i)
			out <- Step{From: name, Output: fmt.Sprintf("%s step %d", name, i), Done: done}
			<-done
			fmt.Printf("[%s] coordinator acknowledged step %d\n", name, i)
			time.Sleep(50 * time.Millisecond)
		}
	}()
	return out
}

func main() {
	a := runner("A")
	b := runner("B")

	fmt.Println("[main] sequencing A and B step by step")
	for round := 1; round <= 3; round++ {
		stepA := <-a
		stepB := <-b
		fmt.Printf("[main] round %d: %s | %s\n", round, stepA.Output, stepB.Output)
		close(stepA.Done)
		close(stepB.Done)
	}
	fmt.Println("[main] done")
}
