package main

import "fmt"

// we range over an array fruits and have a dummy channel case written in the select statement.
// The default case prints out the elements of the array. Hence, this technique can be used if
// we already know how many times we have to run our select statement and also if we need the
// iterative values in our code.
func main() {
	done := make(chan string)
	for _, fruit := range []string{"apple", "banana", "cherry"} {
		select {
		case <-done:
			return
		default:
			fmt.Println(fruit)
		}
	}
}
