package main

//Generators return the next value in a sequence each time they are called.
//This means that each value is available as an output before the generator computes the next value.
//Hence, this pattern is used to introduce parallelism in our program.

//both the goroutine and the main routine can execute concurrently as we print the value of the counter
//as soon as we receive it while the goroutine simultaneously computes the next value.

import (
	"fmt"
)

func foo() <-chan string {
	mychannel := make(chan string)

	go func() {
		for i := 0; ; i++ {
			mychannel <- fmt.Sprintf("%s %d", "Counter at : ", i)
		}
	}()

	return mychannel // returns the channel as returning argument
}

func main() {
	mychannel := foo() // foo() returns a channel.

	for i := 0; i < 5; i++ {
		fmt.Printf("%q\n", <-mychannel)
	}

	fmt.Println("Done with Counter")
}
