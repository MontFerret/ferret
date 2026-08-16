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
	catalogOutput := flags.String("catalog", "", "output Standard Library catalog path")
	if err := flags.Parse(arguments); err != nil {
		return err
	}

	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected positional arguments: %v", flags.Args())
	}

	if err := validateVersion(*version); err != nil {
		return err
	}

	if *output == "" || *catalogOutput == "" {
		return fmt.Errorf("-o and -catalog are required")
	}

	if filepath.Clean(*output) == filepath.Clean(*catalogOutput) {
		return fmt.Errorf("-o and -catalog must identify different files")
	}

	library := runtime.NewLibrary()
	if err := stdlib.Full().Register(library); err != nil {
		return fmt.Errorf("register stdlib.Full(): %w", err)
	}

	functions, err := library.Build()
	if err != nil {
		return fmt.Errorf("build stdlib.Full() registry: %w", err)
	}

	artifacts, err := analyzer.Generate(ctx, analyzer.Options{
		Root:      filepath.Join("..", ".."),
		Version:   *version,
		Functions: functions,
	})
	if err != nil {
		return err
	}

	if err := writeArtifact(*output, "API Reference", artifacts.Reference); err != nil {
		return err
	}

	return writeArtifact(*catalogOutput, "Standard Library catalog", artifacts.Catalog)
}

func validateVersion(value string) error {
	parsed, err := semver.StrictNewVersion(value)
	if err != nil || parsed.String() != value {
		return fmt.Errorf("-version must be an unprefixed canonical semantic version")
	}

	return nil
}

func writeArtifact(path, label string, value any) error {
	directory := filepath.Dir(path)
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	temporary, err := os.CreateTemp(directory, ".stdlib-artifact-*.json")
	if err != nil {
		return fmt.Errorf("create temporary %s: %w", label, err)
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

		return fmt.Errorf("set %s permissions: %w", label, err)
	}

	encoder := json.NewEncoder(temporary)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(value); err != nil {
		_ = temporary.Close()

		return fmt.Errorf("encode %s: %w", label, err)
	}

	if err := temporary.Sync(); err != nil {
		_ = temporary.Close()

		return fmt.Errorf("sync %s: %w", label, err)
	}

	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close %s: %w", label, err)
	}

	if err := os.Rename(temporaryPath, path); err != nil {
		return fmt.Errorf("replace %s: %w", label, err)
	}

	removeTemporary = false

	return nil
}
