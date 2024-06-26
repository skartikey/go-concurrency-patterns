package main

import "fmt"

func getNews(newsChannel chan string) {
	NewsArray := []string{"Roger Federer wins the Wimbledon", "Space Exploration has reached new heights", "Wandering cat prevents playground accident"}
	for _, news := range NewsArray {
		newsChannel <- news
	}

	newsChannel <- "Done" // prevent an infinite loop from running
	close(newsChannel)

}

func main() {
	myNewsChannel := make(chan string)

	go getNews(myNewsChannel)

	// for-loop below enables us to execute the select statement infinitely
	for {
		select {
		case news := <-myNewsChannel:
			fmt.Println(news)
			if news == "Done" { // prevent an infinite loop from running
				return
			}
		default:
		}
	}
}
