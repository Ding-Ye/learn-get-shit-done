# Upstream reading — s03 · Context engine

Source: `gsd-build/get-shit-done` @ `ba625c0978e1133ead8a26546b5e6a435e886edb`.

## A real context fragment: `get-shit-done/contexts/dev.md`

```markdown
# Dev Context Profile

Agent output guidance for dev mode. Loaded when `context: dev` is set in config.json.

## Output Style

- Concise, action-oriented responses
- Lead with the code change or command, follow with brief rationale
- Skip preamble — assume the developer has full context
- Use inline code references (`file:line`) over prose descriptions
```

This fragment has **no template variables** — it's pure prose. The upstream
relies on the fact that the prompt-assembler inlines this file as-is.

The interesting move in the upstream isn't "render template variables" — it's
"select the right files for the right phase". From `sdk/src/context-engine.ts`:

```ts
const PHASE_FILE_MANIFEST: Record<PhaseType, FileSpec[]> = {
  [PhaseType.Execute]: [
    { key: 'state', filename: 'STATE.md', required: true },
    { key: 'config', filename: 'config.json', required: false },
  ],
  [PhaseType.Research]: [
    { key: 'state', filename: 'STATE.md', required: true },
    { key: 'roadmap', filename: 'ROADMAP.md', required: true },
    { key: 'context', filename: 'CONTEXT.md', required: true },
    { key: 'requirements', filename: 'REQUIREMENTS.md', required: false },
  ],
  // …
};
```

Each phase has its own file manifest. The engine reads only the files that
matter for that phase and concatenates them into the prompt. Truncation
happens on top: large files have their headings + first paragraph kept; the
rest is dropped to keep the prompt cache-friendly.

Our chapter-3 engine is a **strict subset**: one fragment per call, with
template-variable rendering. We earn the conceptual right to talk about
"a context" before we add manifests and truncation in later reading.

## What to read next in upstream

- `sdk/src/context-engine.ts` — the full engine.
- `sdk/src/context-truncation.ts` — heading-aware truncation.
- `get-shit-done/contexts/` — the three actual context profiles.
