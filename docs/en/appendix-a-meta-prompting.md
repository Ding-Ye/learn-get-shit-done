# Appendix A — Meta-prompting: the shape of a Claude Code framework

## What distinguishes a framework from a chat client

A chat client takes a user message, prepends a fixed system prompt, sends the
result to a model, and renders the response. The user composes the prompt;
the client just transports it.

A **meta-prompting framework** inverts the relationship: the *framework*
composes the prompt, often from several files. The user types a short
command — `/gsd plan`, `@gsd-executor do task 2`, `gsd:capture --note hi` —
and the framework expands that into a long, structured prompt with system
instructions, context fragments, allowed-tool declarations, a parsed
invocation, and the user's free-form text. The model is told not just *what
to do* but *how to be while doing it* and *what tools are available*.

The upstream `get-shit-done` project is one such framework. Stripping away
its (substantial) TypeScript runtime, packaging, and Claude Code event
plumbing, the conceptual shape is small:

- A **skill** is a markdown file with YAML frontmatter declaring metadata
  and a body that is prose the model will see.
- A **command** is one kind of skill — `/gsd <name>` — wired to a markdown
  file under `commands/gsd/`.
- An **agent** is another kind — `@agent-name` — wired to a markdown file
  under `agents/`. Agents add a `<role>` block.
- A **context** is a reusable prose fragment under `contexts/`, optionally
  with template variables.
- A **workflow** is a markdown file under `workflows/` that declares a chain
  of commands.
- A **spec** is a markdown file under `templates/` that decomposes a goal
  into phases.
- A **hook** is a function that transforms a prompt or response.
- A **transport** is the adapter that actually invokes the model.

## Anatomy of a skill file

```yaml
---
name: gsd:capture                              # how users invoke it
description: Capture ideas to their destination  # one-line summary
argument-hint: "[--note] [text]"               # tab-completion hint
allowed-tools: [Read, Write, Bash]             # what the model may use
---

<objective>...</objective>                      # the prompt body
<process>...</process>
```

Three observations:

1. The **frontmatter is for the framework**; the body is for the model. The
   split lets us validate, render, and route in code while still shipping
   the prose unchanged to the model.
2. The body can use **XML-style tags** as soft section headings. Claude
   handles them well, and the framework can grep them when needed
   (`<role>…</role>`, `<process>…</process>`).
3. The same file is **the documentation**. There is no second source of
   truth; you can `cat` a skill file and understand it.

## Anatomy of a workflow file

```markdown
1. id: append
   command: capture
   args: "--note {{.Text}}"
   next: list

2. id: list
   command: capture
   args: "--list"
```

The workflow file is **declarative**: every step is structured key/value
data. The runtime walks the list, substitutes variables, and dispatches to
the dispatcher. The user is not in the loop.

(The upstream workflow files are looser — mostly prose with embedded `<step>`
tags — because they're often interpreted by the model itself rather than a
strict parser. Our chapter-7 format is the same shape, hardened for code.)

## Anatomy of a spec file

```markdown
# Auth refactor — Specification

## Goal
Replace cookie auth with JWTs.

## Phase P01: Token issuance
Goal: Mint signed JWTs on login.
Tasks:
- add /token endpoint
- sign with HS256
Verify: hitting /token returns a JWT whose exp is ~12h away.
```

A spec is the framework's commitment device: **every requirement must be
falsifiable**. "Improve performance" is not a requirement. "P95 latency
drops from 800ms to ≤300ms" is. The verifier (chapter 8's callback) decides
PASS or FAIL and there is no third option.

## Composition vs concatenation

A naive system would concatenate everything into one giant string.
`get-shit-done` composes:

- The same **context fragment** can be reused by ten different commands.
- The same **agent persona** can be invoked by three different workflows.
- The same **spec** can be planned, executed, and verified by three
  different roles.

Composition demands typed values flowing between stages — exactly what the
`gsd/` package gives us. It's the difference between a Unix pipeline and a
single big shell script.

## When meta-prompting earns its complexity

Meta-prompting is overkill for one-off conversations. It earns its keep when:

- You repeat the same kind of task across many sessions.
- You want the same input to produce *byte-identical* prompts on every run
  (caching, replay, regression testing).
- You want to swap personas without rewriting prompts.
- You want hooks (redaction, logging) you can add once and reuse everywhere.
- You want a junior contributor to be productive without first internalising
  every prompt-engineering trick the team has discovered.

If you're nodding to most of those, meta-prompting is your tool. If not,
write a chat loop and move on.
