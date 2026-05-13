# Upstream reading — s05 · Hooks

Source: `gsd-build/get-shit-done` @ `ba625c0978e1133ead8a26546b5e6a435e886edb`.

## Agent-declared hooks: `agents/gsd-planner.md`

```yaml
---
name: gsd-planner
description: Creates executable phase plans …
tools: Read, Write, Bash, Glob, Grep, WebFetch, mcp__context7__*
color: green
# hooks:
#   PostToolUse:
#     - matcher: "Write|Edit"
#       hooks:
#         - type: command
#           command: "npx eslint --fix $FILE 2>/dev/null || true"
---
```

The block is commented out in the upstream snapshot, but it shows the design:
**agents declare hooks per tool-use event** (matched by a regex on the tool
name) and the framework executes them as shell commands. Our Go version is
deliberately simpler — just pre/post on the whole payload — but the lifecycle
shape is the same:

```
user → assemble → [pre-hooks] → transport → response → [post-hooks] → user
```

## What kinds of transforms hooks do

Skimming `sdk/src/prompt-sanitizer.ts` (and the test counterpart) tells you
the common categories:

1. **Sanitize** — strip control characters, normalize line endings.
2. **Budget** — truncate or summarize to fit a token cap (our `Truncate`).
3. **Redact** — remove secrets, PII, or paths that shouldn't leave the host
   (our `Redact`).
4. **Observe** — append the payload to a log file for replay (our `Log`).
5. **Reject** — return an error to abort the dispatch entirely (we model
   this via `HookFunc` returning a non-nil error).

The pattern is consistent: take a string, return a (maybe transformed)
string-or-error.

## What to read next in upstream

- `hooks/lib/git-cmd.js` — example hook implementation.
- `sdk/src/prompt-sanitizer.ts` — pre-send scrubbing.
- `sdk/src/logger.ts` — observability hook stack.
