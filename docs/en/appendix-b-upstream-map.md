# Appendix B — Upstream map: where to read what

This is a one-page index from each chapter to the upstream files most worth
reading. Upstream commit: `ba625c0978e1133ead8a26546b5e6a435e886edb` (MIT).

## Top-level entry points

- `CLAUDE.md` — top-level instructions Claude reads first.
- `CONTEXT.md` — the example project context file (also one of the files
  the upstream's chapter-3 equivalent loads per phase).
- `package.json`, `tsconfig.json` — package shape (huge dep graph; skip
  until you've understood the rest).
- `sdk/src/cli.ts` — main CLI entry. The `parseArgs` shape mirrors our
  chapter-2 `gsd.Command`.
- `sdk/src/index.ts` — the top-level `GSD` class that wires it all together.

## Per chapter

### s01 — Skill loader

- `commands/gsd/capture.md`, `commands/gsd/plan-phase.md`,
  `commands/gsd/execute-phase.md` — typical command shape.
- `agents/gsd-planner.md`, `agents/gsd-executor.md` — typical agent shape.
- `sdk/src/prompt-builder.ts` — see `parseAgentTools`, `parseAgentRole`
  (same regex-on-frontmatter trick we use).

### s02 — Command parser

- `sdk/src/cli.ts` — `ParsedCliArgs`, `parseCliArgsQueryPermissive`.
- `sdk/src/init-runner.ts` — special-cased parsing for `init` (handles `@file`).
- `commands/gsd/capture.md` — the `<context>Arguments: $ARGUMENTS</context>`
  block shows what the model sees of the parsed input.

### s03 — Context engine

- `sdk/src/context-engine.ts` — full engine with per-phase manifests.
- `sdk/src/context-truncation.ts` — heading-aware truncation.
- `get-shit-done/contexts/dev.md`, `get-shit-done/contexts/research.md`,
  `get-shit-done/contexts/review.md` — the actual context profiles.
- `get-shit-done/templates/context.md` — the template used to seed a new
  project's `CONTEXT.md`.

### s04 — Prompt assembler

- `sdk/src/prompt-builder.ts` — `buildExecutorPrompt`, `formatTask`.
- `sdk/src/assembled-prompts.test.ts` — golden-file tests on the assembled output.
- `sdk/src/phase-prompt.ts` — phase-specific prompt assembly layered above
  the base builder.
- `sdk/src/prompt-sanitizer.ts` — final scrub before send.

### s05 — Hook runner

- `hooks/lib/git-cmd.js` — example hook implementation.
- `sdk/src/prompt-sanitizer.ts` — pre-send scrubbing layer.
- `sdk/src/logger.ts` — observability hook stack.
- The commented `# hooks:` block in `agents/gsd-planner.md` — how an agent
  declares its own per-event hooks.

### s06 — Agent dispatcher

- `sdk/src/cli-transport.ts` — ANSI/terminal transport.
- `sdk/src/gsd-transport.ts` — the transport contract.
- `sdk/src/event-stream.ts` — event types and streaming.
- `sdk/src/query/query-subprocess-adapter.ts` — how the upstream shells out
  to the real `claude` binary.
- `agents/gsd-*.md` — every agent definition the dispatcher would resolve.

### s07 — Workflow runner

- `get-shit-done/workflows/note.md` — the actual note workflow.
- `get-shit-done/workflows/add-todo.md` — multi-step with area inference.
- `get-shit-done/workflows/plan-phase.md` — calls the planner sub-agent.
- `sdk/src/session-runner.ts` — session-level driver.
- `sdk/src/milestone-runner.test.ts` — tests of multi-workflow milestone runs.

### s08 — Spec-driven planner

- `get-shit-done/templates/spec.md` — canonical SPEC template.
- `get-shit-done/templates/phase-prompt.md` — phase execution prompt template.
- `sdk/src/plan-parser.ts` — full plan parser with dependency analysis.
- `sdk/src/phase-runner.ts` — full phase loop.
- `sdk/src/planning-journal.ts` — append-only audit log of phase transitions.
- `agents/gsd-planner.md`, `agents/gsd-verifier.md`, `agents/gsd-executor.md`
  — the three agents that produce, check, and execute phases.

## How to read the upstream after you've finished the chapters

1. Start at `sdk/src/cli.ts`. Identify `parseArgs` — it's chapter 2.
2. Follow the `GSD` import to `sdk/src/index.ts`. Identify the call to the
   prompt builder — it's chapter 4.
3. Open `sdk/src/prompt-builder.ts`. The role of agent-loading there is
   chapter 1.
4. Trace the transport call — `sdk/src/cli-transport.ts` is chapter 6.
5. Find the spec/plan parser — `sdk/src/plan-parser.ts` is chapter 8.
6. Skim the workflow files under `get-shit-done/workflows/` — chapter 7.

You should be able to read every file with comprehension by the end of this
order.
