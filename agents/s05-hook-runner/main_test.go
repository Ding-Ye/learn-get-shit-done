package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Ding-Ye/learn-get-shit-done/gsd"
)

func TestRedact(t *testing.T) {
	r := New()
	r.Add(Redact())
	out, err := r.Run(gsd.HookPre, "use api_key=sk-1234 and password=hunter2")
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if strings.Contains(out, "sk-1234") || strings.Contains(out, "hunter2") {
		t.Errorf("secrets not redacted: %s", out)
	}
}

func TestTruncate(t *testing.T) {
	r := New()
	r.Add(Truncate(10))
	out, _ := r.Run(gsd.HookPre, "12345678901234567890")
	if !strings.HasSuffix(out, "[truncated]") {
		t.Errorf("expected truncation marker, got %q", out)
	}
}

func TestLog_doesNotMutate(t *testing.T) {
	var buf bytes.Buffer
	r := New()
	r.Add(Log(&buf, "test"))
	out, _ := r.Run(gsd.HookPost, "hello")
	if out != "hello" {
		t.Errorf("Log mutated payload: %q", out)
	}
	if !strings.Contains(buf.String(), "test") {
		t.Errorf("log didn't write: %q", buf.String())
	}
}

func TestRun_phaseFilter(t *testing.T) {
	r := New()
	r.Add(Redact())                  // HookPre
	r.Add(Log(new(bytes.Buffer), ""))// HookPost
	out, err := r.Run(gsd.HookPre, "api_key=sk-X")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "sk-X") {
		t.Errorf("pre run should have applied redact")
	}
}

func TestRun_propagatesError(t *testing.T) {
	r := New()
	r.Add(gsd.Hook{Name: "boom", Phase: gsd.HookPre, Fn: func(s string) (string, error) {
		return "", errBoom
	}})
	if _, err := r.Run(gsd.HookPre, "x"); err == nil {
		t.Fatal("expected error")
	}
}

var errBoom = stringError("boom")

type stringError string

func (e stringError) Error() string { return string(e) }
