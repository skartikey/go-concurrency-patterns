package main

import (
	"fmt"
	"sync"
	"time"
)

// Broker fans one stream of messages out to N subscribers. A single owner
// goroutine holds the subscriber set, so no mutex is needed: subscribe,
// publish, and close all serialize through the same select.
type Broker struct {
	subscribeCh chan chan string
	publishCh   chan string
	quit        chan struct{}
}

func NewBroker() *Broker {
	b := &Broker{
		subscribeCh: make(chan chan string),
		publishCh:   make(chan string),
		quit:        make(chan struct{}),
	}
	go b.run()
	return b
}

func (b *Broker) run() {
	subs := make(map[chan string]struct{})
	for {
		select {
		case ch := <-b.subscribeCh:
			subs[ch] = struct{}{}
		case msg := <-b.publishCh:
			for ch := range subs {
				// Non-blocking send: if a subscriber is full, drop the
				// message for THAT subscriber and keep going. Slow
				// subscribers must not stall the broker.
				select {
				case ch <- msg:
				default:
					fmt.Printf("[broker] dropping message for slow subscriber\n")
				}
			}
		case <-b.quit:
			for ch := range subs {
				close(ch)
			}
			return
		}
	}
}

func (b *Broker) Subscribe() <-chan string {
	ch := make(chan string, 4) // small buffer absorbs short bursts
	b.subscribeCh <- ch
	return ch
}

func (b *Broker) Publish(msg string) { b.publishCh <- msg }
func (b *Broker) Close()             { close(b.quit) }

func main() {
	b := NewBroker()

	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Go(func() {
			sub := b.Subscribe()
			fmt.Printf("[subscriber %d] online\n", i)
			for msg := range sub {
				fmt.Printf("[subscriber %d] received %q\n", i, msg)
			}
			fmt.Printf("[subscriber %d] channel closed, exiting\n", i)
		})
	}

	// Let subscribers register before publishing so the demo is
	// deterministic. In real code, publishes before any subscriber
	// exists are simply not delivered (which is the correct semantics).
	time.Sleep(100 * time.Millisecond)

	msgs := []string{"hello", "world", "broadcast", "is", "fun"}
	for _, m := range msgs {
		fmt.Printf("[publisher] publishing %q\n", m)
		b.Publish(m)
		time.Sleep(60 * time.Millisecond)
	}

	// Give subscribers a moment to drain their buffers before shutdown.
	time.Sleep(100 * time.Millisecond)

	fmt.Println("[main] closing broker; subscribers will see channel close")
	b.Close()
	wg.Wait()
	fmt.Println("[main] done")
}
