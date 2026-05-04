package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func work(ctx context.Context, name string, dur time.Duration, fail bool) error {
	fmt.Printf("[%s] starting (will take %v)\n", name, dur)
	select {
	case <-time.After(dur):
		if fail {
			fmt.Printf("[%s] failing\n", name)
			return fmt.Errorf("%s failed after %v", name, dur)
		}
		fmt.Printf("[%s] done\n", name)
		return nil
	case <-ctx.Done():
		fmt.Printf("[%s] cancelled by sibling failure: %v\n", name, ctx.Err())
		return ctx.Err()
	}
}

func main() {
	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error { return work(ctx, "fast-failing-task", 150*time.Millisecond, true) })
	g.Go(func() error { return work(ctx, "slow-task-1      ", 1*time.Second, false) })
	g.Go(func() error { return work(ctx, "slow-task-2      ", 2*time.Second, false) })

	fmt.Println("[main] waiting for group")
	if err := g.Wait(); err != nil {
		fmt.Printf("[main] group failed: %v\n", err)
		return
	}
	fmt.Println("[main] all tasks succeeded")
}
