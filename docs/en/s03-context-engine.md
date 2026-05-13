# s03 — Context engine

## Problem

A meta-prompting framework rarely produces a final prompt from scratch each
time. Instead, it composes one out of named, reusable fragments — "the dev
context", "the review context", "the planning context". Each fragment is a
markdown file with placeholders for the variables that change per call (the
user's name, the task, the target file). We need a small engine that can
**load a fragment by name and render it with caller-supplied variables**.

## Solution

`agents/s03-context-engine/main.go` exposes an `Engine` rooted at a directory.
`Engine.Load(name)` reads `<name>.md`, peels off an optional `---` frontmatter
block, and stores the body plus any frontmatter `key: value` pairs as default
variables. `Engine.Render(name, vars)` parses the body as a Go `text/template`,
merges defaults with caller-supplied `vars` (caller wins), and executes the
template with `Option("missingkey=error")` so a missing variable produces a
clear error rather than the literal string `<no value>`.

## How It Works

`parseContext` is essentially the chapter-1 frontmatter parser, simplified for
this use case: we only care about `key: value` pairs, not list-of-strings.
Anything we can't parse we ignore on purpose.

`Render` does three things in order:

1. Load the template (calls `Load`).
2. Build a merged `map[string]string` of defaults + caller vars, with the
   caller's value taking precedence when both are set.
3. Parse the body as a `text/template`, configure it with
   `Option("missingkey=error")`, and execute it against the merged map. The
   result is the rendered fragment.

Why `text/template` and not regex substitution? Two reasons. First, the
delimiter `{{.Var}}` is familiar to most Go programmers and survives a
copy-paste round trip through markdown reviewers. Second, we get conditionals
and ranges for free — handy when future chapters need a context fragment that
loops over, say, a list of files.

## What Changed

We now have a way to produce **text** from named templates. Combined with the
parsed `gsd.Skill` (chapter 1) and `gsd.Command` (chapter 2), we are ready to
build the prompt assembler in chapter 4: command tells us what context to
load, the engine renders it, the assembler glues it all together.

## Try It

```sh
go run ./agents/s03-context-engine dev Ding
go run ./agents/s03-context-engine review
go test ./agents/s03-context-engine
```

Expected output for the `dev` invocation includes "Hello Ding." plus the task
sentence with the placeholder filled in.

## Upstream Source Reading

- `sdk/src/context-engine.ts` — the upstream context engine. Much more
  ambitious than ours: it knows about per-phase manifests, truncates large
  files for cache-friendliness, and narrows ROADMAP.md to the current
  milestone. We deliberately drop all of that and keep only the core idea.
- `get-shit-done/contexts/dev.md`, `get-shit-done/contexts/research.md`,
  `get-shit-done/contexts/review.md` — example fragments.

Annotated excerpts: [`upstream-readings/s03-context-engine.go.md`](../../upstream-readings/s03-context-engine.go.md).
