// Chapter 06 — Agent dispatcher.
//
// Resolve an @agent reference, load the agent's markdown definition, build a
// prompt that includes the agent's persona, run pre-hooks, send via a
// Transport, run post-hooks on the response. Ships a tiny EchoTransport so
// the chapter is self-contained (no real LLM dependency).
//
// The dispatcher composes earlier chapters at the type level — we re-import
// `gsd` types and rely on the same agent-frontmatter shape — but stays small
// by keeping each helper local rather than reaching across packages.
package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Ding-Ye/learn-get-shit-done/gsd"
)

// EchoTransport returns the prompt prefixed with `ECHO:` — a deterministic
// stand-in for a real LLM call. Perfect for tests and demos.
type EchoTransport struct{}

func (EchoTransport) Send(prompt string) (gsd.TransportResult, error) {
	start := time.Now()
	res := "ECHO:\n" + prompt
	return gsd.TransportResult{Response: res, StartedAt: start, EndedAt: time.Now()}, nil
}

// Dispatcher wires together: agent loading, prompt building, hook chain, transport.
type Dispatcher struct {
	AgentsDir string
	Transport gsd.Transport

	preHooks  []gsd.Hook
	postHooks []gsd.Hook
}

func NewDispatcher(agentsDir string, t gsd.Transport) *Dispatcher {
	return &Dispatcher{AgentsDir: agentsDir, Transport: t}
}

func (d *Dispatcher) AddPre(h gsd.Hook)  { d.preHooks = append(d.preHooks, h) }
func (d *Dispatcher) AddPost(h gsd.Hook) { d.postHooks = append(d.postHooks, h) }

// Resolve finds and loads the agent referenced by an `@agent-name` token.
// The leading `@` is optional.
func (d *Dispatcher) Resolve(ref string) (gsd.Agent, error) {
	name := strings.TrimPrefix(ref, "@")
	if name == "" {
		return gsd.Agent{}, errors.New("empty agent reference")
	}
	path := filepath.Join(d.AgentsDir, name+".md")
	raw, err := os.ReadFile(path)
	if err != nil {
		return gsd.Agent{}, fmt.Errorf("load agent %q: %w", name, err)
	}
	return parseAgent(string(raw))
}

// Dispatch is the one-shot entrypoint: resolve the agent, build a prompt with
// the user's message, run pre-hooks, transport, run post-hooks.
func (d *Dispatcher) Dispatch(agentRef, userMsg string) (string, error) {
	agent, err := d.Resolve(agentRef)
	if err != nil {
		return "", err
	}
	prompt := buildPrompt(agent, userMsg)

	for _, h := range d.preHooks {
		if h.Phase != gsd.HookPre {
			continue
		}
		prompt, err = h.Fn(prompt)
		if err != nil {
			return "", fmt.Errorf("pre-hook %s: %w", h.Name, err)
		}
	}

	res, err := d.Transport.Send(prompt)
	if err != nil {
		return "", fmt.Errorf("transport: %w", err)
	}
	response := res.Response

	for _, h := range d.postHooks {
		if h.Phase != gsd.HookPost {
			continue
		}
		response, err = h.Fn(response)
		if err != nil {
			return "", fmt.Errorf("post-hook %s: %w", h.Name, err)
		}
	}
	return response, nil
}

var roleRe = regexp.MustCompile(`(?s)<role>(.*?)</role>`)

func parseAgent(text string) (gsd.Agent, error) {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	lines := strings.Split(text, "\n")
	a := gsd.Agent{}
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		a.Body = text
		return a, nil
	}
	closing := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			closing = i
			break
		}
	}
	if closing == -1 {
		return a, errors.New("unterminated frontmatter")
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
		k, v = strings.TrimSpace(k), strings.Trim(strings.TrimSpace(v), `"'`)
		switch k {
		case "name":
			a.Name = v
		case "description":
			a.Description = v
		case "tools":
			for _, t := range strings.Split(v, ",") {
				if tt := strings.TrimSpace(t); tt != "" {
					a.Tools = append(a.Tools, tt)
				}
			}
		}
	}
	a.Body = strings.TrimLeft(strings.Join(lines[closing+1:], "\n"), "\n")
	if m := roleRe.FindStringSubmatch(a.Body); m != nil {
		a.Role = strings.TrimSpace(m[1])
	}
	return a, nil
}

func buildPrompt(a gsd.Agent, user string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "You are %s.\n", a.Name)
	if a.Description != "" {
		fmt.Fprintf(&b, "%s\n", a.Description)
	}
	if a.Role != "" {
		b.WriteString("\nRole:\n")
		b.WriteString(a.Role)
		b.WriteString("\n")
	}
	if len(a.Tools) > 0 {
		fmt.Fprintf(&b, "\nAllowed tools: %s\n", strings.Join(a.Tools, ", "))
	}
	b.WriteString("\n=== USER ===\n")
	b.WriteString(strings.TrimSpace(user))
	b.WriteString("\n")
	return b.String()
}

func main() {
	dir := filepath.Join("agents", "s06-agent-dispatcher", "testdata")
	d := NewDispatcher(dir, EchoTransport{})

	// One observability hook on each side, just to prove the wiring.
	d.AddPre(gsd.Hook{Name: "tag-pre", Phase: gsd.HookPre, Fn: func(s string) (string, error) {
		return "(routed via dispatcher)\n" + s, nil
	}})
	d.AddPost(gsd.Hook{Name: "tag-post", Phase: gsd.HookPost, Fn: func(s string) (string, error) {
		return s + "\n(processed)", nil
	}})

	ref := "@gsd-planner"
	user := "Outline the next phase for the auth refactor."
	if len(os.Args) > 1 {
		ref = os.Args[1]
	}
	if len(os.Args) > 2 {
		user = strings.Join(os.Args[2:], " ")
	}
	out, err := d.Dispatch(ref, user)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(out)
}
