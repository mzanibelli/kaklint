package main

import (
	"fmt"
	"kaklint"
	"os"
)

func main() {
	if err := run(os.Args...); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		os.Exit(1)
	}
}

func run(args ...string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: %s FILETYPE FILE", args[0])
	}
	return kaklint.Default.Lint(args[1], args[2])
}
