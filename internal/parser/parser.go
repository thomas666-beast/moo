// Package parser reads a single Go source file and produces a model.File IR.
package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strconv"
	"strings"

	"github.com/thomas666-beast/moo/internal/model"
	"github.com/thomas666-beast/moo/internal/tag"
)

// ParseFile parses the given file and returns the IR for all structs that
// carry at least one `moo:` tagged field.
func ParseFile(path string) (*model.File, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}

	out := &model.File{Package: f.Name.Name}

	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}
		for _, spec := range gd.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				continue
			}
			ss, ok, err := parseStruct(ts.Name.Name, st)
			if err != nil {
				return nil, fmt.Errorf("struct %s: %w", ts.Name.Name, err)
			}
			if ok {
				out.Structs = append(out.Structs, ss)
			}
		}
	}

	return out, nil
}

// parseStruct returns (spec, true, nil) if the struct has at least one
// `moo:` tagged field; otherwise (nil, false, nil).
func parseStruct(name string, st *ast.StructType) (*model.StructSpec, bool, error) {
	ss := &model.StructSpec{Name: name}
	hasMoo := false

	for _, field := range st.Fields.List {
		if len(field.Names) == 0 {
			// embedded field — skip for M0
			continue
		}

		var (
			raw     string
			hasTag  bool
		)
		if field.Tag != nil {
			if unq, err := strconv.Unquote(field.Tag.Value); err == nil {
				if v, ok := reflect.StructTag(unq).Lookup("moo"); ok {
					raw = v
					hasTag = true
					hasMoo = true
				}
			}
		}

		opts, err := tag.Parse(raw)
		if err != nil {
			return nil, false, err
		}

		for _, ident := range field.Names {
			fs := &model.FieldSpec{
				Name:     ident.Name,
				GoType:   exprString(field.Type),
				ReadOnly: opts.ReadOnly,
				Skip:     opts.Skip,
				Required: opts.Required,
				Default:  opts.Default,
			}
			if !ident.IsExported() {
				fs.Skip = true
			}
			ss.Fields = append(ss.Fields, fs)
		}
		_ = hasTag // kept for clarity; presence is what flips hasMoo
	}

	if !hasMoo {
		return nil, false, nil
	}
	return ss, true, nil
}

// exprString renders a Go type expression back into source form.
func exprString(e ast.Expr) string {
	switch t := e.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + exprString(t.X)
	case *ast.ArrayType:
		if t.Len == nil {
			return "[]" + exprString(t.Elt)
		}
		return "[...]" + exprString(t.Elt)
	case *ast.MapType:
		return "map[" + exprString(t.Key) + "]" + exprString(t.Value)
	case *ast.SelectorExpr:
		return exprString(t.X) + "." + t.Sel.Name
	case *ast.InterfaceType:
		return "any"
	case *ast.ChanType:
		return "chan " + exprString(t.Value)
	case *ast.Ellipsis:
		return "..." + exprString(t.Elt)
	case *ast.ParenExpr:
		return "(" + exprString(t.X) + ")"
	case *ast.FuncType:
		return "func()" // best effort; we don't support funcs in M0 output
	default:
		return strings.TrimSpace(fmt.Sprintf("%T", e))
	}
}
