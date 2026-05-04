# go-concurrency-patterns

Standalone, runnable Go examples of common concurrency patterns, plus real-world examples that combine them.

Every pattern lives in its own folder under [`patterns/`](patterns) with a `main.go` you can run and a `README.md` that explains the problem, when to use it, the goroutine/channel topology (Mermaid diagram), and what running the example looks like.

## Patterns

| Pattern | What it shows |
|---------|---------------|
| [generator](patterns/generator) | A goroutine returns a channel that yields successive values. |
| [channel-range-close](patterns/channel-range-close) | Producer closes the channel; consumer ranges until close. |
| [select-loop](patterns/select-loop) | The `for { select {} }` idiom for multiplexing channel reads. |
| [timeout](patterns/timeout) | Bound a channel receive with `time.After` inside a select. |
| [graceful-shutdown](patterns/graceful-shutdown) | Tell long-running goroutines to stop cleanly via a quit channel. |
| [sequencing](patterns/sequencing) | Coordinate goroutines step by step by sending a channel inside a message. |
| [fan-in](patterns/fan-in) | Multiplex N input channels into one stream. |
| [fan-out](patterns/fan-out) | Spread one stream of work across N concurrent workers. |
| [pipeline](patterns/pipeline) | Multi-stage pipeline connected by channels with a shared `done` for cancellation. |
| [worker-pool](patterns/worker-pool) | Fixed-size pool of long-lived workers consuming jobs from a shared channel. |
| [concurrency-limit-semaphore](patterns/concurrency-limit-semaphore) | Cap concurrent goroutines with a buffered `chan struct{}`. |
| [context-cancel](patterns/context-cancel) | `context.WithCancel` for propagating cancellation across goroutines. |
| [context-timeout](patterns/context-timeout) | `context.WithTimeout` for bounding an operation's total time. |
| [errgroup](patterns/errgroup) | `golang.org/x/sync/errgroup` for parallel tasks that can fail. |
| [mutex-vs-channel](patterns/mutex-vs-channel) | The same counter solved with `sync.Mutex` and with channel-owned state. |

## Real-world examples

| Example | What it shows |
|---------|---------------|
| [concurrent-fetcher](examples/concurrent-fetcher) | Fetch N URLs in parallel with errgroup + context, fail-fast cancellation. |
| [bounded-downloader](examples/bounded-downloader) | Download a list of files with semaphore-capped concurrency. |
| [job-queue-worker-pool](examples/job-queue-worker-pool) | Background job processing with a fixed pool of workers and a results aggregator. |

## Running an example

Each pattern and example is a `main` package, so run them with `go run`:

```bash
go run ./patterns/worker-pool
go run ./examples/job-queue-worker-pool
```

Build everything in one go:

```bash
go build ./...
```

## Requires

Go 1.26 or newer.
