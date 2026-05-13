// Chapter 03 — Context engine.
//
// A "context" in the upstream framework is a reusable prompt fragment stored
// on disk. The context engine loads such a fragment by name, parses the
// optional YAML frontmatter that lists template variables, and renders the
// body via Go's text/template package with the variables the caller provides.
//
// This is the smallest sensible model. The upstream context-engine.ts is a
// per-phase file manifest with truncation; here we strip all of that and keep
// only what is conceptually new: a template directory + variable substitution.
package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Ding-Ye/learn-get-shit-done/gsd"
)

// Engine resolves contexts by name from a single root directory.
type Engine struct {
	Root string // directory containing <name>.md context templates
}

// New constructs an Engine rooted at dir.
func New(dir string) *Engine { return &Engine{Root: dir} }

// Load returns the parsed context for `name` (without `.md`), with no
// rendering applied yet — useful for inspection.
func (e *Engine) Load(name string) (gsd.Context, error) {
	path := filepath.Join(e.Root, name+".md")
	raw, err := os.ReadFile(path)
	if err != nil {
		return gsd.Context{}, fmt.Errorf("read context %q: %w", name, err)
	}
	return parseContext(name, string(raw))
}

// Render loads `name` and renders its body using `vars` as the template data.
// Unknown variables left in the template surface as text/template errors.
func (e *Engine) Render(name string, vars map[string]string) (string, error) {
	ctx, err := e.Load(name)
	if err != nil {
		return "", err
	}
	// Merge: caller-provided vars win over frontmatter defaults.
	merged := map[string]string{}
	for k, v := range ctx.Variables {
		merged[k] = v
	}
	for k, v := range vars {
		merged[k] = v
	}
	tmpl, err := template.New(name).Option("missingkey=error").Parse(ctx.Body)
	if err != nil {
		return "", fmt.Errorf("parse template %q: %w", name, err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, merged); err != nil {
		return "", fmt.Errorf("execute template %q: %w", name, err)
	}
	return buf.String(), nil
}

// parseContext peels off an optional `--- … ---` frontmatter block. Every
// frontmatter `key: value` pair becomes a default value in ctx.Variables.
func parseContext(name, text string) (gsd.Context, error) {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	lines := strings.Split(text, "\n")
	ctx := gsd.Context{Name: name, Variables: map[string]string{}}

	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "---" {
		closing := -1
		for i := 1; i < len(lines); i++ {
			if strings.TrimSpace(lines[i]) == "---" {
				closing = i
				break
			}
		}
		if closing == -1 {
			return ctx, errors.New("unterminated frontmatter")
		}
		for _, line := range lines[1:closing] {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			key, value, ok := strings.Cut(line, ":")
			if !ok {
				continue
			}
			ctx.Variables[strings.TrimSpace(key)] = strings.Trim(strings.TrimSpace(value), `"'`)
		}
		ctx.Body = strings.TrimLeft(strings.Join(lines[closing+1:], "\n"), "\n")
	} else {
		ctx.Body = text
	}
	return ctx, nil
}

func main() {
	root := filepath.Join("agents", "s03-context-engine", "testdata")
	name := "dev"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	user := "Ding"
	if len(os.Args) > 2 {
		user = os.Args[2]
	}
	eng := New(root)
	rendered, err := eng.Render(name, map[string]string{"User": user, "Task": "fix the login bug"})
	if err != nil {
		fmt.Fprintf(os.Stderr, "render: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(rendered)
}
