// Package tag parses the value of a `moo:"..."` struct tag into structured
// options. Grammar (M0):
//
//	empty            -> zero Options
//	"-"              -> Skip = true
//	"readonly"       -> ReadOnly = true
//	"readonly" only  -> (future keys land here)
//
// Multiple keys are separated by ";". Unknown keys are an error for now so
// typos surface early.
package tag

import (
	"fmt"
	"strings"
)

// Options is the parsed form of a `moo:"..."` tag.
type Options struct {
	Skip     bool
	ReadOnly bool
}

// Parse parses a raw tag value. It returns an error for malformed input or
// unknown keys.
func Parse(raw string) (Options, error) {
	var opts Options
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return opts, nil
	}
	if raw == "-" {
		opts.Skip = true
		return opts, nil
	}

	for _, part := range strings.Split(raw, ";") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		switch part {
		case "readonly":
			opts.ReadOnly = true
		default:
			return opts, fmt.Errorf("unknown moo option: %q", part)
		}
	}
	return opts, nil
}
