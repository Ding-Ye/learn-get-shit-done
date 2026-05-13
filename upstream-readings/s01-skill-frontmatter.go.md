# Upstream reading — s01 · Skill frontmatter

Source: `gsd-build/get-shit-done` @ `ba625c0978e1133ead8a26546b5e6a435e886edb`
(MIT). Excerpts quoted under fair use for educational purposes.

## A real skill file: `commands/gsd/capture.md`

```markdown
---
name: gsd:capture
description: Capture ideas, tasks, notes, and seeds to their destination
argument-hint: "[--note | --backlog | --seed | --list] [text]"
allowed-tools:
  - Read
  - Write
  - Edit
  - Bash
  - Glob
  - Grep
  - AskUserQuestion
---

<objective>
Capture ideas, tasks, notes, and seeds to their appropriate destination in the GSD system.
...
</objective>
```

Why this shape works for a meta-prompting framework:

- The YAML frontmatter gives the **framework** the structured metadata it
  needs to render menus, validate tool permissions, and decide which body to
  insert when the model invokes the skill.
- The markdown body is the **prompt** the model will see — it can use any
  tags the framework recognises (`<objective>`, `<process>`, `<routing>`).
- The two halves are physically inseparable, so the description and the
  prompt cannot drift apart.

## How upstream parses the same kind of file

From `sdk/src/prompt-builder.ts` (paraphrased — we only quote the regex shape):

```ts
// Look for "tools:" in the YAML frontmatter
const frontmatterMatch = agentDef.match(/^---\s*\n([\s\S]*?)\n---/);
if (!frontmatterMatch) return DEFAULT_ALLOWED_TOOLS;

const toolsMatch = frontmatterMatch[1].match(/^tools:\s*(.+)$/m);
```

This is exactly the trick we use in Go: pin the frontmatter between two
`---` fences and parse only the keys we care about. We use `strings.Cut`
where the upstream uses a regex; the result is identical.

Note that upstream's `agents/*.md` files use the key `tools:` while
`commands/gsd/*.md` files use `allowed-tools:`. Our `gsd.Skill` collapses
both into `AllowedTools` for simplicity; an `gsd.Agent` (chapter 6) will
also reuse the same parser.

## What to read next in upstream

- `commands/gsd/*.md` — every command is in this shape. Skim a few.
- `agents/gsd-planner.md`, `agents/gsd-executor.md` — agents have a slightly
  richer frontmatter (`color`, optional `hooks`).
- `sdk/src/prompt-builder.ts` — `parseAgentTools`, `parseAgentRole`.
- `sdk/src/init-runner.ts` — where skills are discovered at startup.
