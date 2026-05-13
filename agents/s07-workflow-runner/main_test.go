package main

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadWorkflow_note(t *testing.T) {
	w, err := LoadWorkflow(filepath.Join("testdata", "note.md"))
	if err != nil {
		t.Fatalf("LoadWorkflow: %v", err)
	}
	if w.Name != "note" {
		t.Errorf("name = %q", w.Name)
	}
	if len(w.Steps) != 2 {
		t.Fatalf("steps = %d, want 2", len(w.Steps))
	}
	if w.Steps[0].ID != "append" || w.Steps[0].Next != "list" {
		t.Errorf("step 0 wrong: %+v", w.Steps[0])
	}
}

func TestRun_orderedAndSubstituted(t *testing.T) {
	w, err := LoadWorkflow(filepath.Join("testdata", "note.md"))
	if err != nil {
		t.Fatalf("LoadWorkflow: %v", err)
	}
	var ran []string
	exec := func(cmd, args string) (string, error) {
		ran = append(ran, cmd+":"+args)
		return "ok", nil
	}
	if _, err := Run(w, map[string]string{"Text": "hi"}, exec); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if len(ran) != 2 {
		t.Fatalf("ran %d steps, want 2 (%v)", len(ran), ran)
	}
	if !strings.Contains(ran[0], "--note hi") {
		t.Errorf("first arg didn't substitute Text: %v", ran[0])
	}
}

func TestRun_onError(t *testing.T) {
	w, err := LoadWorkflow(filepath.Join("testdata", "with-error.md"))
	if err != nil {
		t.Fatalf("LoadWorkflow: %v", err)
	}
	ran := 0
	exec := func(cmd, args string) (string, error) {
		ran++
		if cmd == "broken" {
			return "", errors.New("boom")
		}
		return "ok", nil
	}
	trace, err := Run(w, nil, exec)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if !strings.Contains(trace, "routing to") {
		t.Errorf("trace should show error routing: %s", trace)
	}
	if ran < 2 {
		t.Errorf("expected at least 2 step executions, got %d", ran)
	}
}
