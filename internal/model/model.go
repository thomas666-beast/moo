// Package model defines the intermediate representation (IR) that the
// parser produces and the generator consumes.
package model

// File is the top-level IR for one source file.
type File struct {
	Package string        // Go package name, e.g. "models"
	Structs []*StructSpec // structs we will generate code for
}

// StructSpec describes one struct to generate code for.
type StructSpec struct {
	Name   string       // e.g. "Point"
	Fields []*FieldSpec // in declaration order
}

// FieldSpec describes one field of a struct.
type FieldSpec struct {
	Name     string
	GoType   string
	ReadOnly bool
	Skip     bool
	Required bool
	Default  string
	Min      *int
	Max      *int
	Enum     []string
	Match    string
}

// WritableFields returns fields that should get a With<Name> option.
func (s *StructSpec) WritableFields() []*FieldSpec {
	out := make([]*FieldSpec, 0, len(s.Fields))
	for _, f := range s.Fields {
		if !f.Skip && !f.ReadOnly {
			out = append(out, f)
		}
	}
	return out
}
