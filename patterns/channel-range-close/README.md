# channel-range-close

## Problem
A producer sends a finite stream of values; the consumer needs to know when there are no more.

## When to use
- The producer knows when it is done (finite work).
- The consumer wants to use `for x := range ch` and exit naturally on completion.
- You want a clear ownership rule: the sender, and only the sender, closes the channel.

## How it works
```mermaid
flowchart LR
    P[producer] -->|send| C(chan)
    C -->|range| R[consumer]
    P -. close ch .-> C
```

`range` over a channel keeps reading until the channel is closed AND drained. Closing is the unambiguous "no more values" signal; closing twice or sending on a closed channel panics, so only the producer should close.

## Example output
```
[consumer] ranging until channel closes
[producer] sending "apple"
[consumer] got "apple"
[producer] sending "banana"
[consumer] got "banana"
[producer] sending "cherry"
[consumer] got "cherry"
[producer] sending "date"
[consumer] got "date"
[producer] no more items, closing channel
[consumer] channel closed, exiting
```

## Run it
```bash
go run ./patterns/channel-range-close
```
