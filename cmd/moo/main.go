package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: moo <file.go>")
		os.Exit(2)
	}
	file := os.Args[1]
	fmt.Printf("moo: would process %s\n", file)
}
