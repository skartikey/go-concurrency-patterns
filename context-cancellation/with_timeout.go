package main

// context.WithTimeout automatically cancels its context after a fixed duration.
// This is the idiomatic way to bound the time spent on an operation
// (an RPC, DB query, external call) without manually wiring a timer + select.
// Always defer cancel() to release timer resources, even if the timeout fires first.

import (
	"context"
	"fmt"
	"time"
)

func slowOperation(ctx context.Context) (string, error) {
	select {
	case <-time.After(2 * time.Second):
		return "operation completed", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	result, err := slowOperation(ctx)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(result)
}
