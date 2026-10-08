package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/thomas666-beast/moo/internal/gen"
	"github.com/thomas666-beast/moo/internal/model"
	"github.com/thomas666-beast/moo/internal/parser"
)

func main() {
	var (
		debugFlag  = flag.Bool("debug", false, "print parsed IR instead of generating code")
		stdoutFlag = flag.Bool("stdout", false, "write generated code to stdout instead of a file")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: moo [-debug] [-stdout] <file.go>\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	path := flag.Arg(0)

	f, err := parser.ParseFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "moo: %v\n", err)
		os.Exit(1)
	}

	if *debugFlag {
		dump(f)
		return
	}

	if len(f.Structs) == 0 {
		fmt.Fprintln(os.Stderr, "moo: no moo-tagged structs found; nothing to do")
		return
	}

	src, err := gen.Generate(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "moo: %v\n", err)
		os.Exit(1)
	}

	if *stdoutFlag {
		os.Stdout.Write(src)
		return
	}

	outPath := outputPath(path)
	if err := os.WriteFile(outPath, src, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "moo: write %s: %v\n", outPath, err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "moo: wrote %s\n", outPath)
}

// outputPath returns the sibling path for generated code:
//
//	models/point.go     -> models/zz_moo_point.go
//	models/point_test.go -> models/zz_moo_point_test.go
func outputPath(src string) string {
	dir := filepath.Dir(src)
	base := filepath.Base(src)
	name := strings.TrimSuffix(base, ".go")
	return filepath.Join(dir, "zz_moo_"+name+".go")
}

func dump(f *model.File) {
	fmt.Printf("package %s\n", f.Package)
	if len(f.Structs) == 0 {
		fmt.Println("  (no moo structs)")
		return
	}
	for _, s := range f.Structs {
		fmt.Printf("  struct %s\n", s.Name)
		for _, fld := range s.Fields {
			fmt.Printf("    - %-10s %-15s", fld.Name, fld.GoType)
			if fld.ReadOnly {
				fmt.Print(" readonly")
			}
			if fld.Required {
				fmt.Print(" required")
			}
			if fld.Default != "" {
				fmt.Printf(" default=%s", fld.Default)
			}
			if fld.Skip {
				fmt.Print(" skip")
			}
			fmt.Println()
		}
	}
}
