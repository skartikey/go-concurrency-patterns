package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type DownloadResult struct {
	URL      string
	Bytes    int
	Duration time.Duration
}

func download(url string) DownloadResult {
	start := time.Now()
	work := time.Duration(200+rand.Intn(400)) * time.Millisecond
	time.Sleep(work)
	return DownloadResult{
		URL:      url,
		Bytes:    1024 + rand.Intn(8192),
		Duration: time.Since(start),
	}
}

func main() {
	urls := []string{
		"file-01.zip", "file-02.zip", "file-03.zip", "file-04.zip", "file-05.zip",
		"file-06.zip", "file-07.zip", "file-08.zip", "file-09.zip", "file-10.zip",
	}

	const maxConcurrent = 3
	sem := make(chan struct{}, maxConcurrent)

	var wg sync.WaitGroup
	results := make([]DownloadResult, len(urls))

	fmt.Printf("[main] downloading %d files, %d at a time\n", len(urls), maxConcurrent)

	for i, u := range urls {
		wg.Go(func() {
			fmt.Printf("[%s] waiting for slot\n", u)
			sem <- struct{}{}
			fmt.Printf("[%s] downloading\n", u)

			results[i] = download(u)

			fmt.Printf("[%s] done in %v (%d bytes), releasing slot\n",
				u, results[i].Duration.Round(time.Millisecond), results[i].Bytes)
			<-sem
		})
	}

	wg.Wait()

	total := 0
	for _, r := range results {
		total += r.Bytes
	}
	fmt.Printf("[main] downloaded %d files, %d total bytes\n", len(results), total)
}
