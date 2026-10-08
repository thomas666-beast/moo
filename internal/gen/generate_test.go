package gen

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/thomas666-beast/moo/internal/model"
)

func TestGenerate_Point(t *testing.T) {
	f := &model.File{
		Package: "models",
		Structs: []*model.StructSpec{
			{
				Name: "Point",
				Fields: []*model.FieldSpec{
					{Name: "X", GoType: "int", Default: "0"},
					{Name: "Y", GoType: "int", Default: "0", Required: true},
					{Name: "Name", GoType: "string", Default: `"origin"`},
					{Name: "ID", GoType: "int64", ReadOnly: true},
					{Name: "Tags", GoType: "[]string"},
					{Name: "secret", GoType: "string", Skip: true},
					{Name: "Ignore", GoType: "int", Skip: true},
				},
			},
		},
	}

	got, err := Generate(f)
	if err != nil {
		t.Fatal(err)
	}

	wantPath := filepath.Join("testdata", "point.golden.go")
	want, err := os.ReadFile(wantPath)
	if err != nil {
		t.Fatalf("read golden: %v", err)
	}
	if string(got) != string(want) {
		t.Fatalf("mismatch.\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}
