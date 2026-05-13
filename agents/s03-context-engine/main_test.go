package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestRender_dev(t *testing.T) {
	eng := New(filepath.Join("testdata"))
	out, err := eng.Render("dev", map[string]string{"User": "Alice", "Task": "ship it"})
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(out, "Alice") || !strings.Contains(out, "ship it") {
		t.Errorf("expected vars to be substituted, got:\n%s", out)
	}
}

func TestRender_missingVar(t *testing.T) {
	eng := New(filepath.Join("testdata"))
	// `dev.md` references {{.User}} and {{.Task}}; omit Task.
	_, err := eng.Render("dev", map[string]string{"User": "Alice"})
	if err == nil {
		t.Fatal("expected missingkey error")
	}
}

func TestLoad_defaultsFromFrontmatter(t *testing.T) {
	eng := New(filepath.Join("testdata"))
	// `review.md` has a default verbosity in frontmatter.
	out, err := eng.Render("review", map[string]string{"Target": "PR-42"})
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(out, "verbosity: medium") {
		t.Errorf("expected default verbosity from frontmatter, got:\n%s", out)
	}
}
