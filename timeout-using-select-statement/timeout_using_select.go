package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Break out of channel communications after a certain period of time
// If the dynamite channel is able to deliver the Dynamite Diffused message
// before the random period of time that will execute the second case, we are safe!
// But if the time.After function delivers the current time after a random duration, the bomb will explode
func main() {
	dynamite := make(chan string)

	go func() {
		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
		dynamite <- "Dynamite Diffused!"
	}()

	for {
		select {
		case s := <-dynamite:
			fmt.Println(s)
			return
		case <-time.After(time.Duration(rand.Intn(500)) * time.Millisecond):
			fmt.Println("Dynamite Explodes!")
			return
		}
	}
}
