# s08 — Spec-driven planner

## Problem

The upstream framework is opinionated about how big work happens: write a SPEC
first, lock its requirements, then break it into phases that are each small
enough to verify and ship independently. A planner reads the spec, the
executor implements one phase, the verifier checks it, and only then does the
next phase begin. This is the spec-driven discipline. This chapter ports its
smallest credible form: parse a spec into structured phases and execute one
phase at a time with a pluggable verifier.

## Solution

`agents/s08-spec-driven-planner/main.go` exposes three functions:

- `LoadSpec(path)` — parse a markdown file into a `gsd.Spec`.
- `ExecutePhase(spec, phaseID, verifier)` — run a single phase, mark it
  `Complete` on success, return the updated spec + a log.
- `NextIncomplete(spec)` — find the next phase whose `Complete` flag is
  still false; what a UI would use to drive "run the next phase" buttons.

The spec format is a flat markdown shape with a top-level `# Title`, a
`## Goal` section, and one `## Phase <ID>: <Title>` heading per phase, each
containing `Goal:` / `Tasks:` (bulleted list) / `Verify:` lines.

## How It Works

`parseSpec` is a small state machine. It walks the file line by line keeping a
state variable (`outer`, `goal`, `phase`, `tasks`). The phase heading regex
`^##\s*Phase\s+(\S+):\s*(.+?)\s*$` captures both the ID (`P01`) and the title
(`Token issuance`). Any other `## ` line closes the current Goal/Phase block.
Inside a phase, prefix-matched lines (`Goal:`, `Verify:`) populate the
respective field, `Tasks:` flips into the tasks sub-state, and dashed bullets
under it accumulate into `current.Tasks`.

`ExecutePhase` is intentionally tiny and pure-ish — it doesn't itself decide
whether a phase passes. It takes a `verifier func(gsd.SpecPhase) error`
callback. The callback is the seam: in a test the verifier returns `nil` or
`errors.New(...)` deterministically; in real life it might call `go test`,
hit a healthcheck URL, or grep a log. On success ExecutePhase flips
`p.Complete` and writes a log line; on failure it returns the verifier's
error wrapped with the phase ID.

`NextIncomplete` is one for-loop. It exists so that the integration story
("run the next phase, see what happens") has a clean entry point.

## What Changed

Chapter 7 chained commands a user knew to type; chapter 8 chains commands a
user **didn't write** — the structure comes from the SPEC file. With this
chapter the framework can take "here is my goal" and produce "here is what
the first executable step is, and here is how we will know it worked". That
is the gist of spec-driven development.

## Try It

```sh
go run ./agents/s08-spec-driven-planner
go test ./agents/s08-spec-driven-planner
```

The default demo loads `auth-refactor.md`, prints the parsed spec, then
executes `P01` with a stub verifier that succeeds when `Verify:` is non-empty.

## Upstream Source Reading

- `get-shit-done/templates/spec.md` — the canonical SPEC template. Much
  richer (ambiguity reports, falsifiability requirements, boundaries) but
  the section-shape we parse is recognisable.
- `sdk/src/plan-parser.ts` — the upstream parser. Where ours splits on `## `,
  the upstream tracks task dependencies and wave assignment for parallelism.
- `sdk/src/phase-runner.ts` — phase-by-phase execution loop.
- `agents/gsd-planner.md`, `agents/gsd-verifier.md` — the agents that produce
  and check phases.

Annotated excerpts: [`upstream-readings/s08-spec-phase.go.md`](../../upstream-readings/s08-spec-phase.go.md).
