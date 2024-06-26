package main

import "fmt"

// In Go, we have a range function which lets us iterate over elements in different data structures.
// Using this function, we can range over the items we receive on a channel until it is closed.
// Also, note that only the sender, not the receiver, should close the channel when it feels that it has no more values to send.
type Money struct {
	amount int
	year   int
}

func sendMoney(parent chan Money) {

	for i := 0; i <= 18; i++ {
		parent <- Money{5000, i}
	}
	close(parent)
}

func main() {
	money := make(chan Money)

	go sendMoney(money)

	for kidMoney := range money {
		fmt.Printf("Money received by kid in year %d : %d\n", kidMoney.year, kidMoney.amount)
	}
}
