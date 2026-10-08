package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, src string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "in.go")
	if err := os.WriteFile(p, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestParseFile_Basic(t *testing.T) {
	src := `package models

type Point struct {
	X int    ` + "`moo:\"default=0\"`" + `
	Y int    ` + "`moo:\"readonly\"`" + `
	name string
	SkipMe int ` + "`moo:\"-\"`" + `
}

type Ignored struct {
	A int
}
`
	p := writeTemp(t, src)
	f, err := ParseFile(p)
	if err != nil {
		t.Fatal(err)
	}

	if f.Package != "models" {
		t.Fatalf("package: got %q", f.Package)
	}
	if len(f.Structs) != 1 {
		t.Fatalf("expected 1 struct, got %d", len(f.Structs))
	}
	s := f.Structs[0]
	if s.Name != "Point" {
		t.Fatalf("name: got %q", s.Name)
	}
	if len(s.Fields) != 4 {
		t.Fatalf("expected 4 fields, got %d", len(s.Fields))
	}

	want := map[string]struct {
		typ      string
		readonly bool
		skip     bool
	}{
		"X":      {"int", false, false},
		"Y":      {"int", true, false},
		"name":   {"string", false, true}, // unexported → skip
		"SkipMe": {"int", false, true},
	}
	for _, fld := range s.Fields {
		w, ok := want[fld.Name]
		if !ok {
			t.Fatalf("unexpected field %q", fld.Name)
		}
		if fld.GoType != w.typ {
			t.Errorf("%s type: got %q want %q", fld.Name, fld.GoType, w.typ)
		}
		if fld.ReadOnly != w.readonly {
			t.Errorf("%s readonly: got %v want %v", fld.Name, fld.ReadOnly, w.readonly)
		}
		if fld.Skip != w.skip {
			t.Errorf("%s skip: got %v want %v", fld.Name, fld.Skip, w.skip)
		}
	}
}

func TestParseFile_TypesRendered(t *testing.T) {
	src := `package models

type Thing struct {
	A *int             ` + "`moo:\"\"`" + `
	B []string         ` + "`moo:\"\"`" + `
	C map[string]int   ` + "`moo:\"\"`" + `
	D []*Foo           ` + "`moo:\"\"`" + `
	E time.Time        ` + "`moo:\"\"`" + `
}
`
	p := writeTemp(t, src)
	f, err := ParseFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(f.Structs) != 1 {
		t.Fatalf("expected 1 struct, got %d", len(f.Structs))
	}
	got := map[string]string{}
	for _, fld := range f.Structs[0].Fields {
		got[fld.Name] = fld.GoType
	}
	want := map[string]string{
		"A": "*int",
		"B": "[]string",
		"C": "map[string]int",
		"D": "[]*Foo",
		"E": "time.Time",
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("%s: got %q want %q", k, got[k], v)
		}
	}
}

func TestParseFile_NoMooTags(t *testing.T) {
	src := `package models

type Plain struct {
	A int
}
`
	p := writeTemp(t, src)
	f, err := ParseFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(f.Structs) != 0 {
		t.Fatalf("expected 0 structs, got %d", len(f.Structs))
	}
}
