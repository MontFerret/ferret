package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver/v3"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib"
	"github.com/MontFerret/ferret/v2/tools/apiref/internal/analyzer"
)

func main() {
	if err := run(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}

func run(ctx context.Context, arguments []string) error {
	flags := flag.NewFlagSet("apiref", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	version := flags.String("version", "", "canonical unprefixed semantic version")
	output := flags.String("o", "", "output API Reference path")
	if err := flags.Parse(arguments); err != nil {
		return err
	}

	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected positional arguments: %v", flags.Args())
	}

	if err := validateVersion(*version); err != nil {
		return err
	}

	if *output == "" {
		return fmt.Errorf("-o is required")
	}

	library := runtime.NewLibrary()
	if err := stdlib.Full().Register(library); err != nil {
		return fmt.Errorf("register stdlib.Full(): %w", err)
	}

	functions, err := library.Build()
	if err != nil {
		return fmt.Errorf("build stdlib.Full() registry: %w", err)
	}

	reference, err := analyzer.Generate(ctx, analyzer.Options{
		Root:      ".",
		Version:   *version,
		Functions: functions,
	})
	if err != nil {
		return err
	}

	return writeReference(*output, reference)
}

func validateVersion(value string) error {
	parsed, err := semver.StrictNewVersion(value)
	if err != nil || parsed.String() != value {
		return fmt.Errorf("-version must be an unprefixed canonical semantic version")
	}

	return nil
}

func writeReference(path string, reference any) error {
	directory := filepath.Dir(path)
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	temporary, err := os.CreateTemp(directory, ".api-reference-*.json")
	if err != nil {
		return fmt.Errorf("create temporary API Reference: %w", err)
	}

	temporaryPath := temporary.Name()
	removeTemporary := true
	defer func() {
		if removeTemporary {
			_ = os.Remove(temporaryPath)
		}
	}()

	if err := temporary.Chmod(0o644); err != nil {
		_ = temporary.Close()

		return fmt.Errorf("set API Reference permissions: %w", err)
	}

	encoder := json.NewEncoder(temporary)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(reference); err != nil {
		_ = temporary.Close()

		return fmt.Errorf("encode API Reference: %w", err)
	}

	if err := temporary.Sync(); err != nil {
		_ = temporary.Close()

		return fmt.Errorf("sync API Reference: %w", err)
	}

	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close API Reference: %w", err)
	}

	if err := os.Rename(temporaryPath, path); err != nil {
		return fmt.Errorf("replace API Reference: %w", err)
	}

	removeTemporary = false

	return nil
}
