# Upstream reading — s02 · CLI argument parsing

Source: `gsd-build/get-shit-done` @ `ba625c0978e1133ead8a26546b5e6a435e886edb`.

## Upstream entry point: `sdk/src/cli.ts`

```ts
export interface ParsedCliArgs {
  command: string | undefined;
  prompt: string | undefined;
  initInput: string | undefined;
  init: string | undefined;
  projectDir: string;
  wsPort: number | undefined;
  model: string | undefined;
  maxBudget: number | undefined;
  ws: string | undefined;
  help: boolean;
  version: boolean;
  queryArgv?: string[];
}
```

Two things to notice:

1. The parsed value is a **plain interface** with optional fields — exactly
   the shape our `gsd.Command` is, just with a different vocabulary.
2. The upstream also keeps `queryArgv` — a raw passthrough of tokens it didn't
   know how to interpret. This is the same forward-compat trick we use by
   collecting unknown positional tokens into `cmd.Args`.

## How `$ARGUMENTS` ends up here

Inside `commands/gsd/capture.md` the body contains:

```
<context>
Arguments: $ARGUMENTS

Parse the first token of $ARGUMENTS:
- If it is `--note`: strip the flag, pass remainder to note workflow
- If it is `--backlog`: strip the flag, pass remainder to add-backlog workflow
- If it is `--seed`: strip the flag, pass remainder to plant-seed workflow
- If it is `--list`: pass remainder (optional area filter) to check-todos workflow
- Otherwise: pass all of $ARGUMENTS to add-todo workflow
</context>
```

The framework substitutes `$ARGUMENTS` with whatever the user typed after the
command name. Our `gsd.Command.Flags` + `gsd.Command.Args` together represent
exactly the same data, but in a structured way — so the dispatcher in
chapter 6 doesn't have to parse strings at runtime.

## What to read next in upstream

- `sdk/src/cli.ts` — top-level CLI entry; main flag dispatch.
- `sdk/src/init-runner.ts` — special-case parsing for `init` (reads `@file`).
- `sdk/src/query/` — the `query` subcommand keeps unknown flags verbatim.
