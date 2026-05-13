package main

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/Ding-Ye/learn-get-shit-done/gsd"
)

func TestResolve(t *testing.T) {
	d := NewDispatcher(filepath.Join("testdata"), EchoTransport{})
	a, err := d.Resolve("@gsd-planner")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if a.Name != "gsd-planner" || a.Role == "" {
		t.Errorf("agent not fully populated: %+v", a)
	}
}

func TestResolve_missing(t *testing.T) {
	d := NewDispatcher(filepath.Join("testdata"), EchoTransport{})
	if _, err := d.Resolve("@does-not-exist"); err == nil {
		t.Fatal("expected error for missing agent")
	}
}

func TestDispatch_endToEnd(t *testing.T) {
	d := NewDispatcher(filepath.Join("testdata"), EchoTransport{})
	called := 0
	d.AddPre(gsd.Hook{Name: "p", Phase: gsd.HookPre, Fn: func(s string) (string, error) {
		called++
		return "[PRE]\n" + s, nil
	}})
	d.AddPost(gsd.Hook{Name: "q", Phase: gsd.HookPost, Fn: func(s string) (string, error) {
		called++
		return s + "\n[POST]", nil
	}})
	out, err := d.Dispatch("@gsd-planner", "plan it")
	if err != nil {
		t.Fatalf("Dispatch: %v", err)
	}
	if called != 2 {
		t.Errorf("hook called count = %d, want 2", called)
	}
	if !strings.HasPrefix(out, "ECHO:\n[PRE]\nYou are gsd-planner.") {
		t.Errorf("unexpected output:\n%s", out)
	}
	if !strings.HasSuffix(strings.TrimRight(out, "\n"), "[POST]") {
		t.Errorf("post-hook didn't apply:\n%s", out)
	}
}
