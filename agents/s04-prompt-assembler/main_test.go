package main

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/Ding-Ye/learn-get-shit-done/gsd"
)

func TestAssemble_allSectionsPresent(t *testing.T) {
	out := Assemble(Inputs{
		Skill:   gsd.Skill{Name: "gsd:plan", Description: "Plan a phase", Body: "plan body", AllowedTools: []string{"Read", "Write"}},
		Command: gsd.Command{Name: "plan", Flags: map[string]string{"phase": "03"}, Args: []string{"03-auth"}},
		Agent:   gsd.Agent{Name: "gsd-planner", Description: "Phase planner", Role: "you plan", Tools: []string{"Read"}},
		Context: "be concise",
		User:    "do the thing",
	})
	for _, want := range []string{"=== SYSTEM ===", "=== CONTEXT ===", "=== COMMAND ===", "=== USER ===", "gsd-planner", "be concise", "plan body", "--phase = 03", "03-auth", "do the thing"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n---\n%s", want, out)
		}
	}
}

func TestAssemble_deterministicOrder(t *testing.T) {
	in := Inputs{User: "u", Context: "c", Skill: gsd.Skill{Name: "n", Body: "b"}, Command: gsd.Command{Name: "n"}, Agent: gsd.Agent{Name: "a"}}
	got := Assemble(in)
	iSys := strings.Index(got, "=== SYSTEM ===")
	iCtx := strings.Index(got, "=== CONTEXT ===")
	iCmd := strings.Index(got, "=== COMMAND ===")
	iUsr := strings.Index(got, "=== USER ===")
	if !(iSys < iCtx && iCtx < iCmd && iCmd < iUsr) {
		t.Fatalf("section order broken: sys=%d ctx=%d cmd=%d user=%d", iSys, iCtx, iCmd, iUsr)
	}
}

func TestLoadAgent(t *testing.T) {
	a, err := loadAgent(filepath.Join("testdata", "agent-planner.md"))
	if err != nil {
		t.Fatalf("loadAgent: %v", err)
	}
	if a.Name != "gsd-planner" {
		t.Errorf("name = %q", a.Name)
	}
	if a.Role == "" {
		t.Errorf("role should be populated from <role>…</role>")
	}
}
