package tag

import "testing"

func TestParse(t *testing.T) {
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
			if got != tc.want {
				t.Fatalf("got %+v, want %+v", got, tc.want)
			}
		})
	}
}
