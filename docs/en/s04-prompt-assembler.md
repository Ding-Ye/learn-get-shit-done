# s04 — Prompt assembler

## Problem

By the end of chapter 3 we have three structured values: a `gsd.Skill`, a
`gsd.Command`, and a rendered context string. We also have agent personas
sitting on disk in the same markdown-frontmatter shape. Before any transport
can send anything, we need to **assemble** these into one final prompt — a
single string with deterministic section order, so caching and diffing work
predictably. This chapter writes that assembler.

## Solution

`Assemble(in Inputs) string` is the public surface. Its `Inputs` struct
carries the parsed skill, parsed command, agent persona, an already-rendered
context fragment, and the raw user text. The function writes a single string
with four sections in fixed order: `=== SYSTEM ===` (agent persona),
`=== CONTEXT ===` (the rendered fragment), `=== COMMAND ===` (the skill body
+ structured invocation), and `=== USER ===` (the raw user message).

The chapter also ships a small `loadAgent` helper — about 20 lines — that
parses an agent markdown file (same frontmatter shape as a skill, but with
`tools` instead of `allowed-tools` and an optional `<role>…</role>` block in
the body). We keep it inline here on purpose: chapter 4 stands on its own,
and chapter 6 (the dispatcher) will replace this helper with a richer one.

## How It Works

The assembler is a `strings.Builder` plus a series of `fmt.Fprintf` calls. The
key decision is the **section order**. Why this order?

- SYSTEM first because a real Claude API call requires the system prompt to
  precede everything else, and we want the assembled text to read like the
  actual call structure even though we're not building an API payload yet.
- CONTEXT second because it's reusable, deterministic, and worth caching.
- COMMAND third because it's the most "this call" thing — skill body plus
  the actual parsed invocation.
- USER last because that's where free-form text lives.

Two small touches: we trim trailing whitespace on each section before writing
so we don't accumulate blank lines, and the COMMAND section prints flags as
`--name = value` (with `--name` alone for boolean flags). This makes the
output trivially scannable when you `go run` the chapter and pipe to `less`.

## What Changed

We finally have a function whose return value would be safe to hand to a
transport. The next chapters layer on top: chapter 5 wraps `Assemble` with
pre/post hooks; chapter 6 replaces the hard-coded fixture loading with a real
dispatcher; chapters 7 and 8 add workflow + spec orchestration. Everything
downstream depends on this one assembler producing the same string for the
same inputs every time.

## Try It

```sh
go run ./agents/s04-prompt-assembler
go run ./agents/s04-prompt-assembler | less
go test ./agents/s04-prompt-assembler
```

Expected output starts with `=== SYSTEM ===`, names `gsd-planner`, prints the
agent's `<role>` block, then a CONTEXT section, then a COMMAND block listing
the parsed flags (`--phase = 03`) and args (`03-auth`), and finally the USER
message.

## Upstream Source Reading

- `sdk/src/prompt-builder.ts` — see `buildExecutorPrompt`, `formatTask`,
  `parseAgentRole`. The upstream prompt-builder is the spiritual cousin of
  this chapter.
- `sdk/src/assembled-prompts.test.ts` — golden-file tests that fix the
  upstream's assembled output. Worth a skim to see how seriously they take
  determinism.
- `sdk/src/phase-prompt.ts` — phase-specific prompt assembly, layered above
  the base builder.

Annotated excerpts: [`upstream-readings/s04-prompt-builder.go.md`](../../upstream-readings/s04-prompt-builder.go.md).
