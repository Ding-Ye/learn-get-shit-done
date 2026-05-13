package main

import (
	"reflect"
	"testing"
)

func TestParseCommand_table(t *testing.T) {
	cases := []struct {
		in        string
		wantName  string
		wantFlags map[string]string
		wantArgs  []string
	}{
		{
			in:        "/gsd capture --note \"Pay the bill\"",
			wantName:  "capture",
			wantFlags: map[string]string{"note": "Pay the bill"},
		},
		{
			in:        "gsd:plan --phase=03 --gaps",
			wantName:  "plan",
			wantFlags: map[string]string{"phase": "03", "gaps": ""},
		},
		{
			in:        "/gsd execute 03-auth",
			wantName:  "execute",
			wantFlags: map[string]string{},
			wantArgs:  []string{"03-auth"},
		},
		{
			in:        "fast --text hello world",
			wantName:  "fast",
			wantFlags: map[string]string{"text": "hello"},
			wantArgs:  []string{"world"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			cmd, err := ParseCommand(tc.in)
			if err != nil {
				t.Fatalf("ParseCommand: %v", err)
			}
			if cmd.Name != tc.wantName {
				t.Errorf("name = %q, want %q", cmd.Name, tc.wantName)
			}
			if !reflect.DeepEqual(cmd.Flags, tc.wantFlags) {
				t.Errorf("flags = %v, want %v", cmd.Flags, tc.wantFlags)
			}
			if !reflect.DeepEqual(cmd.Args, tc.wantArgs) {
				t.Errorf("args = %v, want %v", cmd.Args, tc.wantArgs)
			}
		})
	}
}

func TestParseCommand_errors(t *testing.T) {
	for _, in := range []string{"", "   ", `/gsd capture --note "unterminated`} {
		if _, err := ParseCommand(in); err == nil {
			t.Errorf("expected error for %q", in)
		}
	}
}
