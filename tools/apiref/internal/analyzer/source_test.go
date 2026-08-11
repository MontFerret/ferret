package analyzer

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadSourceCatalogNormalizesLineAndBlockComments(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, filepath.Join(root, "go.mod"), "module example.com/fixture\n\ngo 1.25.0\n")
	writeFixture(t, filepath.Join(root, "pkg", "stdlib", "sample", "sample.go"), `package sample

// LINE documents a line declaration.
// @param value {Any} Value.
// @return {Any} Value.
func Line(value any) any { return value }

/*
BLOCK documents a block declaration.
@return {None} None.
*/
func Block() {}
`)

	catalog, err := loadSourceCatalog(context.Background(), root)
	if err != nil {
		t.Fatal(err)
	}

	line := catalog.Declarations["example.com/fixture/pkg/stdlib/sample.Line"]
	block := catalog.Declarations["example.com/fixture/pkg/stdlib/sample.Block"]
	if line == nil || !strings.Contains(line.Documentation, "@param value {Any} Value.") {
		t.Fatalf("line documentation = %#v", line)
	}

	if block == nil || block.Documentation != "BLOCK documents a block declaration.\n@return {None} None." {
		t.Fatalf("block documentation = %#v", block)
	}
}

func writeFixture(t *testing.T, path, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
