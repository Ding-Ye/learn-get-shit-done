# Upstream reading — s07 · Workflows

Source: `gsd-build/get-shit-done` @ `ba625c0978e1133ead8a26546b5e6a435e886edb`.

## A real workflow: `get-shit-done/workflows/note.md`

```markdown
<purpose>
Zero-friction idea capture. One Write call, one confirmation line. No questions, no prompts.
…
Runs inline — no Task, no AskUserQuestion, no Bash.
</purpose>

<process>
<step name="storage_format">
**Note storage format.**
Notes are stored as individual markdown files:
- Project scope: `.planning/notes/{YYYY-MM-DD}-{slug}.md` …
- Global scope: `~/.claude/notes/{YYYY-MM-DD}-{slug}.md` …
</step>

<step name="parse_subcommand">
**Parse subcommand from $ARGUMENTS (after stripping --global).**
| Condition | Subcommand |
| Arguments are exactly `list` (case-insensitive) | **list** |
| Arguments are exactly `promote <N>` where N is a number | **promote** |
| Arguments are empty (no text at all) | **list** |
| Anything else | **append** (the text IS the note) |
</step>
</process>
```

A few observations:

- The upstream workflow file is **mostly prose**. Each `<step>` block is
  human-readable instructions that the model will follow during execution.
  The framework isn't parsing the steps into structured data — it's
  injecting the whole file as instructions for Claude.
- Our chapter-7 format is more rigid (numbered steps with key/value blocks)
  because we're executing in Go, not delegating to a model. The shape is
  intentionally tighter so the runner can deterministically walk it.
- The `<step name="parse_subcommand">` decision table maps very neatly to a
  flow of `next:` transitions in our format. You could mechanically translate
  it.

## Why the upstream still chains steps explicitly

Even though each step is prose, the workflow file controls the **order** of
execution. The framework will not "skip ahead" because the model thinks it
knows better — the prose-with-numbered-steps shape forces a predictable
sequence. Our runner enforces that order in code; the upstream enforces it
in language.

## What to read next in upstream

- `get-shit-done/workflows/note.md` — the file in full.
- `get-shit-done/workflows/add-todo.md` — multi-step with area inference.
- `sdk/src/session-runner.ts` — how a workflow becomes a session.
- `commands/gsd/capture.md` — the slash-command that picks which workflow to run.
