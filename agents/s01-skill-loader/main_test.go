package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestParseSkill_capture(t *testing.T) {
	const src = `---
name: gsd:capture
description: Capture ideas, tasks, notes, and seeds
argument-hint: "[--note] [text]"
allowed-tools:
  - Read
  - Write
  - Bash
---

# Capture

Capture ideas to the right place.
`
	skill, err := ParseSkill(src)
	if err != nil {
		t.Fatalf("ParseSkill: %v", err)
	}
	if skill.Name != "gsd:capture" {
		t.Errorf("name = %q, want gsd:capture", skill.Name)
	}
	if skill.ArgumentHint != "[--note] [text]" {
		t.Errorf("argument-hint = %q", skill.ArgumentHint)
	}
	if got, want := strings.Join(skill.AllowedTools, ","), "Read,Write,Bash"; got != want {
		t.Errorf("allowed-tools = %q, want %q", got, want)
	}
	if !strings.Contains(skill.Body, "Capture ideas to the right place.") {
		t.Errorf("body missing expected text: %q", skill.Body)
	}
}

func TestParseSkill_inlineTools(t *testing.T) {
	const src = `---
name: gsd:plan
description: Plan a phase
allowed-tools: Read, Write, Edit, Bash
---

body
`
	skill, err := ParseSkill(src)
	if err != nil {
		t.Fatalf("ParseSkill: %v", err)
	}
	if got, want := len(skill.AllowedTools), 4; got != want {
		t.Fatalf("len(AllowedTools) = %d, want %d (%v)", got, want, skill.AllowedTools)
	}
}

func TestParseSkill_missingFrontmatter(t *testing.T) {
	_, err := ParseSkill("no fence here")
	if err == nil {
		t.Fatal("expected error for missing frontmatter")
	}
}

func TestLoadSkill_fixture(t *testing.T) {
	skill, err := LoadSkill(filepath.Join("testdata", "capture.md"))
	if err != nil {
		t.Fatalf("LoadSkill: %v", err)
	}
	if skill.Name == "" {
		t.Fatal("fixture produced empty Name")
	}
}
