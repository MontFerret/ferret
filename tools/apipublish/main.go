package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/MontFerret/ferret/v2/tools/apipublish/internal/publisher"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}

func run(arguments []string) error {
	flags := flag.NewFlagSet("apipublish", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	referencePath := flags.String("reference", "", "generated API Reference path")
	catalogPath := flags.String("catalog", "", "generated API Catalog path")
	pagesRoot := flags.String("pages", "", "checked-out gh-pages root")
	if err := flags.Parse(arguments); err != nil {
		return err
	}

	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected positional arguments: %v", flags.Args())
	}

	if *referencePath == "" || *catalogPath == "" || *pagesRoot == "" {
		return fmt.Errorf("-reference, -catalog, and -pages are required")
	}

	referenceData, err := os.ReadFile(*referencePath)
	if err != nil {
		return fmt.Errorf("read generated API Reference: %w", err)
	}

	catalogData, err := os.ReadFile(*catalogPath)
	if err != nil {
		return fmt.Errorf("read generated API Catalog: %w", err)
	}

	return publisher.Publish(*pagesRoot, referenceData, catalogData)
}
