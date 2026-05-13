# Upstream reading — s08 · Spec & phases

Source: `gsd-build/get-shit-done` @ `ba625c0978e1133ead8a26546b5e6a435e886edb`.

## The canonical SPEC template: `get-shit-done/templates/spec.md`

```markdown
# Phase [X]: [Name] — Specification

**Created:** [date]
**Ambiguity score:** [score] (gate: ≤ 0.20)
**Requirements:** [N] locked

## Goal
[One precise sentence — specific and measurable. NOT "improve X" — instead "X changes from A to B".]

## Background
[Current state from codebase — what exists today, what's broken or missing, what triggers this work.]

## Requirements
1. **[Short label]**: [Specific, testable statement.]
   - Current: [what exists or does NOT exist today]
   - Target: [what it should become after this phase]
   - Acceptance: [concrete pass/fail check — how a verifier confirms this was met]

## Boundaries
**In scope:**
- [Explicit list of what this phase produces]

**Out of scope:**
- [Explicit list of what this phase does NOT do] — [brief reason why it's excluded]

## Acceptance Criteria
- [ ] [Pass/fail criterion — unambiguous, verifiable]
- [ ] [Pass/fail criterion]
```

Notice the discipline:

1. **Goal is one sentence**, specific, measurable. "Improve performance" is
   not a goal; "P95 latency drops from 800ms to ≤300ms" is.
2. **Requirements have Current/Target/Acceptance**. The acceptance line is
   a concrete pass/fail check that a verifier can run.
3. **Boundaries are explicit**. "In scope" and "Out of scope" lists pin the
   blast radius of one phase.
4. **Acceptance criteria are checkboxes**. Each resolves to PASS or FAIL.

Our chapter-8 spec format collapses this to `Goal:` / `Tasks:` / `Verify:` —
the minimum that demonstrates the mechanism. The conceptual machinery —
"each phase has a verifier that returns PASS or FAIL" — is intact.

## Phase runner: `sdk/src/phase-runner.ts`

The upstream phase-runner loops "plan → execute → verify" until the SPEC's
acceptance criteria are all green:

- on PASS, mark the phase complete, persist `STATE.md`, move on;
- on FAIL, route to `gap-closure` (a re-planning loop);
- both transitions log to `planning-journal.ts` for replay.

Our `ExecutePhase` + `NextIncomplete` is the simplest equivalent: one step,
one verifier, mark complete on success.

## What to read next in upstream

- `get-shit-done/templates/spec.md` — full template.
- `sdk/src/plan-parser.ts` — parser with dependency analysis.
- `sdk/src/phase-runner.ts` — full phase loop.
- `sdk/src/planning-journal.ts` — append-only audit log of phase transitions.
- `agents/gsd-verifier.md` — the agent that owns PASS/FAIL.
