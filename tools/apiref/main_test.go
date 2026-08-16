package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateVersionRequiresCanonicalUnprefixedSemVer(t *testing.T) {
	for _, version := range []string{"2.0.0-alpha.45", "1.2.3", "1.2.3+build.1"} {
		if err := validateVersion(version); err != nil {
			t.Fatalf("validate %q: %v", version, err)
		}
	}

	for _, version := range []string{"", "v1.2.3", "1.2", "01.2.3", "1.2.3-alpha.01"} {
		if err := validateVersion(version); err == nil {
			t.Fatalf("expected %q to be rejected", version)
		}
	}
}

func TestWriteArtifactIsIndentedNewlineTerminatedAndReplaceable(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "api.json")
	if err := writeArtifact(path, "test artifact", map[string]any{"version": "1.0.0", "schemaVersion": 1}); err != nil {
		t.Fatal(err)
	}

	first, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasSuffix(string(first), "\n") || !strings.Contains(string(first), "\n  \"schemaVersion\"") {
		t.Fatalf("output is not indented with one trailing newline:\n%s", first)
	}

	if err := writeArtifact(path, "test artifact", map[string]any{"schemaVersion": 2}); err != nil {
		t.Fatalf("replace output: %v", err)
	}

	second, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(string(second), "1.0.0") || strings.Count(string(second), "\n") != 3 {
		t.Fatalf("replacement retained old or trailing content:\n%s", second)
	}
}

func TestRunRequiresDistinctArtifactPaths(t *testing.T) {
	err := run(context.Background(), []string{"-version", "2.0.0-alpha.47", "-o", "same.json", "-catalog", "same.json"})
	if err == nil || !strings.Contains(err.Error(), "different files") {
		t.Fatalf("run error = %v, want distinct path rejection", err)
	}
}
