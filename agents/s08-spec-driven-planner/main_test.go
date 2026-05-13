package main

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Ding-Ye/learn-get-shit-done/gsd"
)

func TestLoadSpec_authRefactor(t *testing.T) {
	s, err := LoadSpec(filepath.Join("testdata", "auth-refactor.md"))
	if err != nil {
		t.Fatalf("LoadSpec: %v", err)
	}
	if !strings.Contains(s.Title, "Auth refactor") {
		t.Errorf("title = %q", s.Title)
	}
	if len(s.Phases) != 2 {
		t.Fatalf("phases = %d, want 2", len(s.Phases))
	}
	p1 := s.Phases[0]
	if p1.ID != "P01" || !strings.Contains(p1.Goal, "Mint") {
		t.Errorf("phase 1 wrong: %+v", p1)
	}
	if len(p1.Tasks) < 2 {
		t.Errorf("phase 1 tasks = %d, want >=2 (%v)", len(p1.Tasks), p1.Tasks)
	}
}

func TestExecutePhase_marksComplete(t *testing.T) {
	s, _ := LoadSpec(filepath.Join("testdata", "auth-refactor.md"))
	updated, log, err := ExecutePhase(s, "P01", func(p gsd.SpecPhase) error { return nil })
	if err != nil {
		t.Fatalf("ExecutePhase: %v", err)
	}
	if !updated.Phases[0].Complete {
		t.Errorf("P01 not marked complete")
	}
	if !strings.Contains(log, "Phase P01 complete") {
		t.Errorf("log missing completion line:\n%s", log)
	}
}

func TestExecutePhase_verifierError(t *testing.T) {
	s, _ := LoadSpec(filepath.Join("testdata", "auth-refactor.md"))
	_, _, err := ExecutePhase(s, "P01", func(p gsd.SpecPhase) error { return errors.New("fail") })
	if err == nil {
		t.Fatal("expected verifier error")
	}
}

func TestNextIncomplete(t *testing.T) {
	s, _ := LoadSpec(filepath.Join("testdata", "auth-refactor.md"))
	p, ok := NextIncomplete(s)
	if !ok || p.ID != "P01" {
		t.Errorf("NextIncomplete = %s,%v; want P01,true", p.ID, ok)
	}
	updated, _, _ := ExecutePhase(s, "P01", func(p gsd.SpecPhase) error { return nil })
	p, ok = NextIncomplete(updated)
	if !ok || p.ID != "P02" {
		t.Errorf("NextIncomplete after P01 = %s,%v; want P02,true", p.ID, ok)
	}
}
