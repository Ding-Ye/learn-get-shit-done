// Chapter 01 — Skill loader.
//
// A "skill" in the upstream framework is a markdown file with YAML
// frontmatter. This program reads such a file from disk and produces a
// runtime gsd.Skill. We deliberately hand-roll a tiny frontmatter parser
// (no yaml.v3 dependency) — the upstream grammar we care about is small:
//
//   - the file MUST start with a line `---`;
//   - frontmatter lines are either `key: value` or, for lists, a header line
//     `key:` followed by `  - item` bullets;
//   - the frontmatter ends at the next standalone `---` line;
//   - everything after that is the markdown body.
package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Ding-Ye/learn-get-shit-done/gsd"
)

// LoadSkill reads a skill markdown file and returns a populated gsd.Skill.
func LoadSkill(path string) (gsd.Skill, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return gsd.Skill{}, fmt.Errorf("read skill: %w", err)
	}
	return ParseSkill(string(raw))
}

// ParseSkill is the pure-string entry point — useful for tests.
func ParseSkill(text string) (gsd.Skill, error) {
	// Normalise line endings and split.
	text = strings.ReplaceAll(text, "\r\n", "\n")
	lines := strings.Split(text, "\n")

	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return gsd.Skill{}, errors.New("missing opening --- frontmatter fence")
	}

	// Find the closing fence.
	closing := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			closing = i
			break
		}
	}
	if closing == -1 {
		return gsd.Skill{}, errors.New("missing closing --- frontmatter fence")
	}

	skill := gsd.Skill{}
	if err := parseFrontmatter(lines[1:closing], &skill); err != nil {
		return gsd.Skill{}, err
	}
	skill.Body = strings.TrimLeft(strings.Join(lines[closing+1:], "\n"), "\n")
	return skill, nil
}

// parseFrontmatter handles the subset of YAML we need: `key: value` and a
// header line `key:` followed by `  - item` bullets. Anything else is ignored
// for forward-compatibility.
func parseFrontmatter(lines []string, out *gsd.Skill) error {
	var currentList *[]string // points to the list we're currently appending to
	for _, line := range lines {
		// List item — must precede with two-space indent.
		if currentList != nil && strings.HasPrefix(line, "  - ") {
			item := strings.TrimSpace(strings.TrimPrefix(line, "  - "))
			item = strings.Trim(item, `"'`)
			*currentList = append(*currentList, item)
			continue
		}

		// Non-list line ends any active list collection.
		currentList = nil

		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		key, value, ok := strings.Cut(trimmed, ":")
		if !ok {
			return fmt.Errorf("invalid frontmatter line: %q", line)
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		value = strings.Trim(value, `"'`)

		switch key {
		case "name":
			out.Name = value
		case "description":
			out.Description = value
		case "argument-hint":
			out.ArgumentHint = value
		case "allowed-tools":
			if value == "" {
				// Header for a list — collect subsequent bullet lines.
				currentList = &out.AllowedTools
			} else {
				// Inline form: `allowed-tools: Read, Write` — split on commas.
				for _, t := range strings.Split(value, ",") {
					if tt := strings.TrimSpace(t); tt != "" {
						out.AllowedTools = append(out.AllowedTools, tt)
					}
				}
			}
		}
	}
	return nil
}

func main() {
	// Default to the bundled fixture so `go run ./agents/s01-skill-loader` Just Works.
	path := filepath.Join("agents", "s01-skill-loader", "testdata", "capture.md")
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	skill, err := LoadSkill(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading %s: %v\n", path, err)
		os.Exit(1)
	}
	fmt.Printf("Loaded skill: %s\n", skill.Name)
	fmt.Printf("  description : %s\n", skill.Description)
	fmt.Printf("  arg-hint    : %s\n", skill.ArgumentHint)
	fmt.Printf("  tools       : %s\n", strings.Join(skill.AllowedTools, ", "))
	fmt.Printf("  body (first 80 chars): %s\n", firstLine(skill.Body, 80))
}

func firstLine(s string, n int) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		s = s[:i]
	}
	if len(s) > n {
		s = s[:n] + "…"
	}
	return s
}
