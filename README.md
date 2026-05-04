# go-concurrency-patterns

Standalone, runnable Go examples of common concurrency patterns. Each pattern lives in its own directory.

## Patterns

| Directory | What it shows |
|-----------|---------------|
| [generator-pattern](generator-pattern) | A goroutine returns a channel that yields successive values. |
| [range-and-close](range-and-close) | Producer closes the channel; consumer ranges until close. |
| [for-select-loop](for-select-loop) | The `for { select {} }` idiom for multiplexing channel reads. |
| [timeout-using-select-statement](timeout-using-select-statement) | Bound a channel receive with `time.After` inside a select. |
| [quit-channel](quit-channel) | Signal goroutine shutdown via a dedicated channel. |
| [sequencing](sequencing) | Coordinate goroutines by sending a channel over a channel. |
| [fan-in-fan-out](fan-in-fan-out) | Multiplex N inputs into one (fan-in); split work across N workers (fan-out). |
| [pipeline](pipeline) | Multi-stage pipeline connected by channels with a shared `done` for cancellation. |
| [worker-pool](worker-pool) | Fixed-size pool of long-lived workers consuming jobs from a shared channel. |
| [semaphore](semaphore) | Cap concurrent goroutines with a buffered `chan struct{}` (counting semaphore). |
| [context-cancellation](context-cancellation) | `context.WithCancel` and `context.WithTimeout`, the modern replacement for the quit-channel. |
| [errgroup](errgroup) | `golang.org/x/sync/errgroup` for parallel tasks that can fail. |

## Running an example

Each `.go` file is its own `package main`, so run them individually:

```bash
go run generator-pattern/fibonacci.go
go run pipeline/pipeline_parallel.go
go run errgroup/errgroup_with_context.go
```

`go build ./...` and `go test ./...` will fail because several directories declare more than one `main`. This is intentional: each file is a self-contained demo.

## Requires

Go 1.26 or newer.
