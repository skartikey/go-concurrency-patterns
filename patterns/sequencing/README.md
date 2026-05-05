# sequencing

## Problem
Two or more producers should advance in lockstep: each waits for the coordinator to acknowledge its current step before producing the next.

## When to use
- Round-based coordination (turn-taking games, paired pipelines).
- Backpressure where the coordinator decides when each producer may proceed.
- Demonstrating "channels as first-class values" by sending a channel inside a message.

## How it works
```mermaid
flowchart LR
    A[runner A] -->|Step{ack chan}| M[main coordinator]
    B[runner B] -->|Step{ack chan}| M
    M -.->|close ack| A
    M -.->|close ack| B
```

Each runner sends a struct that carries its result PLUS a `Done` channel. The runner blocks on `<-Done` until the coordinator closes it, then produces the next step. The ack channel is per-message, so coordination is naturally per-round.

## Example output
```
[main] sequencing A and B step by step
[B] producing step 1
[A] producing step 1
[main] round 1: A step 1 | B step 1
[A] coordinator acknowledged step 1
[B] coordinator acknowledged step 1
...
[main] done
```

## Run it
```bash
go run ./patterns/sequencing
```
