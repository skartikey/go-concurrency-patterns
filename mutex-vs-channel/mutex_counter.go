package main

// sync.Mutex protects shared state by serializing access. It is the right
// tool when:
//
//   - the state is small and accessed in tight bursts (a counter, a map),
//   - operations are short and synchronous,
//   - you just need to prevent torn reads or writes, not coordinate the
//     lifecycles of goroutines.
//
// Compare with channel_counter.go, which owns the same state inside a
// dedicated goroutine and serves requests over channels. Both produce the
// correct result; the trade-offs are about clarity, ownership, and
// overhead. For a plain shared variable, the mutex form below is simpler
// and faster: no extra goroutine, no channel allocations, no shutdown to
// wire up.

import (
	"fmt"
	"sync"
)

type counter struct {
	mu sync.Mutex
	n  int
}

func (c *counter) inc() {
	c.mu.Lock()
	c.n++
	c.mu.Unlock()
}

func (c *counter) read() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.n
}

func main() {
	var c counter
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
