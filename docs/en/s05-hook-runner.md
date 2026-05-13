# s05 — Hook runner

## Problem

Once you assemble a prompt (chapter 4), you almost never want to send it
unmodified. Real workflows demand last-mile transforms: scrub secrets,
truncate to fit a budget, log for replay, redact PII, normalize whitespace.
Hard-coding those steps into the assembler is wrong — they're cross-cutting
concerns that come and go per deployment. The standard escape hatch is a
**hook**: a small named function the framework runs before send and after
receive.

## Solution

`agents/s05-hook-runner/main.go` defines a `Runner` with an ordered slice of
`gsd.Hook`. Each hook declares its phase (`HookPre` or `HookPost`) and a
`HookFunc` that transforms a string. `Runner.Run(phase, payload)` walks the
slice, runs every hook of the matching phase in registration order, and
threads the output of one into the input of the next. The first hook to
return an error stops the chain.

We ship three concrete hooks: `Redact()` (regex-based secret scrubber, pre),
`Truncate(max)` (cap at N runes, pre), and `Log(w, label)` (write-through
logger, post).

## How It Works

`Runner.Run` is the core. It filters by `phase`, calls each matching hook
with the running payload, and propagates errors via `fmt.Errorf` so the
caller sees which hook failed. Hooks themselves are values, not interface
implementations — keeping them as plain function pointers means you can
declare a one-off hook inline (we use that trick in the error-propagation
test).

`Redact` compiles its regex once at package init via `regexp.MustCompile` and
uses `ReplaceAllString`. The pattern is intentionally generous: any token
adjacent to `api_key`, `password`, `secret`, or `token` becomes `[REDACTED]`.
False positives are cheap; false negatives are expensive.

`Truncate` operates on `[]rune` rather than bytes so multibyte input doesn't
get split mid-character. If the payload is under the cap it's returned
unchanged.

`Log` is interesting because it's a hook that does **not** mutate. It writes
the payload's size to the configured writer and returns the input as-is.
This lets you compose `Log` between two mutating hooks for free observability.

## What Changed

We now have a way to bracket the assembler output with cross-cutting
transforms. Chapter 6 uses this runner to wrap dispatch: `Run(HookPre, …)`
on the prompt before handing it to a transport, then `Run(HookPost, …)` on
the transport's response.

## Try It

```sh
go run ./agents/s05-hook-runner
go test ./agents/s05-hook-runner
```

The demo runs Redact + Truncate over a prompt that contains a fake API key
and prints the cleaned text, then runs Log over a long synthetic response.

## Upstream Source Reading

- `hooks/lib/git-cmd.js` — upstream ships hook examples as standalone JS
  files. The runtime is Node, not Go, but the contract is the same: read a
  payload from stdin, write a (possibly modified) payload to stdout.
- `sdk/src/prompt-sanitizer.ts` — the upstream's "final scrub before send"
  layer. It's a hook in spirit even when implemented inline.
- `agents/gsd-planner.md` — note the `# hooks:` commented-out block in the
  frontmatter. Agents can declare their own hooks; the framework wires them
  up when the agent is loaded.

Annotated excerpts: [`upstream-readings/s05-hooks.go.md`](../../upstream-readings/s05-hooks.go.md).
