package main

import (
	"fmt"
	"time"
)

func tickerSource(name string, interval time.Duration) <-chan string {
	out := make(chan string)
	go func() {
		i := 0
		for {
			time.Sleep(interval)
			out <- fmt.Sprintf("%s msg %d", name, i)
			i++
		}
	}()
	return out
}

func main() {
	fast := tickerSource("fast", 80*time.Millisecond)
	slow := tickerSource("slow", 200*time.Millisecond)
	stop := time.After(700 * time.Millisecond)

	fmt.Println("[main] multiplexing fast + slow until stop fires")
	for {
		select {
		case msg := <-fast:
			fmt.Println("[fast]  ", msg)
		case msg := <-slow:
			fmt.Println("[slow]  ", msg)
		case <-stop:
			fmt.Println("[main] stop signal received, exiting")
			return
		}
	}
}
