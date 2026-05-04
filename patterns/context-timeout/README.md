# context-timeout

## Problem
Bound the total time spent on an operation. If it doesn't finish in the budget, abandon it.

## When to use
- RPC calls, DB queries, external HTTP requests.
- Anywhere you'd otherwise wire up `time.After` plus manual cancellation.
- When the deadline must propagate into downstream functions (which accept `context.Context`).

## How it works
```mermaid
flowchart LR
    M[main] -->|WithTimeout d| CTX(ctx with deadline)
    CTX --> OP[operation]
    OP -->|returns| M
    CTX -. ctx.Done after d .-> OP
```

`context.WithTimeout` builds a context that auto-cancels after the given duration. The operation `select`s on its work and `ctx.Done()`; whichever fires first wins. Always `defer cancel()` so the context's resources are released even when the operation completes early.

## Example output
```
[main] running op with 300ms budget but it actually needs 1s
[op] starting work that will take 1s
[main] operation failed: context deadline exceeded
```

## Run it
```bash
go run ./patterns/context-timeout
```
