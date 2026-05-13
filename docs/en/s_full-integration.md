# Full integration — `/gsd plan "fix the auth bug"` end to end

## Problem

We have eight chapters, each demonstrating one mechanism in isolation. The
exam question is: when a user types a single command, how do those mechanisms
compose? This walk traces a single hypothetical invocation —
`/gsd plan "fix the auth bug"` — through every stage. No new code; this is a
reading guide stitched from the chapters you've already built.

## Solution

Follow the data flow:

```
user input string
   │
   ▼  chapter 2: ParseCommand
gsd.Command{Name:"plan", Flags:{}, Args:["fix","the","auth","bug"]}
   │
   ▼  chapter 1: LoadSkill("commands/gsd/plan.md")
gsd.Skill{Name:"gsd:plan", AllowedTools:[...], Body:"..."}
   │
   ▼  chapter 3: engine.Render("dev", {User:"Ding"})
context string ("Output style: concise…")
   │
   ▼  chapter 6: dispatcher.Resolve("@gsd-planner")
gsd.Agent{Name:"gsd-planner", Role:"You produce PLAN.md..."}
   │
   ▼  chapter 4: Assemble(Inputs{skill, command, agent, context, user})
final prompt string
   │
   ▼  chapter 5: runner.Run(HookPre, prompt) — redact, truncate
sanitised prompt string
   │
   ▼  chapter 6: transport.Send(prompt)
gsd.TransportResult{Response:"..."}
   │
   ▼  chapter 5: runner.Run(HookPost, response) — log
final response shown to user
```

## How It Works

Step by step, with the chapter that introduced the piece:

1. **Parse the user input.** `agents/s02-command-parser` turns the raw line
   into `gsd.Command{Name: "plan", Args: ["fix", "the", "auth", "bug"]}`.
   Quoting and flag detection happen here so no later stage needs to think
   about strings.

2. **Load the matching skill.** The `Name` field of the command tells us to
   load `commands/gsd/plan.md` (we'd configure a skills directory). The
   chapter-1 `LoadSkill` returns a `gsd.Skill` with `AllowedTools` and `Body`.

3. **Decide which context to render.** The skill body, or a config file,
   chooses `"dev"` or `"review"`. We pass it through the chapter-3
   `engine.Render("dev", vars)`, which loads `contexts/dev.md`, merges
   variables (caller wins over frontmatter defaults), and executes a
   `text/template`. Output is a ready-to-paste string.

4. **Resolve the agent persona.** The skill may reference `@gsd-planner`
   (workflows often do). The chapter-6 dispatcher's `Resolve("@gsd-planner")`
   reads `agents/gsd-planner.md` and returns a `gsd.Agent` with `Tools`,
   `Description`, and the `<role>` block.

5. **Assemble the final prompt.** Chapter 4's `Assemble(Inputs{...})` joins
   SYSTEM (agent), CONTEXT (rendered), COMMAND (skill body + parsed
   invocation), USER (the original input). The output is the deterministic
   string a transport would send.

6. **Run pre-hooks.** Chapter 5's `Runner.Run(HookPre, prompt)` walks the
   registered hooks: `Redact()` strips API keys, `Truncate(N)` keeps it under
   budget, optional inline hooks tag the prompt for tracing. The output is
   the bytes that actually leave the host.

7. **Dispatch.** Chapter 6's `Transport.Send(prompt)` returns a
   `gsd.TransportResult`. In tests this is `EchoTransport`; in production it
   is the real Claude Code transport.

8. **Run post-hooks.** Chapter 5's `Runner.Run(HookPost, response)` is the
   symmetric step on the way back: log, redact, normalize.

If the user instead typed `/gsd workflow note "Pay the bill"`, chapter 7's
`Run(workflow, vars, dispatcher.Dispatch)` would replace step 7, chaining
several commands. And if the user typed `/gsd plan-phase auth-refactor`,
chapter 8's `LoadSpec` + `ExecutePhase` would run before any of the rest:
first carve the user goal into phases, then for each phase invoke the same
pipeline above.

## What Changed

This document doesn't change the code. It's the **map** that lets you
navigate the eight chapters in execution order rather than in the order I
introduced them. If you can recite this trace, you've internalised the
framework.

## Try It

You can hand-stitch the trace in a single short Go program (left as an
exercise). Or read each chapter's `main.go` in the order listed above —
each one prints its piece of the trace, so reading them sequentially
recreates the integration in your head.

```sh
go run ./agents/s02-command-parser '/gsd plan fix the auth bug'
go run ./agents/s01-skill-loader  agents/s04-prompt-assembler/testdata/plan.md
go run ./agents/s03-context-engine dev Ding
go run ./agents/s06-agent-dispatcher @gsd-planner "fix the auth bug"
```

## Upstream Source Reading

- `sdk/src/cli.ts` — the actual end-to-end entry point. Read it last; its
  shape will now feel familiar because you wrote each piece yourself.
- `sdk/src/index.ts` — top-level `GSD` class that wires everything together.
- `sdk/src/session-runner.ts` — long-lived sessions that drive workflows.

Annotated excerpts: each upstream-reading file in
[`upstream-readings/`](../../upstream-readings) corresponds to one chapter.
