package analyzer

import (
	"reflect"
	"strings"
	"testing"

	"github.com/MontFerret/specs/pkg/api"
	apicatalog "github.com/MontFerret/specs/pkg/api/catalog"
)

func TestBuildStandardLibraryCatalogDerivesCategoriesAndQualifiedOverrides(t *testing.T) {
	metadata := []apicatalog.Category{
		{ID: "arrays", Title: "Arrays", Description: "Array functions."},
		{ID: "io", Title: "I/O", Description: "I/O functions."},
		{ID: "math", Title: "Math", Description: "Math functions."},
		{ID: "testing", Title: "Testing", Description: "Testing functions."},
	}
	registered := []registeredSignature{
		{QualifiedName: "append", Name: "append", PackagePath: "github.com/MontFerret/ferret/v2/pkg/stdlib/arrays"},
		{QualifiedName: "abs", Name: "abs", PackagePath: "github.com/MontFerret/ferret/v2/pkg/stdlib/utils"},
		{QualifiedName: "abs", Name: "abs", PackagePath: "github.com/MontFerret/ferret/v2/pkg/stdlib/math"},
		{QualifiedName: "io::fs::read", Namespace: "io::fs", Name: "read", PackagePath: "github.com/MontFerret/ferret/v2/pkg/stdlib/io/fs"},
		{QualifiedName: "t::read", Namespace: "t", Name: "read", PackagePath: "github.com/MontFerret/ferret/v2/pkg/stdlib/testing/base"},
	}
	overrides := map[functionIdentity]string{{Name: "abs"}: "math"}

	catalog, err := buildStandardLibraryCatalog("2.0.0-alpha.47", registered, metadata, overrides)
	if err != nil {
		t.Fatal(err)
	}

	want := [][]apicatalog.FunctionRef{
		{{Namespace: "", Name: "append"}},
		{{Namespace: "io::fs", Name: "read"}},
		{{Namespace: "", Name: "abs"}},
		{{Namespace: "t", Name: "read"}},
	}
	for index := range want {
		if got := catalog.Categories[index].Functions; !reflect.DeepEqual(got, want[index]) {
			t.Fatalf("category %q functions = %v, want %v", catalog.Categories[index].ID, got, want[index])
		}
	}
}

func TestBuildStandardLibraryCatalogRejectsInvalidMembership(t *testing.T) {
	metadata := []apicatalog.Category{
		{ID: "arrays", Title: "Arrays", Description: "Array functions."},
		{ID: "math", Title: "Math", Description: "Math functions."},
	}

	tests := []struct {
		overrides  map[functionIdentity]string
		name       string
		want       string
		registered []registeredSignature
	}{
		{name: "outside stdlib", registered: []registeredSignature{{QualifiedName: "abs", Name: "abs", PackagePath: "example/math"}}, want: "outside pkg/stdlib"},
		{name: "unknown category", registered: []registeredSignature{{QualifiedName: "abs", Name: "abs", PackagePath: "example/pkg/stdlib/unknown"}}, want: "unknown category"},
		{name: "conflicting overloads", registered: []registeredSignature{
			{QualifiedName: "abs", Name: "abs", PackagePath: "example/pkg/stdlib/math"},
			{QualifiedName: "abs", Name: "abs", PackagePath: "example/pkg/stdlib/arrays"},
		}, overrides: map[functionIdentity]string{}, want: "overloads resolve"},
		{name: "empty category", registered: []registeredSignature{{QualifiedName: "abs", Name: "abs", PackagePath: "example/pkg/stdlib/math"}}, want: `category "arrays" has no functions`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			catalog, err := buildStandardLibraryCatalog("1.0.0", test.registered, metadata, test.overrides)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("build error = %v, catalog = %#v, want %q", err, catalog, test.want)
			}
		})
	}
}

func TestValidateCatalogAgainstReference(t *testing.T) {
	reference := &api.Reference{
		ID:      moduleID,
		Version: "1.2.3",
		Namespaces: []api.Namespace{
			{Name: "", Functions: []api.Function{{Name: "read"}}},
			{Name: "io::fs", Functions: []api.Function{{Name: "read"}}},
			{Name: "t", Functions: []api.Function{{Name: "read"}}},
		},
	}
	catalog := &apicatalog.Catalog{
		ID:      moduleID,
		Version: "1.2.3",
		Categories: []apicatalog.Category{
			{ID: "io", Functions: []apicatalog.FunctionRef{{Namespace: "io::fs", Name: "read"}}},
			{ID: "testing", Functions: []apicatalog.FunctionRef{{Namespace: "t", Name: "read"}}},
			{ID: "utils", Functions: []apicatalog.FunctionRef{{Namespace: "", Name: "read"}}},
		},
	}

	if err := validateCatalogAgainstReference(reference, catalog); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		mutate func(*apicatalog.Catalog)
		want   string
	}{
		{name: "id", mutate: func(value *apicatalog.Catalog) { value.ID = "other/module" }, want: "does not match API id"},
		{name: "version", mutate: func(value *apicatalog.Catalog) { value.Version = "1.2.4" }, want: "does not match API version"},
		{name: "unknown namespace", mutate: func(value *apicatalog.Catalog) { value.Categories[0].Functions[0].Namespace = "io::missing" }, want: "unknown API namespace"},
		{name: "unknown function", mutate: func(value *apicatalog.Catalog) { value.Categories[0].Functions[0].Name = "missing" }, want: "unknown function"},
		{name: "duplicate", mutate: func(value *apicatalog.Catalog) { value.Categories[1].Functions[0] = value.Categories[0].Functions[0] }, want: "assigned to categories"},
		{name: "uncategorized global", mutate: func(value *apicatalog.Catalog) { value.Categories[2].Functions = nil }, want: `function "read" is not assigned`},
		{name: "uncategorized namespaced", mutate: func(value *apicatalog.Catalog) { value.Categories[1].Functions = nil }, want: `function "t::read" is not assigned`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			copyValue := *catalog
			copyValue.Categories = append([]apicatalog.Category(nil), catalog.Categories...)
			for index := range copyValue.Categories {
				copyValue.Categories[index].Functions = append([]apicatalog.FunctionRef(nil), catalog.Categories[index].Functions...)
			}
			test.mutate(&copyValue)

			err := validateCatalogAgainstReference(reference, &copyValue)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("validation error = %v, want %q", err, test.want)
			}
		})
	}
}
