// Chapter 02 — Command parser.
//
// Turn a raw input line like `/gsd capture --note "Pay the bill"` into a
// structured gsd.Command. We support:
//
//   - the leading slash is optional;
//   - an optional namespace prefix (`gsd ` or `gsd:`);
//   - `--flag` (boolean), `--flag value`, and `--flag=value`;
//   - quoted strings ("…" or '…') so multi-word values stay together.
//
// The parser is intentionally not as forgiving as a shell — we want
// deterministic behavior for tests, not feature parity with bash.
package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/Ding-Ye/learn-get-shit-done/gsd"
)

// ParseCommand parses a single command line into a gsd.Command.
func ParseCommand(line string) (gsd.Command, error) {
	raw := line
	line = strings.TrimSpace(line)
	if line == "" {
		return gsd.Command{}, errors.New("empty command line")
	}

	// Optional leading slash.
	line = strings.TrimPrefix(line, "/")

	tokens, err := tokenize(line)
	if err != nil {
		return gsd.Command{}, err
	}
	if len(tokens) == 0 {
		return gsd.Command{}, errors.New("no tokens after slash")
	}

	cmd := gsd.Command{Raw: raw, Flags: map[string]string{}}

	// First token is the command name. We allow either `gsd capture` (two
	// tokens, namespace + name) or `gsd:capture` (colon form).
	first := tokens[0]
	rest := tokens[1:]
	if strings.Contains(first, ":") {
		_, name, _ := strings.Cut(first, ":")
		cmd.Name = name
	} else if first == "gsd" && len(rest) > 0 {
		cmd.Name = rest[0]
		rest = rest[1:]
	} else {
		cmd.Name = first
	}

	// Walk the remaining tokens, recognizing flags.
	for i := 0; i < len(rest); i++ {
		t := rest[i]
		if !strings.HasPrefix(t, "--") {
			cmd.Args = append(cmd.Args, t)
			continue
		}
		// Strip leading dashes; the upstream uses long flags only.
		flag := strings.TrimPrefix(t, "--")
		if eq := strings.IndexByte(flag, '='); eq >= 0 {
			cmd.Flags[flag[:eq]] = flag[eq+1:]
			continue
		}
		// Look ahead: if the next token is a value (not another flag), consume it.
		if i+1 < len(rest) && !strings.HasPrefix(rest[i+1], "--") {
			cmd.Flags[flag] = rest[i+1]
			i++
		} else {
			cmd.Flags[flag] = ""
		}
	}

	return cmd, nil
}

// tokenize splits a string into shell-ish tokens, honouring quotes. We do not
// implement escape sequences — quotes are paired, and that is it.
func tokenize(s string) ([]string, error) {
	var tokens []string
	var cur strings.Builder
	inQuote := rune(0)
	for _, r := range s {
		switch {
		case inQuote != 0:
			if r == inQuote {
				inQuote = 0
			} else {
				cur.WriteRune(r)
			}
		case r == '"' || r == '\'':
			inQuote = r
		case unicode.IsSpace(r):
			if cur.Len() > 0 {
				tokens = append(tokens, cur.String())
				cur.Reset()
			}
		default:
			cur.WriteRune(r)
		}
	}
	if inQuote != 0 {
		return nil, fmt.Errorf("unterminated quote: %s", s)
	}
	if cur.Len() > 0 {
		tokens = append(tokens, cur.String())
	}
	return tokens, nil
}

func main() {
	line := `/gsd capture --note "Pay the electric bill"`
	if len(os.Args) > 1 {
		line = strings.Join(os.Args[1:], " ")
	}
	cmd, err := ParseCommand(line)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("raw   : %s\n", cmd.Raw)
	fmt.Printf("name  : %s\n", cmd.Name)
	fmt.Printf("flags : %v\n", cmd.Flags)
	fmt.Printf("args  : %v\n", cmd.Args)
}
