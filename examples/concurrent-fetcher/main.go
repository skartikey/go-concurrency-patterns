package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

type Result struct {
	URL    string
	Status int
	Bytes  int
}

func fetch(ctx context.Context, url string) (Result, error) {
	fmt.Printf("[fetch] starting %s\n", url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Result{}, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{}, err
	}
	fmt.Printf("[fetch] done     %s -> %d (%d bytes)\n", url, resp.StatusCode, len(body))
	return Result{URL: url, Status: resp.StatusCode, Bytes: len(body)}, nil
}

func main() {
	urls := []string{
		"https://example.com",
		"https://www.iana.org/help/example-domains",
		"https://httpbin.org/status/200",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)
	results := make([]Result, len(urls))

	fmt.Printf("[main] fetching %d URLs in parallel (errgroup + context)\n", len(urls))
	for i, u := range urls {
		i, u := i, u
		g.Go(func() error {
			r, err := fetch(ctx, u)
			if err != nil {
				return fmt.Errorf("%s: %w", u, err)
			}
			results[i] = r
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("[main] group failed: %v\n", err)
		return
	}

	fmt.Println("[main] summary:")
	for _, r := range results {
		fmt.Printf("  %-50s -> %d (%d bytes)\n", r.URL, r.Status, r.Bytes)
	}
}
