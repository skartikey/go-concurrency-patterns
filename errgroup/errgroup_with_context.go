package main

// errgroup.WithContext returns a Group whose context is cancelled the
// moment any goroutine returns a non-nil error. Sibling goroutines that
// observe ctx.Done() can then bail out early, preventing wasted work and
// unbounded resource use after a failure.
//
// Pattern: every goroutine should select on ctx.Done() inside any blocking
// operation (or pass ctx into stdlib calls that accept it, like
// http.NewRequestWithContext or db.QueryContext).

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func work(ctx context.Context, name string, dur time.Duration, fail bool) error {
	select {
	case <-time.After(dur):
		if fail {
			return fmt.Errorf("%s failed after %v", name, dur)
		}
		fmt.Printf("%s done\n", name)
		return nil
	case <-ctx.Done():
		fmt.Printf("%s cancelled: %v\n", name, ctx.Err())
		return ctx.Err()
	}
}

func main() {
	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error { return work(ctx, "fast-task", 100*time.Millisecond, true) })
	g.Go(func() error { return work(ctx, "slow-task-1", 1*time.Second, false) })
	g.Go(func() error { return work(ctx, "slow-task-2", 2*time.Second, false) })

	if err := g.Wait(); err != nil {
		fmt.Println("group error:", err)
		return
	}
	fmt.Println("all tasks succeeded")
}
