package main

// Channel-owned state: a single goroutine owns the variable and serves
// requests over channels. Callers send messages instead of taking a lock,
// so only one goroutine ever touches the data ("share memory by
// communicating"). No mutex is needed because there is no shared write.
//
// This shape is preferable when:
//
//   - the state already coordinates goroutine lifecycles (workers,
//     pipelines, event loops),
//   - operations are naturally request/response (commands, events),
//   - you want a clear single-owner boundary that is easy to reason about
//     and easy to extend with new operations later.
//
// For a plain counter or map, sync.Mutex (see mutex_counter.go) is
// usually simpler and faster. Pick the channel form when communication
// is the point, not when it is overhead wrapped around a guarded
// variable.

import (
	"fmt"
	"sync"
)

type counter struct {
	incCh  chan struct{}
	readCh chan int
	quit   chan struct{}
}

func newCounter() *counter {
	c := &counter{
		incCh:  make(chan struct{}),
		readCh: make(chan int),
		quit:   make(chan struct{}),
	}
	go func() {
		n := 0
		for {
			select {
			case <-c.incCh:
				n++
			case c.readCh <- n:
			case <-c.quit:
				return
			}
		}
	}()
	return c
}

func (c *counter) inc()      { c.incCh <- struct{}{} }
func (c *counter) read() int { return <-c.readCh }
func (c *counter) close()    { close(c.quit) }

func main() {
	c := newCounter()
	defer c.close()

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.inc()
		}()
	}
	wg.Wait()
	fmt.Println("final count:", c.read())
}
