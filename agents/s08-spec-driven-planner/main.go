// Chapter 08 — Spec-driven planner.
//
// Take a SPEC markdown file, parse it into a Spec (title + goal + phases),
// then execute a single phase: print its tasks, run a verification stub, and
// flip its Complete flag. The upstream's "phase" idea is a way of breaking a
// large goal into testable, falsifiable chunks. Our chapter ports the
// parsing + one-phase-at-a-time execution discipline.
//
// SPEC format (see testdata/auth-refactor.md):
//
//   # Auth refactor — Specification
//
//   ## Goal
//   Replace cookie auth with JWT.
//
//   ## Phase P01: Token issuance
//   Goal: mint tokens on login.
//   Tasks:
//   - add /token endpoint
//   - sign with HS256
//   Verify: hitting /token returns a valid JWT.
//
//   ## Phase P02: Token validation
//   Goal: validate tokens in middleware.
//   ...
package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Ding-Ye/learn-get-shit-done/gsd"
)

// LoadSpec reads a spec markdown file and returns the parsed Spec.
func LoadSpec(path string) (gsd.Spec, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return gsd.Spec{}, err
	}
	return parseSpec(string(raw))
}

// ExecutePhase runs the verification step for a single phase, returns the
// updated Spec with that phase marked Complete on success, and a human log.
//
// `verifier` is the testable seam — the runner doesn't itself decide whether
// a phase passes. Real callers would invoke `go test` or shell out to a
// build, then report back through this callback.
func ExecutePhase(s gsd.Spec, phaseID string, verifier func(p gsd.SpecPhase) error) (gsd.Spec, string, error) {
	var log strings.Builder
	for i, p := range s.Phases {
		if p.ID != phaseID {
			continue
		}
		fmt.Fprintf(&log, "Executing phase %s — %s\n", p.ID, p.Title)
		fmt.Fprintf(&log, "Goal: %s\n", p.Goal)
		fmt.Fprintf(&log, "Tasks (%d):\n", len(p.Tasks))
		for j, t := range p.Tasks {
			fmt.Fprintf(&log, "  %d. %s\n", j+1, t)
		}
		fmt.Fprintf(&log, "Verify: %s\n", p.Verify)
		if err := verifier(p); err != nil {
			return s, log.String(), fmt.Errorf("phase %s verification failed: %w", p.ID, err)
		}
		s.Phases[i].Complete = true
		fmt.Fprintf(&log, "Phase %s complete.\n", p.ID)
		return s, log.String(), nil
	}
	return s, log.String(), fmt.Errorf("phase %q not found", phaseID)
}

// NextIncomplete returns the first phase whose Complete is false, or false if all are done.
func NextIncomplete(s gsd.Spec) (gsd.SpecPhase, bool) {
	for _, p := range s.Phases {
		if !p.Complete {
			return p, true
		}
	}
	return gsd.SpecPhase{}, false
}

// ─── Parsing ───────────────────────────────────────────────────────────────

var phaseHeading = regexp.MustCompile(`^##\s*Phase\s+(\S+):\s*(.+?)\s*$`)

func parseSpec(text string) (gsd.Spec, error) {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	lines := strings.Split(text, "\n")
	spec := gsd.Spec{}

	// Title — first `# ` line.
	for _, line := range lines {
		if strings.HasPrefix(line, "# ") {
			spec.Title = strings.TrimSpace(strings.TrimPrefix(line, "# "))
			break
		}
	}
	if spec.Title == "" {
		return spec, errors.New("missing # Title line")
	}

	// We walk sections. State: either "outer", "in goal", or "in phase".
	const (
		sOuter = iota
		sGoal
		sPhase
		sTasks
	)
	state := sOuter
	var current *gsd.SpecPhase
	var goalLines []string

	flushPhase := func() {
		if current != nil {
			spec.Phases = append(spec.Phases, *current)
			current = nil
		}
	}

	for _, line := range lines {
		trimmed := strings.TrimRight(line, " \t")

		// Phase heading takes precedence regardless of state.
		if m := phaseHeading.FindStringSubmatch(trimmed); m != nil {
			flushPhase()
			current = &gsd.SpecPhase{ID: m[1], Title: m[2]}
			state = sPhase
			continue
		}

		// Top-level Goal section.
		if strings.TrimSpace(trimmed) == "## Goal" {
			flushPhase()
			state = sGoal
			continue
		}

		// Any other H2 closes both Goal and Phase contexts.
		if strings.HasPrefix(trimmed, "## ") {
			flushPhase()
			state = sOuter
			continue
		}

		switch state {
		case sGoal:
			if strings.TrimSpace(trimmed) != "" {
				goalLines = append(goalLines, strings.TrimSpace(trimmed))
			}
		case sPhase, sTasks:
			if current == nil {
				continue
			}
			low := strings.TrimSpace(trimmed)
			switch {
			case strings.HasPrefix(low, "Goal:"):
				current.Goal = strings.TrimSpace(strings.TrimPrefix(low, "Goal:"))
				state = sPhase
			case strings.HasPrefix(low, "Verify:"):
				current.Verify = strings.TrimSpace(strings.TrimPrefix(low, "Verify:"))
				state = sPhase
			case low == "Tasks:":
				state = sTasks
			case state == sTasks && strings.HasPrefix(low, "- "):
				current.Tasks = append(current.Tasks, strings.TrimPrefix(low, "- "))
			}
		}
	}
	flushPhase()
	spec.Goal = strings.Join(goalLines, " ")
	if len(spec.Phases) == 0 {
		return spec, errors.New("spec has no phases")
	}
	return spec, nil
}

// ─── Demo ──────────────────────────────────────────────────────────────────

func main() {
	path := filepath.Join("agents", "s08-spec-driven-planner", "testdata", "auth-refactor.md")
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	spec, err := LoadSpec(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("Spec: %s\nGoal: %s\nPhases: %d\n\n", spec.Title, spec.Goal, len(spec.Phases))

	next, ok := NextIncomplete(spec)
	if !ok {
		fmt.Println("All phases already complete.")
		return
	}
	updated, log, err := ExecutePhase(spec, next.ID, func(p gsd.SpecPhase) error {
		// Stub verifier: succeeds if Verify is non-empty.
		if strings.TrimSpace(p.Verify) == "" {
			return errors.New("no Verify statement")
		}
		return nil
	})
	fmt.Print(log)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}
	if next, ok := NextIncomplete(updated); ok {
		fmt.Printf("\nNext phase: %s — %s\n", next.ID, next.Title)
	} else {
		fmt.Println("\nAll phases complete.")
	}
}
