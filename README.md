# learn-get-shit-done

A hands-on, from-zero **Go** learning companion to the upstream
[gsd-build/get-shit-done](https://github.com/gsd-build/get-shit-done)
JavaScript/TypeScript project — a Claude Code meta-prompting / context-engineering
framework.

This repository walks you through building a small Claude Code-style
meta-prompting framework in Go, **one mechanism per chapter**. Each chapter is a
tiny, focused Go program (~80–150 LOC) that demonstrates exactly one idea from
the upstream system.

> Looking for the Chinese version? See [`README.zh.md`](./README.zh.md).

![ci](https://github.com/Ding-Ye/learn-get-shit-done/actions/workflows/ci.yml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/Ding-Ye/learn-get-shit-done.svg)](https://pkg.go.dev/github.com/Ding-Ye/learn-get-shit-done)
![Go 1.23+](https://img.shields.io/badge/go-1.23%2B-00ADD8)
![License: Apache 2.0](https://img.shields.io/badge/license-Apache--2.0-blue)

## Why this exists

The upstream project is large, polyglot, and TypeScript-heavy. If you've never
read a meta-prompting framework before, jumping into ~190K lines of code is
discouraging. This companion strips the system down to its bones, ports each
core mechanism into ~100 lines of idiomatic Go, and points you back at the
upstream source so you can read it with confidence afterwards.

You will learn:

- How a markdown file with YAML frontmatter becomes a **runtime Skill**.
- How a `/gsd <command>` slash invocation is **parsed**.
- How **context templates** compose with variables to produce prompt fragments.
- How a final prompt is **assembled** from command + context + agent + user input.
- How **hooks** transform prompts and responses on the way in and out.
- How **agents** are looked up and **dispatched** through a Transport interface.
- How a **workflow** declaratively chains commands.
- How a **spec** turns into phases that you can execute one at a time.

## Curriculum

| #   | Chapter            | What you build                                                   |
| --- | ------------------ | ---------------------------------------------------------------- |
| 01  | Skill loader       | Parse a markdown skill file with YAML frontmatter → `gsd.Skill`. |
| 02  | Command parser     | Parse `/gsd <name> [args]` → structured `gsd.Command`.           |
| 03  | Context engine     | Load context templates from disk + render with `text/template`.  |
| 04  | Prompt assembler   | Combine command + context + agent persona + user input.          |
| 05  | Hook runner        | Pre/post hooks: redact, truncate, log.                           |
| 06  | Agent dispatcher   | Resolve `@agent` → load → dispatch via an `EchoTransport`.       |
| 07  | Workflow runner    | Execute commands chained in a workflow markdown file.            |
| 08  | Spec-driven plan   | Parse a spec → planned phases → execute one phase at a time.     |

End-to-end integration walk: [`docs/en/s_full-integration.md`](./docs/en/s_full-integration.md).

Two appendices:

- [Appendix A — Meta-prompting: the shape of a Claude Code framework](./docs/en/appendix-a-meta-prompting.md)
- [Appendix B — Upstream map (where to read what)](./docs/en/appendix-b-upstream-map.md)

## How to use this repo

```sh
git clone https://github.com/Ding-Ye/learn-get-shit-done.git
cd learn-get-shit-done
go vet ./...
go test ./...

# Run any chapter as a tiny demo
go run ./agents/s01-skill-loader
go run ./agents/s02-command-parser '/gsd capture --note Pay the bill'
# ...
```

Each chapter directory contains:

- `main.go` — a runnable demo, intentionally small (~80–150 LOC).
- `*_test.go` — a unit test covering the main behavior.
- `README.md` — pointer to the long-form doc.
- `testdata/` — small fixture files modeled on the upstream's markdown shapes.

Each chapter also ships an English long-form doc under
`docs/en/sNN-<name>.md` and a Chinese mirror under `docs/zh/sNN-<name>.md`,
both organised in six sections:

1. **Problem** — what mechanism are we modelling, and why does it matter?
2. **Solution** — the smallest Go construct that demonstrates it.
3. **How It Works** — line-by-line walkthrough of the code.
4. **What Changed** — the delta from the previous chapter.
5. **Try It** — copy-pasteable shell commands.
6. **Upstream Source Reading** — pointers into the upstream tree.

Annotated upstream excerpts (each commented heavily) live under
[`upstream-readings/`](./upstream-readings).

## Project layout

```
.
├── go.mod
├── gsd/                       # shared domain types (Skill, Command, Agent, …)
├── agents/                    # one chapter per directory
│   ├── s01-skill-loader/
│   ├── s02-command-parser/
│   ├── s03-context-engine/
│   ├── s04-prompt-assembler/
│   ├── s05-hook-runner/
│   ├── s06-agent-dispatcher/
│   ├── s07-workflow-runner/
│   └── s08-spec-driven-planner/
├── docs/                      # en/ + zh/ long-form docs
├── upstream-readings/         # annotated upstream excerpts (.go.md)
├── web/                       # static doc viewer
└── .github/workflows/ci.yml
```

## License & upstream credit

- This repository's Go source and documentation are licensed under **Apache 2.0**.
- The upstream project [`gsd-build/get-shit-done`](https://github.com/gsd-build/get-shit-done)
  (commit `ba625c0978e1133ead8a26546b5e6a435e886edb`) is distributed under the
  **MIT License**. We do not redistribute its source; we only quote small,
  annotated excerpts under fair use for educational purposes.

See [`LICENSE`](./LICENSE) for the full text.
