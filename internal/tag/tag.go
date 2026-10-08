// Package tag parses the value of a `moo:"..."` struct tag into structured
// options.
//
// Grammar:
//
//	empty                 -> zero Options
//	"-"                   -> Skip = true
//	"readonly"            -> ReadOnly = true
//	"required"            -> Required = true
//	"default=<expr>"      -> Default = "<expr>"   (raw, unparsed)
//
// Multiple keys are separated by ";". Unknown keys are an error so typos
// surface early.
package tag

import (
	"fmt"
	"strings"
)

// Options is the parsed form of a `moo:"..."` tag.
type Options struct {
	Skip     bool
	ReadOnly bool
	Required bool
	// Default is the raw right-hand side of `default=...`, unparsed. The
	// generator is responsible for interpreting it as a Go expression.
	Default string
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

		key, val, hasVal := strings.Cut(part, "=")
		key = strings.TrimSpace(key)

		switch key {
		case "readonly":
			if hasVal {
				return opts, fmt.Errorf("moo: readonly takes no value: %q", part)
			}
			opts.ReadOnly = true
		case "required":
			if hasVal {
				return opts, fmt.Errorf("moo: required takes no value: %q", part)
			}
			opts.Required = true
		case "default":
			if !hasVal {
				return opts, fmt.Errorf("moo: default requires a value: %q", part)
			}
			opts.Default = strings.TrimSpace(val)
		default:
			return opts, fmt.Errorf("unknown moo option: %q", part)
		}
	}
	return opts, nil
}
