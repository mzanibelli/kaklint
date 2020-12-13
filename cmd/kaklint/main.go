package main

import (
	"fmt"
	"kaklint"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		usage()
	}
	if err := kaklint.Default.Lint(os.Args[1], os.Args[2]); err != nil {
		exit(err)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s FILETYPE FILE\n", os.Args[0])
	os.Exit(1)
}

func exit(err error) {
	fmt.Fprintf(os.Stderr, "%s: %v", os.Args[0], err)
	os.Exit(1)
}
