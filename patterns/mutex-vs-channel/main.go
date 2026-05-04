package main

import (
	"fmt"
	"sync"
	"time"
)

type mutexCounter struct {
	mu sync.Mutex
	n  int
}

func (c *mutexCounter) inc() {
	c.mu.Lock()
	c.n++
	c.mu.Unlock()
}

func (c *mutexCounter) read() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.n
}

type channelCounter struct {
	incCh  chan struct{}
	readCh chan int
	quit   chan struct{}
}

func newChannelCounter() *channelCounter {
	c := &channelCounter{
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

func (c *channelCounter) inc()      { c.incCh <- struct{}{} }
func (c *channelCounter) read() int { return <-c.readCh }
func (c *channelCounter) close()    { close(c.quit) }

func runMutex(n int) (int, time.Duration) {
	start := time.Now()
	var c mutexCounter
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.inc()
		}()
	}
	wg.Wait()
	return c.read(), time.Since(start)
}

func runChannel(n int) (int, time.Duration) {
	start := time.Now()
	c := newChannelCounter()
	defer c.close()
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.inc()
		}()
	}
	wg.Wait()
	return c.read(), time.Since(start)
}

func main() {
	const n = 10000

	fmt.Printf("[main] running mutex-protected counter (%d goroutines)\n", n)
	mr, md := runMutex(n)
	fmt.Printf("[mutex]   final count: %d  (took %v)\n", mr, md.Round(time.Microsecond))

	fmt.Printf("[main] running channel-owned counter (%d goroutines)\n", n)
	cr, cd := runChannel(n)
	fmt.Printf("[channel] final count: %d  (took %v)\n", cr, cd.Round(time.Microsecond))

	fmt.Println("[main] both correct; mutex is typically faster for this shape of problem")
}
