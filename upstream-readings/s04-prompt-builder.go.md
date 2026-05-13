# Upstream reading — s04 · Prompt builder

Source: `gsd-build/get-shit-done` @ `ba625c0978e1133ead8a26546b5e6a435e886edb`.

## Upstream task formatter: `sdk/src/prompt-builder.ts`

```ts
function formatTask(task: PlanTask, index: number): string {
  const lines: string[] = [];
  lines.push(`### Task ${index + 1}: ${task.name}`);

  if (task.files.length > 0) {
    lines.push(`**Files:** ${task.files.join(', ')}`);
  }

  if (task.read_first.length > 0) {
    lines.push(`**Read first:** ${task.read_first.join(', ')}`);
  }

  lines.push('');
  lines.push('**Action:**');
  lines.push(task.action);

  if (task.verify) {
    lines.push('');
    lines.push('**Verify:**');
    lines.push(task.verify);
  }
  // … done, acceptance_criteria
  return lines.join('\n');
}
```

Two things to notice:

1. The upstream builder is also a **string-builder pattern** — same shape
   as our `Assemble` function, just for tasks instead of full prompts.
2. The fields it includes (`Files`, `Read first`, `Action`, `Verify`,
   `Done when`, `Acceptance criteria`) come straight from the SPEC/PLAN
   markdown templates. The structure of the markdown becomes the structure
   of the prompt.

## Agent role extraction

```ts
export function parseAgentRole(agentDef: string): string {
  const match = agentDef.match(/<role>([\s\S]*?)<\/role>/i);
  return match ? match[1].trim() : '';
}
```

Our `loadAgent` does the same thing with `strings.Index` on `<role>` and
`</role>`. We don't need a regex for it.

## Why determinism matters

The upstream ships `sdk/src/assembled-prompts.test.ts` with **golden-file**
assertions: parse a fixture, build a prompt, compare bytes-for-bytes against
a checked-in expected output. If our assembler isn't deterministic, that test
flakes; if a future refactor accidentally re-orders sections, the prompt
caching that depends on identical prefixes silently breaks.

Our chapter-4 `TestAssemble_deterministicOrder` is the tiny version of that
discipline.

## What to read next in upstream

- `sdk/src/prompt-builder.ts` — full prompt builder.
- `sdk/src/assembled-prompts.test.ts` — golden tests.
- `sdk/src/phase-prompt.ts` — phase-prompt layer.
- `sdk/src/prompt-sanitizer.ts` — final scrubbing pass before send (relevant
  in chapter 5 too).
