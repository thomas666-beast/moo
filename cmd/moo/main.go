package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/thomas666-beast/moo/internal/model"
	"github.com/thomas666-beast/moo/internal/parser"
)

func main() {
	debug := flag.Bool("debug", true, "print parsed IR instead of generating code (M0 default)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: moo [-debug] <file.go>\n")
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

	if *debug {
		dump(f)
		return
	}

	// generation path lands in Task 6
	fmt.Fprintln(os.Stderr, "moo: code generation not implemented yet")
	os.Exit(1)
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
