# s01 — Skill loader

## Problem

A Claude Code framework like upstream `get-shit-done` exposes capabilities to
the model as **skills**. Each skill is a markdown file that begins with a YAML
frontmatter block describing the skill's name, summary, expected arguments, and
which tools it is allowed to use; the body is the prose the model will see when
the skill is invoked. Before any prompt assembly, hook, or dispatch happens, we
need a way to take that file on disk and produce a typed, in-memory value the
rest of the system can reason about. The first chapter is that: read a file,
return a `gsd.Skill`.

## Solution

We write a tiny, single-purpose loader in `agents/s01-skill-loader/main.go`.
It accepts a path, reads the bytes, looks for the `---` frontmatter fence,
parses the small subset of YAML the upstream actually uses (`key: value`,
plus list-of-strings under `allowed-tools:`), and returns a `gsd.Skill` struct
defined in the shared `gsd/` package. We deliberately avoid `gopkg.in/yaml.v3`
— the upstream grammar is tiny, and writing the parser ourselves keeps the
chapter to ~100 lines of pure standard-library Go.

## How It Works

`LoadSkill(path)` is the public entry point. It reads the file, then delegates
to `ParseSkill(text)` so tests can pass an in-memory string. ParseSkill:

1. Normalises line endings and splits into lines.
2. Requires the first line to be `---`. If not, returns an error.
3. Walks forward until it finds the next `---`. That index marks the end of
   frontmatter; the body is everything after.
4. Calls `parseFrontmatter` on the slice between the fences.

`parseFrontmatter` is a single pass. It tracks one bit of state — a pointer
called `currentList` — which is non-nil when the previous non-empty line was
a list-header like `allowed-tools:`. While `currentList` is active and the
current line starts with `"  - "`, we append a trimmed bullet item. Any other
line resets `currentList` to nil and is parsed as a `key: value` pair. We also
support the inline list form `allowed-tools: Read, Write, Edit` by splitting
on commas when `value` is non-empty for that key. Unknown keys are ignored on
purpose — that keeps the loader forward-compatible with new fields the
upstream may add.

The `main` function provides a thin CLI: with no args it loads the bundled
fixture under `testdata/capture.md`; with one arg it loads that path. It prints
the parsed fields so you can `go run` it and see immediately that the round
trip works.

## What Changed

This is the first chapter, so the "baseline" is the empty `gsd/` package of
shared types. The delta is: we now have a real function that turns bytes on
disk into a `gsd.Skill`. Subsequent chapters will consume that struct without
having to re-parse the file.

## Try It

```sh
go run ./agents/s01-skill-loader
go run ./agents/s01-skill-loader agents/s01-skill-loader/testdata/capture.md
go test ./agents/s01-skill-loader
```

Expected output (truncated):

```
Loaded skill: gsd:capture
  description : Capture ideas, tasks, notes, and seeds to their destination
  arg-hint    : [--note | --backlog | --seed | --list] [text]
  tools       : Read, Write, Edit, Bash
  body (first 80 chars): # Capture
```

## Upstream Source Reading

- `commands/gsd/capture.md` — the file shape we are modelling.
- `agents/gsd-planner.md` — an agent file with the same frontmatter convention
  but slightly different fields (`tools`, `color`, optional `hooks`).
- `sdk/src/prompt-builder.ts` — see `parseAgentTools` and `parseAgentRole`;
  the upstream uses regex to pull the same fields out of the same kind of file.

Annotated excerpts: [`upstream-readings/s01-skill-frontmatter.go.md`](../../upstream-readings/s01-skill-frontmatter.go.md).
