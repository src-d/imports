package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/src-d/imports"
	_ "github.com/src-d/imports/languages/all"
)

var (
	fRoot    = flag.String("root", ".", "root directory with the analyzed project")
	fRel     = flag.String("rel", "", "a directory relative to the root to analyze (recursively)")
	fNum     = flag.Int("n", 0, "max allowed concurrency (0 means use the number of CPUs)")
	fSymLink = flag.Bool("sym-links", false, "allow sym-links traversal")
)

func main() {
	flag.Parse()
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	e := imports.NewExtractor(imports.Config{
		Out:      os.Stdout,
		Num:      *fNum,
		SymLinks: *fSymLink,
	})
	return e.Extract(*fRoot, *fRel)
}
