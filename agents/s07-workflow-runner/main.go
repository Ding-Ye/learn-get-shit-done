// Chapter 07 — Workflow runner.
//
// A workflow chains commands declaratively. The upstream `workflows/*.md`
// files mix narrative prose with executable steps; our simplified format is
// a YAML-frontmatter workflow definition plus a numbered list of steps in
// the markdown body. The runner reads the file, resolves each step's command
// to a (pretend) executor, runs them in order, and returns a log.
//
// Workflow file shape (see testdata/note.md):
//
//   ---
//   name: note
//   description: Append a note then list them
//   ---
//
//   1. id: append
//      command: capture
//      args: "--note {{.Text}}"
//      next: list
//
//   2. id: list
//      command: capture
//      args: "--list"
package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Ding-Ye/learn-get-shit-done/gsd"
)

// CommandRunner is the dependency the workflow runner needs: a way to execute
// one command by name with an args string and return its output.
type CommandRunner func(cmd, args string) (string, error)

// LoadWorkflow reads a workflow markdown file at path and returns the parsed Workflow.
func LoadWorkflow(path string) (gsd.Workflow, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return gsd.Workflow{}, err
	}
	return parseWorkflow(string(raw))
}

// Run executes the workflow steps in order, starting from the first step.
// Each step's command is dispatched to runCmd. The returned trace is a
// concatenated log of each step's input, output, and any transition.
//
// Variables: each `args` field may contain `{{.Key}}` placeholders that are
// substituted from `vars` before dispatch.
func Run(w gsd.Workflow, vars map[string]string, runCmd CommandRunner) (string, error) {
	if len(w.Steps) == 0 {
		return "", errors.New("workflow has no steps")
	}
	byID := map[string]gsd.WorkflowStep{}
	for _, s := range w.Steps {
		byID[s.ID] = s
	}

	var trace strings.Builder
	current := w.Steps[0].ID
	visited := 0
	for current != "" {
		step, ok := byID[current]
		if !ok {
			return trace.String(), fmt.Errorf("unknown step id %q", current)
		}
		args := substitute(step.Args, vars)
		fmt.Fprintf(&trace, "[step %s] command=%s args=%q\n", step.ID, step.Command, args)

		out, err := runCmd(step.Command, args)
		if err != nil {
			if step.OnError != "" {
				fmt.Fprintf(&trace, "  -> error: %v; routing to %q\n", err, step.OnError)
				current = step.OnError
				continue
			}
			return trace.String(), fmt.Errorf("step %s: %w", step.ID, err)
		}
		fmt.Fprintf(&trace, "  -> ok: %s\n", strings.TrimSpace(out))

		current = step.Next

		visited++
		if visited > 100 {
			return trace.String(), errors.New("workflow exceeded 100 steps (loop?)")
		}
	}
	return trace.String(), nil
}

// substitute does cheap {{.Key}} replacement — we don't bring in text/template
// here because workflow args are typically short single-line strings.
func substitute(s string, vars map[string]string) string {
	for k, v := range vars {
		s = strings.ReplaceAll(s, "{{."+k+"}}", v)
	}
	return s
}

// parseWorkflow reads our simplified workflow markdown. Frontmatter gives the
// workflow's name + description. The body is a numbered list of steps where
// each step is a small key/value block.
func parseWorkflow(text string) (gsd.Workflow, error) {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	lines := strings.Split(text, "\n")
	w := gsd.Workflow{}
	body := lines

	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "---" {
		closing := -1
		for i := 1; i < len(lines); i++ {
			if strings.TrimSpace(lines[i]) == "---" {
				closing = i
				break
			}
		}
		if closing == -1 {
			return w, errors.New("unterminated frontmatter")
		}
		for _, line := range lines[1:closing] {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			k, v, ok := strings.Cut(line, ":")
			if !ok {
				continue
			}
			k = strings.TrimSpace(k)
			v = strings.Trim(strings.TrimSpace(v), `"'`)
			switch k {
			case "name":
				w.Name = v
			case "description":
				w.Description = v
			}
		}
		body = lines[closing+1:]
	}

	var current *gsd.WorkflowStep
	flush := func() {
		if current != nil {
			w.Steps = append(w.Steps, *current)
		}
		current = nil
	}
	for _, raw := range body {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		// A numbered marker like `1.` or `1)` starts a new step.
		if isStepHeader(line) {
			flush()
			current = &gsd.WorkflowStep{}
			rest := strings.TrimSpace(stripStepHeader(line))
			if rest != "" {
				applyStepLine(current, rest)
			}
			continue
		}
		if current == nil {
			continue // ignore prose before the first step
		}
		applyStepLine(current, line)
	}
	flush()
	if w.Name == "" {
		return w, errors.New("workflow has no name in frontmatter")
	}
	return w, nil
}

func isStepHeader(line string) bool {
	// `1.` / `1)` / `12.` etc.
	for i, r := range line {
		if r >= '0' && r <= '9' {
			continue
		}
		if i == 0 {
			return false
		}
		return r == '.' || r == ')'
	}
	return false
}

func stripStepHeader(line string) string {
	for i, r := range line {
		if r == '.' || r == ')' {
			return line[i+1:]
		}
	}
	return line
}

func applyStepLine(s *gsd.WorkflowStep, line string) {
	k, v, ok := strings.Cut(line, ":")
	if !ok {
		return
	}
	k = strings.TrimSpace(k)
	v = strings.Trim(strings.TrimSpace(v), `"'`)
	switch k {
	case "id":
		s.ID = v
	case "command":
		s.Command = v
	case "args":
		s.Args = v
	case "next":
		s.Next = v
	case "on_error", "on-error":
		s.OnError = v
	case "describe":
		s.Describe = v
	}
}

// ─── Demo executor & main ───────────────────────────────────────────────────

func main() {
	path := filepath.Join("agents", "s07-workflow-runner", "testdata", "note.md")
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	w, err := LoadWorkflow(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Demo CommandRunner: pretend each command is a function that just echoes.
	exec := func(cmd, args string) (string, error) {
		return fmt.Sprintf("(%s ran with %q)", cmd, args), nil
	}
	out, err := Run(w, map[string]string{"Text": "Pay the bill"}, exec)
	fmt.Print(out)
	if err != nil {
		fmt.Fprintln(os.Stderr, "workflow error:", err)
		os.Exit(1)
	}
}
