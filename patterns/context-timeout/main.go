package main

import (
	"context"
	"fmt"
	"time"
)

func slowOperation(ctx context.Context, dur time.Duration) (string, error) {
	fmt.Printf("[op] starting work that will take %v\n", dur)
	select {
	case <-time.After(dur):
		return fmt.Sprintf("completed after %v", dur), nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func main() {
	const budget = 300 * time.Millisecond
	const workDuration = 1 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), budget)
	defer cancel()

	fmt.Printf("[main] running op with %v budget but it actually needs %v\n", budget, workDuration)
	result, err := slowOperation(ctx, workDuration)
	if err != nil {
		fmt.Printf("[main] operation failed: %v\n", err)
		return
	}
	fmt.Printf("[main] operation succeeded: %s\n", result)
}
