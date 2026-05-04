# Go Concurrency Patterns (Practical Guide)

Stop guessing Go concurrency.

This repo shows **production-style concurrency patterns** with:

* runnable examples (`go run`)
* real-world use cases
* visual diagrams (Mermaid)
* clear “when to use what” guidance

---

## ⚡ Try it in 10 seconds

```bash
go run ./patterns/worker-pool
```

Example output (abridged from a real run):

```
[worker 1] online
[worker 2] online
[worker 3] online
[producer] enqueue job 1
[worker 1] picked up job 1
[producer] enqueue job 2
[worker 2] picked up job 2
[producer] enqueue job 3
[worker 3] picked up job 3
[result]   worker 1 finished job 1 -> 1
[result]   worker 2 finished job 2 -> 4
[result]   worker 3 finished job 3 -> 9
...
[main] done
```

---

## 🧠 Which pattern should I use?

| Problem | Pattern |
|---------|---------|
| Stream a sequence of values from a goroutine | [generator](patterns/generator) |
| Signal end-of-stream from a producer | [channel-range-close](patterns/channel-range-close) |
| React to multiple channels in one goroutine | [select-loop](patterns/select-loop) |
| Time-bound a single channel receive | [timeout](patterns/timeout) |
| Stop long-running goroutines cleanly | [graceful-shutdown](patterns/graceful-shutdown) |
| Coordinate goroutines step by step | [sequencing](patterns/sequencing) |
| Merge multiple channels into one | [fan-in](patterns/fan-in) |
| Spread one stream of work across N workers | [fan-out](patterns/fan-out) |
| Chain stream-processing stages | [pipeline](patterns/pipeline) |
| Process many jobs with a fixed pool of workers | [worker-pool](patterns/worker-pool) |
| Cap how many goroutines run a section at once | [concurrency-limit-semaphore](patterns/concurrency-limit-semaphore) |
| Cancel work across goroutines | [context-cancel](patterns/context-cancel) |
| Bound an operation's total time | [context-timeout](patterns/context-timeout) |
| Run parallel tasks, fail-fast on any error | [errgroup](patterns/errgroup) |
| Pick between mutex and channel for shared state | [mutex-vs-channel](patterns/mutex-vs-channel) |

---

## 📦 Patterns

Each pattern:

* is in its own folder under [`patterns/`](patterns)
* has a runnable `main.go`
* includes a README with problem, diagram, and output

| Pattern                                                             | What it shows                                               |
| ------------------------------------------------------------------- | ----------------------------------------------------------- |
| [generator](patterns/generator)                                     | A goroutine returns a channel that yields successive values |
| [channel-range-close](patterns/channel-range-close)                 | Close channel and range until done                          |
| [select-loop](patterns/select-loop)                                 | Multiplex multiple channels with `select`                   |
| [timeout](patterns/timeout)                                         | Add timeouts using `time.After`                             |
| [graceful-shutdown](patterns/graceful-shutdown)                     | Stop goroutines cleanly                                     |
| [sequencing](patterns/sequencing)                                   | Coordinate ordered execution                                |
| [fan-in](patterns/fan-in)                                           | Merge multiple channels                                     |
| [fan-out](patterns/fan-out)                                         | Distribute work across goroutines                           |
| [pipeline](patterns/pipeline)                                       | Multi-stage processing                                      |
| [worker-pool](patterns/worker-pool)                                 | Fixed pool of workers                                       |
| [concurrency-limit-semaphore](patterns/concurrency-limit-semaphore) | Limit concurrent work                                       |
| [context-cancel](patterns/context-cancel)                           | Cancel work across goroutines                               |
| [context-timeout](patterns/context-timeout)                         | Time-bound operations                                       |
| [errgroup](patterns/errgroup)                                       | Parallel tasks with error handling                          |
| [mutex-vs-channel](patterns/mutex-vs-channel)                       | Compare synchronization approaches                          |

---

## 🌍 Real-world examples

These combine multiple patterns into practical scenarios:

| Example                                                 | What it shows                         |
| ------------------------------------------------------- | ------------------------------------- |
| [concurrent-fetcher](examples/concurrent-fetcher)       | Parallel API calls with cancellation  |
| [bounded-downloader](examples/bounded-downloader)       | Controlled concurrency with semaphore |
| [job-queue-worker-pool](examples/job-queue-worker-pool) | Background job processing system      |

---

## ⚠️ Common mistakes

* **Spawning unbounded goroutines.** One goroutine per request can exhaust memory and file descriptors under load. Use [worker-pool](patterns/worker-pool) or [concurrency-limit-semaphore](patterns/concurrency-limit-semaphore) to cap concurrency.
* **Forgetting `defer cancel()`.** A context whose `cancel` is never called leaks its associated goroutine and timer until the parent context dies. Always pair `context.WithCancel` / `WithTimeout` with `defer cancel()` (see [context-cancel](patterns/context-cancel), [context-timeout](patterns/context-timeout)).
* **Goroutine leaks.** A goroutine blocked forever on a channel that no one closes or reads is invisible until your service runs out of memory. Give every goroutine an exit via [graceful-shutdown](patterns/graceful-shutdown) or context-based cancellation.
* **Blocking unintentionally on channels.** Sends on unbuffered channels block until someone receives; receives block forever if no one closes. Design with [select-loop](patterns/select-loop) and [timeout](patterns/timeout) escape hatches when a hang would be unsafe.

---

## ▶️ Run examples

```bash
go run ./patterns/worker-pool
go run ./examples/job-queue-worker-pool
```

Build everything:

```bash
go build ./...
```

---

## Requires

Go 1.26+
