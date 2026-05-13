// Package gsd holds the shared domain types used across every chapter.
//
// The types here intentionally mirror the conceptual vocabulary of the upstream
// gsd-build/get-shit-done framework (Skill, Command, Context, Agent, Hook,
// Workflow, Phase, Transport). Each chapter under agents/sNN-* imports this
// package and demonstrates ONE mechanism that operates on these types.
//
// Keep this file behavior-free. Concrete parsing, rendering, and dispatch logic
// lives in the chapter that introduces it.
package gsd

import "time"

// ─── Skill ──────────────────────────────────────────────────────────────────

// Skill is the runtime representation of a markdown skill file with YAML
// frontmatter (upstream: commands/gsd/*.md, agents/*.md).
//
// The frontmatter populates the structured fields; everything after the
// closing `---` line is preserved verbatim in Body.
type Skill struct {
	Name         string   // e.g. "gsd:capture"
	Description  string   // one-line summary surfaced to the user
	ArgumentHint string   // e.g. "[--note | --backlog | --seed] [text]"
	AllowedTools []string // tool names the skill is permitted to use
	Body         string   // markdown body (after frontmatter)
}

// ─── Command ────────────────────────────────────────────────────────────────

// Command is a parsed slash-command invocation, e.g. `/gsd capture --note hi`.
//
// Name is the bare command name without the leading slash or namespace prefix.
// Flags collects long-form `--flag` and `--flag value` pairs. Positional words
// that are not flags land in Args in order.
type Command struct {
	Raw   string            // original input line
	Name  string            // e.g. "capture"
	Flags map[string]string // e.g. {"note": ""} or {"text": "hello"}
	Args  []string          // positional words after flags are consumed
}

// ─── Context ────────────────────────────────────────────────────────────────

// Context is a reusable prompt fragment loaded by name (upstream:
// get-shit-done/contexts/*.md). Variables are merged into Body via text/template.
type Context struct {
	Name      string            // file stem, e.g. "dev"
	Body      string            // template body
	Variables map[string]string // variables resolved at render time
}

// ─── Agent ──────────────────────────────────────────────────────────────────

// Agent is a markdown-defined sub-Claude persona (upstream: agents/*.md).
// Role is the prose extracted from the <role>…</role> block in the body, if any.
type Agent struct {
	Name        string   // e.g. "gsd-planner"
	Description string
	Tools       []string
	Role        string // optional <role>…</role> contents
	Body        string // full markdown body for reference
}

// ─── Hook ───────────────────────────────────────────────────────────────────

// HookPhase distinguishes pre-prompt and post-response hooks.
type HookPhase string

const (
	HookPre  HookPhase = "pre"  // runs on the assembled prompt before dispatch
	HookPost HookPhase = "post" // runs on the model's response before return
)

// HookFunc transforms a payload (prompt or response). Returning an error stops
// the chain and surfaces the failure to the caller.
type HookFunc func(payload string) (string, error)

// Hook bundles a name with a HookFunc for a given phase.
type Hook struct {
	Name  string
	Phase HookPhase
	Fn    HookFunc
}

// ─── Workflow ───────────────────────────────────────────────────────────────

// WorkflowStep is a single step inside a workflow: invoke a command and
// optionally transition to another step based on the response.
type WorkflowStep struct {
	ID       string // step identifier (used by Next references)
	Command  string // command name to invoke for this step
	Args     string // raw argument string passed to the command
	Next     string // ID of the next step, or empty for terminal
	OnError  string // ID of the step to jump to on error, or empty
	Describe string // human-readable label
}

// Workflow chains steps declaratively (upstream: get-shit-done/workflows/*.md).
type Workflow struct {
	Name        string
	Description string
	Steps       []WorkflowStep
}

// ─── Spec / Phases ──────────────────────────────────────────────────────────

// SpecPhase is one phase inside a spec-driven plan.
type SpecPhase struct {
	ID       string   // e.g. "P01"
	Title    string   // short human title
	Goal     string   // one-line objective
	Tasks    []string // ordered task list
	Verify   string   // verification statement
	Complete bool     // mutated by the runner after a phase finishes
}

// Spec is the parsed form of a spec markdown document.
type Spec struct {
	Title  string
	Goal   string
	Phases []SpecPhase
}

// ─── Transport ──────────────────────────────────────────────────────────────

// TransportResult bundles the response text and bookkeeping fields returned by
// a Transport implementation. We keep it tiny on purpose — the focus is on
// prompt assembly, not on faithfully modeling the upstream event stream.
type TransportResult struct {
	Response  string
	StartedAt time.Time
	EndedAt   time.Time
}

// Transport sends an assembled prompt and returns a response. Chapter 6 ships
// a simple EchoTransport implementation; real systems would shell out to
// `claude` or POST to an API.
type Transport interface {
	Send(prompt string) (TransportResult, error)
}
