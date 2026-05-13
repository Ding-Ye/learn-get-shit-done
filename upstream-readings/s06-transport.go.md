# Upstream reading — s06 · Transport

Source: `gsd-build/get-shit-done` @ `ba625c0978e1133ead8a26546b5e6a435e886edb`.

## Upstream CLI transport: `sdk/src/cli-transport.ts`

```ts
export class CLITransport implements TransportHandler {
  private readonly out: Writable;
  private runningCostUsd = 0;

  constructor(out?: Writable) {
    this.out = out ?? process.stdout;
  }

  /** Format and write a GSD event as a rich ANSI-colored line. Never throws. */
  onEvent(event: GSDEvent): void {
    try {
      const line = this.formatEvent(event);
      this.out.write(line + '\n');
    } catch {
      // TransportHandler contract: onEvent must never throw
    }
  }
  // …
}
```

The upstream's Transport is **event-driven**: instead of "send a prompt, get
a response", it streams `GSDEvent`s (session-init, tool-use, message-delta,
session-complete). Our `gsd.Transport` is a string-in/string-out simplification.
The upgrade path is straightforward: replace `Send(string) (TransportResult, error)`
with a streaming variant that yields events.

## Why multiple transports?

The upstream ships at least:

- `CLITransport` — for humans at a terminal (ANSI colors, cost totals);
- `WSTransport` — for the web UI (JSON over WebSockets);
- a test/mock transport for golden-file tests.

A common interface (`TransportHandler`) makes them interchangeable. Our
chapter-6 `EchoTransport` plays the same role as a mock transport: it lets
tests assert against the exact prompt without ever calling a network.

## What to read next in upstream

- `sdk/src/cli-transport.ts` — ANSI/terminal transport.
- `sdk/src/event-stream.ts` — event types and the streaming pipeline.
- `sdk/src/gsd-transport.ts` — interface, registration, dispatch.
- `sdk/src/query/query-subprocess-adapter.ts` — how the upstream actually
  shells out to the `claude` binary when running for real.
