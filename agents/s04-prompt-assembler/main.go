// Chapter 04 — Prompt assembler.
//
// Now we put the first three chapters together. Given:
//
//   - a parsed gsd.Command  (chapter 2)
//   - a chosen context name + variables to render (chapter 3)
//   - a gsd.Agent persona   (loaded the same way as a skill in chapter 1)
//   - the raw user input text
//
// produce the final prompt string a transport would actually send.
//
// The assembled format is intentionally human-readable so you can pipe the
// output of this chapter into `less` and read it. Each section is delimited
// with an ASCII heading, in the order: system → context → command → user.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Ding-Ye/learn-get-shit-done/gsd"
)

// Inputs bundles every piece of state the assembler needs.
type Inputs struct {
	Skill   gsd.Skill   // command definition (chapter 1)
	Command gsd.Command // parsed invocation (chapter 2)
	Agent   gsd.Agent   // persona (chapter 6 loads these; here we pass directly)
	Context string      // already-rendered context fragment (chapter 3)
	User    string      // raw user input text
}

// Assemble joins the inputs into a single deterministic string.
func Assemble(in Inputs) string {
	var b strings.Builder

	b.WriteString("=== SYSTEM ===\n")
	if in.Agent.Name != "" {
		fmt.Fprintf(&b, "You are %s.\n", in.Agent.Name)
		if in.Agent.Description != "" {
			fmt.Fprintf(&b, "%s\n", in.Agent.Description)
		}
		if in.Agent.Role != "" {
			b.WriteString("\nRole:\n")
			b.WriteString(in.Agent.Role)
			b.WriteString("\n")
		}
	} else {
		b.WriteString("You are an assistant.\n")
	}
	if len(in.Agent.Tools) > 0 {
		fmt.Fprintf(&b, "\nAllowed tools: %s\n", strings.Join(in.Agent.Tools, ", "))
	}

	if in.Context != "" {
		b.WriteString("\n=== CONTEXT ===\n")
		b.WriteString(strings.TrimRight(in.Context, "\n"))
		b.WriteString("\n")
	}

	if in.Skill.Name != "" {
		b.WriteString("\n=== COMMAND ===\n")
		fmt.Fprintf(&b, "Name: %s\n", in.Skill.Name)
		if in.Skill.Description != "" {
			fmt.Fprintf(&b, "Summary: %s\n", in.Skill.Description)
		}
		if len(in.Skill.AllowedTools) > 0 {
			fmt.Fprintf(&b, "Skill tools: %s\n", strings.Join(in.Skill.AllowedTools, ", "))
		}
		b.WriteString("\nSkill body:\n")
		b.WriteString(strings.TrimRight(in.Skill.Body, "\n"))
		b.WriteString("\n")
	}

	if in.Command.Name != "" {
		b.WriteString("\nInvocation:\n")
		fmt.Fprintf(&b, "  command: %s\n", in.Command.Name)
		if len(in.Command.Flags) > 0 {
			for k, v := range in.Command.Flags {
				if v == "" {
					fmt.Fprintf(&b, "  --%s\n", k)
				} else {
					fmt.Fprintf(&b, "  --%s = %s\n", k, v)
				}
			}
		}
		if len(in.Command.Args) > 0 {
			fmt.Fprintf(&b, "  args: %s\n", strings.Join(in.Command.Args, " "))
		}
	}

	b.WriteString("\n=== USER ===\n")
	b.WriteString(strings.TrimSpace(in.User))
	b.WriteString("\n")

	return b.String()
}

// loadAgent is a 20-line agent loader — same shape as the chapter-1 skill loader.
// We keep it inline here so the chapter compiles on its own without
// reaching across chapter boundaries; chapter 6 will replace it with a real
// dispatcher.
func loadAgent(path string) (gsd.Agent, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return gsd.Agent{}, err
	}
	text := strings.ReplaceAll(string(raw), "\r\n", "\n")
	lines := strings.Split(text, "\n")
	agent := gsd.Agent{}
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		agent.Body = text
		return agent, nil
	}
	closing := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			closing = i
			break
		}
	}
	if closing == -1 {
		return agent, fmt.Errorf("unterminated frontmatter")
	}
	for _, line := range lines[1:closing] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		key, value = strings.TrimSpace(key), strings.Trim(strings.TrimSpace(value), `"'`)
		switch key {
		case "name":
			agent.Name = value
		case "description":
			agent.Description = value
		case "tools":
			for _, t := range strings.Split(value, ",") {
				if tt := strings.TrimSpace(t); tt != "" {
					agent.Tools = append(agent.Tools, tt)
				}
			}
		}
	}
	body := strings.Join(lines[closing+1:], "\n")
	agent.Body = strings.TrimLeft(body, "\n")
	// Extract <role>…</role>.
	if start := strings.Index(agent.Body, "<role>"); start >= 0 {
		end := strings.Index(agent.Body, "</role>")
		if end > start {
			agent.Role = strings.TrimSpace(agent.Body[start+len("<role>") : end])
		}
	}
	return agent, nil
}

func main() {
	// Hard-coded demo using the bundled fixtures. The chapter-by-chapter
	// integration chapter wires this up dynamically.
	td := filepath.Join("agents", "s04-prompt-assembler", "testdata")
	skillBytes, _ := os.ReadFile(filepath.Join(td, "plan.md"))
	skillText := string(skillBytes)

	// Reuse a minimal inline parse for the skill — same shape as chapter 1.
	skill := gsd.Skill{Name: "gsd:plan", Description: "Plan a phase", Body: skillText}
	if i := strings.Index(skillText, "\n---\n"); i >= 0 {
		// keep the frontmatter trick simple
		skill.Body = strings.TrimSpace(skillText[i+len("\n---\n"):])
	}

	agent, err := loadAgent(filepath.Join(td, "agent-planner.md"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	in := Inputs{
		Skill:   skill,
		Command: gsd.Command{Name: "plan", Flags: map[string]string{"phase": "03"}, Args: []string{"03-auth"}, Raw: "/gsd plan --phase=03 03-auth"},
		Agent:   agent,
		Context: "Output style: concise, action-oriented. User: Ding.",
		User:    "Plan the next phase for the auth refactor.",
	}
	fmt.Print(Assemble(in))
}
