package tag

import "testing"

func TestParse(t *testing.T) {
	atoi := func(n int) *int { return &n }

	cases := []struct {
		name    string
		in      string
		want    Options
		wantErr bool
	}{
		{"empty", "", Options{}, false},
		{"dash", "-", Options{Skip: true}, false},
		{"readonly", "readonly", Options{ReadOnly: true}, false},
		{"required", "required", Options{Required: true}, false},
		{"default", "default=0", Options{Default: "0"}, false},
		{"default-string", `default="origin"`, Options{Default: `"origin"`}, false},
		{"spaces", "  readonly  ", Options{ReadOnly: true}, false},
		{"multi-empty", "readonly;", Options{ReadOnly: true}, false},
		{"multi", "required;default=5", Options{Required: true, Default: "5"}, false},
		{"unknown", "banana", Options{}, true},
		{"readonly-val", "readonly=1", Options{}, true},
		{"default-noval", "default", Options{}, true},

		// extended
		{"min", "min=3", Options{Min: atoi(3)}, false},
		{"max", "max=10", Options{Max: atoi(10)}, false},
		{"min-max", "min=1;max=5", Options{Min: atoi(1), Max: atoi(5)}, false},
		{"enum", "enum=a|b|c", Options{Enum: []string{"a", "b", "c"}}, false},
		{"enum-spaces", "enum= a | b ", Options{Enum: []string{"a", "b"}}, false},
		{"match", "match=email", Options{Match: "email"}, false},
		{"combo", "required;min=1;max=9;match=email",
			Options{Required: true, Min: atoi(1), Max: atoi(9), Match: "email"}, false},
		{"bad-min", "min=abc", Options{}, true},
		{"no-min-val", "min", Options{}, true},
		{"empty-enum", "enum=", Options{}, true},
		{"enum-empty-item", "enum=a||b", Options{}, true},
		{"no-match-val", "match", Options{}, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Parse(tc.in)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("want error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !equalOptions(got, tc.want) {
				t.Fatalf("got %+v, want %+v", got, tc.want)
			}
		})
	}
}

func equalOptions(a, b Options) bool {
	if a.Skip != b.Skip || a.ReadOnly != b.ReadOnly || a.Required != b.Required ||
		a.Default != b.Default || a.Match != b.Match {
		return false
	}
	if !equalIntPtr(a.Min, b.Min) || !equalIntPtr(a.Max, b.Max) {
		return false
	}
	if len(a.Enum) != len(b.Enum) {
		return false
	}
	for i := range a.Enum {
		if a.Enum[i] != b.Enum[i] {
			return false
		}
	}
	return true
}

func equalIntPtr(a, b *int) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}
