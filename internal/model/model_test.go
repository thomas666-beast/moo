package model

import "testing"

func TestWritableFields(t *testing.T) {
	s := &StructSpec{
		Name: "Point",
		Fields: []*FieldSpec{
			{Name: "X", GoType: "int"},
			{Name: "Y", GoType: "int", ReadOnly: true},
			{Name: "Z", GoType: "int", Skip: true},
		},
	}
	got := s.WritableFields()
	if len(got) != 1 || got[0].Name != "X" {
		t.Fatalf("expected only X writable, got %+v", got)
	}
}
