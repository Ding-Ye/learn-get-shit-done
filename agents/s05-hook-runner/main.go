// Chapter 05 — Hook runner.
//
// Hooks are the framework's escape hatch: small, named functions that run
// before a prompt is sent (HookPre) and after the response is received
// (HookPost). Each hook gets to inspect, transform, or reject the payload.
//
// Real upstream uses cases:
//   - redact secrets before sending;
//   - truncate context to fit a token budget;
//   - log every prompt/response pair for replay;
//   - sanitize markdown (strip backticks, normalize whitespace).
//
// We ship three example hooks (Redact, Truncate, Log) and a Runner that
// chains them in order. The Runner returns the final payload or the first
// hook error.
package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/Ding-Ye/learn-get-shit-done/gsd"
)

// Runner holds an ordered list of hooks. Hooks for different phases
// (HookPre vs HookPost) are run separately by Run.
type Runner struct {
	Hooks []gsd.Hook
}

// New returns a Runner pre-populated with no hooks.
func New() *Runner { return &Runner{} }

// Add appends a hook in registration order.
func (r *Runner) Add(h gsd.Hook) { r.Hooks = append(r.Hooks, h) }

// Run executes every hook of the given phase against `payload` in registration
// order. The output of one hook becomes the input of the next. If any hook
// returns an error, Run stops and surfaces it.
func (r *Runner) Run(phase gsd.HookPhase, payload string) (string, error) {
	for _, h := range r.Hooks {
		if h.Phase != phase {
			continue
		}
		out, err := h.Fn(payload)
		if err != nil {
			return payload, fmt.Errorf("hook %s: %w", h.Name, err)
		}
		payload = out
	}
	return payload, nil
}

// ─── Sample hooks ────────────────────────────────────────────────────────────

var secretRegex = regexp.MustCompile(`(?i)(api[_-]?key|password|secret|token)[\s:=]+\S+`)

// Redact replaces obvious secrets with `[REDACTED]`. Intentionally pessimistic.
func Redact() gsd.Hook {
	return gsd.Hook{
		Name:  "redact",
		Phase: gsd.HookPre,
		Fn: func(s string) (string, error) {
			return secretRegex.ReplaceAllString(s, "[REDACTED]"), nil
		},
	}
}

// Truncate caps payload size at max runes, appending an ellipsis marker.
func Truncate(max int) gsd.Hook {
	return gsd.Hook{
		Name:  "truncate",
		Phase: gsd.HookPre,
		Fn: func(s string) (string, error) {
			r := []rune(s)
			if len(r) <= max {
				return s, nil
			}
			return string(r[:max]) + "\n…[truncated]", nil
		},
	}
}

// Log writes payload to w (typically os.Stderr) without modifying it.
// Returned as a post-hook so the response gets the spotlight; flip the phase
// to log prompts instead.
func Log(w io.Writer, label string) gsd.Hook {
	return gsd.Hook{
		Name:  "log:" + label,
		Phase: gsd.HookPost,
		Fn: func(s string) (string, error) {
			fmt.Fprintf(w, "[%s] %d bytes\n", label, len(s))
			return s, nil
		},
	}
}

func main() {
	r := New()
	r.Add(Redact())
	r.Add(Truncate(200))
	r.Add(Log(os.Stderr, "response"))

	prompt := "Please call the API with api_key=sk-VERY-SECRET-1234 and summarize the result."
	cleaned, err := r.Run(gsd.HookPre, prompt)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("--- after pre hooks ---")
	fmt.Println(cleaned)

	// Pretend a model responded:
	response := strings.Repeat("Result line.\n", 30)
	logged, _ := r.Run(gsd.HookPost, response)
	fmt.Println("--- after post hooks ---")
	fmt.Printf("(response length: %d bytes)\n", len(logged))
}
