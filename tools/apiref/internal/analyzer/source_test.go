package analyzer

import (
	"context"
	"fmt"
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

func TestLoadSourceCatalogExpandsLiteralAssertionCatalog(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, filepath.Join(root, "go.mod"), "module example.com/fixture\n\ngo 1.25.0\n")
	writeFixture(t, filepath.Join(root, "pkg", "stdlib", "testing", "testing.go"), assertionCatalogFixture(`
	{name: "value", descriptor: unaryAssertion, negatable: true},
	{name: "fail", descriptor: failAssertion, negatable: false},
`))

	catalog, err := loadSourceCatalog(context.Background(), root)
	if err != nil {
		t.Fatal(err)
	}

	want := map[string][2]int{
		"t::value":      {1, 2},
		"t::not::value": {1, 2},
		"t::fail":       {0, 1},
	}
	if len(catalog.Assertions) != len(want) {
		t.Fatalf("assertions = %#v, want %d entries", catalog.Assertions, len(want))
	}

	for name, bounds := range want {
		descriptor, exists := catalog.Assertions[name]
		if !exists {
			t.Errorf("assertion catalog is missing %s", name)
			continue
		}

		if descriptor.Min != bounds[0] || descriptor.Max != bounds[1] {
			t.Errorf("%s bounds = %d..%d, want %d..%d", name, descriptor.Min, descriptor.Max, bounds[0], bounds[1])
		}
	}

	if _, exists := catalog.Assertions["t::not::fail"]; exists {
		t.Fatal("non-negatable assertion t::fail was expanded into t::not::fail")
	}
}

func TestLoadSourceCatalogRejectsMalformedAssertionCatalog(t *testing.T) {
	tests := []struct {
		name    string
		entries string
		want    string
	}{
		{
			name:    "dynamic name",
			entries: `{name: dynamicName, descriptor: unaryAssertion, negatable: true},`,
			want:    "name must be a string literal",
		},
		{
			name:    "dynamic negatable",
			entries: `{name: "value", descriptor: unaryAssertion, negatable: dynamicNegatable},`,
			want:    "negatable must be a boolean literal",
		},
		{
			name:    "dynamic descriptor",
			entries: `{name: "value", descriptor: newUnaryAssertion(), negatable: true},`,
			want:    "must use a statically named descriptor",
		},
		{
			name:    "missing metadata",
			entries: `{name: "value", descriptor: unaryAssertion},`,
			want:    "must declare negatable",
		},
		{
			name:    "unkeyed metadata",
			entries: `{"value", unaryAssertion, true},`,
			want:    "must use keyed fields",
		},
		{
			name: "duplicate name",
			entries: `
	{name: "value", descriptor: unaryAssertion, negatable: true},
	{name: "value", descriptor: unaryAssertion, negatable: false},
`,
			want: "is registered more than once",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			writeFixture(t, filepath.Join(root, "go.mod"), "module example.com/fixture\n\ngo 1.25.0\n")
			writeFixture(
				t,
				filepath.Join(root, "pkg", "stdlib", "testing", "testing.go"),
				assertionCatalogFixture(test.entries),
			)

			_, err := loadSourceCatalog(context.Background(), root)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("load error = %v, want %q", err, test.want)
			}
		})
	}
}

func assertionCatalogFixture(entries string) string {
	return fmt.Sprintf(`package testing

type (
	assertionArgs struct {
		min int
		max int
	}

	assertion struct {
		args assertionArgs
	}

	assertionRegistration struct {
		name string
		descriptor assertion
		negatable bool
	}
)

var (
	dynamicName = "value"
	dynamicNegatable = true
	unaryAssertion = newUnaryAssertion()
	failAssertion = assertion{args: assertionArgs{min: 0, max: 1}}
)

var assertionCatalog = []assertionRegistration{%s}

func newUnaryAssertion() assertion {
	return assertion{args: assertionArgs{min: 1, max: 2}}
}
`, entries)
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
