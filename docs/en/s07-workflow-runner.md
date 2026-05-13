# s07 — Workflow runner

## Problem

Some user intents map cleanly to a single command (chapter 2 covers those).
But many useful operations are chains: append a note, then list. Capture a
todo, then move it to a milestone. Run the planner, then route to the
executor. The upstream framework calls these **workflows** and declares them
as markdown files; the runtime walks the file and dispatches each step.

## Solution

`agents/s07-workflow-runner/main.go` provides `LoadWorkflow(path)` to parse a
small workflow markdown shape (YAML frontmatter for metadata + a numbered
list of `key: value` step blocks) and `Run(w, vars, runCmd)` to execute it.
Each step names a command, an args string (with `{{.Var}}` substitution), an
optional `next` step ID for transitions, and an optional `on_error` step ID
for recovery. The runner returns a textual trace and surfaces the first
unhandled error.

## How It Works

`parseWorkflow` is a two-phase pass: peel off the frontmatter (same fence
detection as chapter 1's loader), then walk the body, treating any line that
starts with `<number>.` or `<number>)` as the start of a new step. Subsequent
indented lines on the same step are read as `key: value` and dispatched into
the `WorkflowStep` fields (`id`, `command`, `args`, `next`, `on_error`,
`describe`). Anything before the first numbered marker — typical prose — is
ignored. This means you can keep narrative explanations next to the executable
parts in the same file without confusing the parser.

`Run` builds a `map[string]WorkflowStep` keyed by ID, starts at the first
step in declaration order, and follows `Next` IDs. For each step:

1. Substitute `{{.Var}}` placeholders in the args string from the `vars` map.
2. Call `runCmd(step.Command, args)`.
3. On success, write a one-line trace, follow `Next`. An empty `Next`
   terminates.
4. On error: if `OnError` is set, jump to that step and continue; otherwise
   wrap the error with the step ID and return.

A safety cap of 100 step transitions stops infinite loops dead. That's
enough headroom for any realistic workflow and small enough to surface
buggy `next:` self-references during development.

## What Changed

We've moved up one level of granularity. Chapters 1–6 each handle a single
mechanism; chapter 7 takes those mechanisms — specifically command dispatch
— and arranges them in a sequence the user didn't have to type out. The
`runCmd` callback is the seam between the workflow layer and the chapter-6
dispatcher: in real use you'd close over a `Dispatcher` and route each step
to it.

## Try It

```sh
go run ./agents/s07-workflow-runner
go run ./agents/s07-workflow-runner agents/s07-workflow-runner/testdata/with-error.md
go test ./agents/s07-workflow-runner
```

The default workflow runs `note.md`: two steps that pretend to append a note
and then list them. The `with-error.md` workflow demonstrates the recovery
path when `OnError` is set.

## Upstream Source Reading

- `get-shit-done/workflows/note.md` — the actual note workflow the upstream
  ships. Much richer than ours (multiple modes, interactive prompts) but the
  step structure is recognisably the same.
- `get-shit-done/workflows/add-todo.md` — example multi-step workflow that
  walks the planner through area inference and file output.
- `sdk/src/session-runner.ts` — session-level driver that calls workflows.

Annotated excerpts: [`upstream-readings/s07-workflow.go.md`](../../upstream-readings/s07-workflow.go.md).
