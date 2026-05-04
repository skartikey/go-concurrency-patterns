package main

// errgroup.Group runs a collection of goroutines and waits for them to
// finish, returning the first non-nil error any of them produced. Use it
// instead of bare goroutines + sync.WaitGroup whenever any of the parallel
// tasks can fail and you want to surface that failure to the caller.
//
// The bare Group does NOT cancel sibling goroutines when one fails; every
// task runs to completion and only the first error is reported. For
// automatic cancellation propagation, use WithContext (see
// errgroup_with_context.go).

import (
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	var g errgroup.Group

	tasks := []struct {
		name  string
		delay time.Duration
		fails bool
	}{
		{"task-A", 100 * time.Millisecond, false},
		{"task-B", 200 * time.Millisecond, true},
		{"task-C", 300 * time.Millisecond, false},
	}

	for _, t := range tasks {
		t := t
		g.Go(func() error {
			time.Sleep(t.delay)
			if t.fails {
				return fmt.Errorf("%s failed", t.name)
			}
			fmt.Printf("%s done\n", t.name)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Println("first error:", err)
		return
	}
	fmt.Println("all tasks succeeded")
}
