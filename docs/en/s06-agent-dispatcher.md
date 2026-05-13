# s06 — Agent dispatcher

## Problem

The upstream framework lets a user (or another command) reference a sub-agent
with `@gsd-planner`. The framework must turn that token into an actual call:
find the markdown file, parse its frontmatter and `<role>` block, assemble a
prompt that includes the persona, run pre-hooks on the prompt, hand it to a
transport that talks to Claude Code, and run post-hooks on the returned text.
That whole pipeline is a single mechanism — the **dispatcher**.

## Solution

`agents/s06-agent-dispatcher/main.go` defines a `Dispatcher` configured with
an agents directory and a `gsd.Transport`. `Dispatcher.Dispatch(ref, userMsg)`
runs the full pipeline. The chapter also ships an `EchoTransport` that just
prefixes the prompt with `ECHO:\n` — useful for tests and for showing the
exact bytes that would have been sent to a real model.

`Dispatch` is the one-shot entry point most users will hit:

```go
d := NewDispatcher("./agents", EchoTransport{})
d.AddPre(s05hooks.Redact())
d.AddPost(s05hooks.Log(os.Stderr, "out"))
out, err := d.Dispatch("@gsd-planner", "Plan the auth phase")
```

## How It Works

`Resolve(ref)` strips a leading `@`, joins the agent name with `.md`, reads
the file, and parses it with the same frontmatter-and-`<role>` parser we
introduced in chapter 4. It returns a `gsd.Agent`.

`Dispatch(ref, userMsg)` is five steps:

1. `Resolve(ref)` — load the agent.
2. `buildPrompt(agent, userMsg)` — assemble the persona + user prompt into a
   deterministic string (a tighter sibling of chapter 4's `Assemble`).
3. Walk `d.preHooks`, threading the prompt through each `HookPre` in order.
4. `d.Transport.Send(prompt)` — for `EchoTransport` this just echoes back.
5. Walk `d.postHooks` on the response.

We keep the hook plumbing local to the dispatcher rather than constructing a
chapter-5 `Runner` here. It's the same shape — `for h := range hooks {…}` —
but inlining makes the chapter readable end to end without flipping between
files. A real codebase would extract this into the shared `gsd/` package.

The `EchoTransport` exists so this chapter has no external dependencies and
the test is fully hermetic. To plug in a real Claude Code or HTTP API
transport later, you only need to implement `Send(prompt string)
(gsd.TransportResult, error)` — the dispatcher doesn't care.

## What Changed

Chapter 5 gave us hooks; chapter 6 gave them somewhere to live. We now have
a function that takes a user message and an agent reference and produces a
response — exactly the surface the workflow runner in chapter 7 wants to
call.

## Try It

```sh
go run ./agents/s06-agent-dispatcher
go run ./agents/s06-agent-dispatcher @gsd-executor "Apply task 2"
go test ./agents/s06-agent-dispatcher
```

Expected output: the `ECHO:` marker followed by the assembled prompt
(`You are gsd-planner.` … `=== USER ===` … your message), with the two
demo hooks framing the prompt and response.

## Upstream Source Reading

- `sdk/src/cli-transport.ts` — the upstream's main transport (rich ANSI
  output to stdout); the upstream actually has multiple transports (`ws-`,
  `cli-`) behind a common interface, exactly like our `gsd.Transport`.
- `sdk/src/gsd-transport.ts` — the generic transport contract.
- `agents/gsd-planner.md`, `agents/gsd-executor.md`, etc. — the agent
  definitions our dispatcher would resolve.
- `sdk/src/prompt-builder.ts` — `parseAgentTools`, `parseAgentRole`.

Annotated excerpts: [`upstream-readings/s06-transport.go.md`](../../upstream-readings/s06-transport.go.md).
