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
	migrateTypes := flags.Bool("migrate-types", false, "migrate indexed legacy API type strings in place")
	check := flags.Bool("check", false, "check whether an API type migration is required without writing")
	if err := flags.Parse(arguments); err != nil {
		return err
	}

	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected positional arguments: %v", flags.Args())
	}

	if *migrateTypes {
		if *pagesRoot == "" {
			return fmt.Errorf("-pages is required with -migrate-types")
		}

		if *referencePath != "" || *catalogPath != "" {
			return fmt.Errorf("-reference and -catalog cannot be used with -migrate-types")
		}

		return publisher.MigrateTypes(*pagesRoot, *check)
	}

	if *check {
		return fmt.Errorf("-check requires -migrate-types")
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
