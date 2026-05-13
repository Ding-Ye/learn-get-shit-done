# s02 — Command parser

## Problem

Once we can load a skill (chapter 1), we need to actually **invoke** one. In
the upstream framework the user types something like `/gsd capture --note "Pay
the electric bill"`. That single line is the entry point of the whole system:
it names the command, optionally carries flags, and trails any positional
arguments. Until we turn it into a structured value, no downstream stage —
context engine, prompt assembler, dispatcher — has anything to act on. This
chapter is the parser that does that conversion.

## Solution

`agents/s02-command-parser/main.go` exposes `ParseCommand(line string) (gsd.Command, error)`.
It tokenises the line (honouring `"…"` and `'…'`), strips an optional leading
slash, accepts both `gsd capture` (space form) and `gsd:capture` (colon form)
as the namespace, then walks the remaining tokens, recognising `--flag`,
`--flag value`, and `--flag=value`. Bare tokens that aren't flags land in
`cmd.Args` in order. The result is a `gsd.Command` with `Name`, `Flags` (a
map), and `Args` (a slice) — exactly the shape every later chapter needs.

## How It Works

`tokenize` is a tiny state machine. It walks the input rune by rune, tracking
whether we're inside a quoted region. On a quote it flips state; on whitespace
outside quotes it flushes the current token; otherwise it accumulates. We
don't implement escape sequences — quotes pair, and that's it. An unterminated
quote returns an error so callers get deterministic failure rather than a
half-parsed result.

`ParseCommand` then peels off the first token as the namespace/name. If the
first token contains `:`, the part after the colon is the command name. If
it's literally `gsd` followed by another token, the second token is the name.
Otherwise the first token *is* the name (this lets you write `fast --text hi`
without a `gsd` prefix in workflow files).

Flag handling: a token starting with `--` becomes a flag. `--flag=value` is
split on `=`. `--flag` followed by a non-flag token consumes that token as
its value; followed by another `--flag` or end of input, it stays empty —
which we treat as "boolean: present". The first non-flag token after flags
have started still goes into Args; we don't try to reorder.

## What Changed

We now have a second runtime type from the `gsd` package being populated by a
real parser. Together with the loader from chapter 1, we can take a user's
typed line plus a skill file on disk and put both into typed Go values —
enough to start building the context and prompt layers in the next chapters.

## Try It

```sh
go run ./agents/s02-command-parser
go run ./agents/s02-command-parser '/gsd capture --note "Pay the bill"'
go run ./agents/s02-command-parser 'gsd:plan --phase=03 --gaps'
go test ./agents/s02-command-parser
```

Expected output for the second invocation:

```
raw   : /gsd capture --note "Pay the bill"
name  : capture
flags : map[note:Pay the bill]
args  : []
```

## Upstream Source Reading

- `sdk/src/cli.ts` — see `parseArgs` / `parseCliArgsQueryPermissive`. The
  upstream uses Node's `parseArgs` plus a custom permissive layer for
  forward-compat with unknown flags.
- `commands/gsd/capture.md` — see the `<context>Arguments: $ARGUMENTS</context>`
  block; the same `$ARGUMENTS` string is what gets parsed.

Annotated excerpts: [`upstream-readings/s02-cli-args.go.md`](../../upstream-readings/s02-cli-args.go.md).
