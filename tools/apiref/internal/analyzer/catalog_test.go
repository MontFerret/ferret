package analyzer

import (
	"reflect"
	"strings"
	"testing"

	"github.com/MontFerret/specs/pkg/api"
	apicatalog "github.com/MontFerret/specs/pkg/api/catalog"
)

func TestBuildStandardLibraryCatalogDerivesCategoriesOverridesAndRoots(t *testing.T) {
	metadata := []apicatalog.Category{
		{ID: "arrays", Title: "Arrays", Description: "Array functions."},
		{ID: "math", Title: "Math", Description: "Math functions."},
	}
	reference := &api.Reference{
		SchemaVersion: api.SchemaVersion,
		ID:            moduleID,
		Version:       "2.0.0-alpha.47",
		Namespaces: []api.Namespace{
			{Name: "", Functions: []api.Function{{Name: "abs"}, {Name: "append"}}},
			{Name: "io::fs", Functions: []api.Function{{Name: "read"}}},
			{Name: "io::net::http", Functions: []api.Function{{Name: "get"}}},
			{Name: "t", Functions: []api.Function{{Name: "eq"}}},
		},
	}
	registered := []registeredSignature{
		{Name: "append", PackagePath: "github.com/MontFerret/ferret/v2/pkg/stdlib/arrays"},
		{Name: "abs", PackagePath: "github.com/MontFerret/ferret/v2/pkg/stdlib/utils"},
		{Name: "abs", PackagePath: "github.com/MontFerret/ferret/v2/pkg/stdlib/math"},
		{Name: "read", Namespace: "io::fs", PackagePath: "github.com/MontFerret/ferret/v2/pkg/stdlib/io"},
	}

	catalog, err := buildStandardLibraryCatalog("2.0.0-alpha.47", registered, reference, metadata, map[string]string{"abs": "math"})
	if err != nil {
		t.Fatal(err)
	}

	if got, want := catalog.Categories[0].Functions, []string{"append"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("arrays = %v, want %v", got, want)
	}

	if got, want := catalog.Categories[1].Functions, []string{"abs"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("math = %v, want %v", got, want)
	}

	if got, want := catalog.NamespaceRoots, []string{"io", "t"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("roots = %v, want %v", got, want)
	}
}

func TestBuildStandardLibraryCatalogRejectsInvalidMembership(t *testing.T) {
	metadata := []apicatalog.Category{
		{ID: "arrays", Title: "Arrays", Description: "Array functions."},
		{ID: "math", Title: "Math", Description: "Math functions."},
	}
	reference := &api.Reference{ID: moduleID, Version: "1.0.0"}

	tests := []struct {
		overrides  map[string]string
		name       string
		want       string
		registered []registeredSignature
	}{
		{name: "outside stdlib", registered: []registeredSignature{{Name: "abs", PackagePath: "example/math"}}, want: "outside pkg/stdlib"},
		{name: "unknown category", registered: []registeredSignature{{Name: "abs", PackagePath: "example/pkg/stdlib/unknown"}}, want: "unknown category"},
		{name: "conflicting overloads", registered: []registeredSignature{
			{Name: "abs", PackagePath: "example/pkg/stdlib/math"},
			{Name: "abs", PackagePath: "example/pkg/stdlib/arrays"},
		}, overrides: map[string]string{}, want: "overloads resolve"},
		{name: "empty category", registered: []registeredSignature{{Name: "abs", PackagePath: "example/pkg/stdlib/math"}}, want: `category "arrays" has no global functions`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			catalog, err := buildStandardLibraryCatalog("1.0.0", test.registered, reference, metadata, test.overrides)
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
			{Name: "", Functions: []api.Function{{Name: "abs"}, {Name: "append"}}},
			{Name: "io::fs", Functions: []api.Function{{Name: "read"}}},
		},
	}
	catalog := &apicatalog.Catalog{
		ID:      moduleID,
		Version: "1.2.3",
		Categories: []apicatalog.Category{
			{ID: "arrays", Functions: []string{"append"}},
			{ID: "math", Functions: []string{"abs"}},
		},
		NamespaceRoots: []string{"io"},
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
		{name: "unknown", mutate: func(value *apicatalog.Catalog) { value.Categories[0].Functions[0] = "missing" }, want: "unknown global function"},
		{name: "duplicate", mutate: func(value *apicatalog.Catalog) { value.Categories[1].Functions[0] = "append" }, want: "assigned to categories"},
		{name: "uncategorized", mutate: func(value *apicatalog.Catalog) { value.Categories[0].Functions = nil }, want: "not assigned"},
		{name: "missing root", mutate: func(value *apicatalog.Catalog) { value.NamespaceRoots = nil }, want: "not declared"},
		{name: "unknown root", mutate: func(value *apicatalog.Catalog) { value.NamespaceRoots = append(value.NamespaceRoots, "t") }, want: "does not cover"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			copyValue := *catalog
			copyValue.Categories = append([]apicatalog.Category(nil), catalog.Categories...)
			for index := range copyValue.Categories {
				copyValue.Categories[index].Functions = append([]string(nil), catalog.Categories[index].Functions...)
			}
			copyValue.NamespaceRoots = append([]string(nil), catalog.NamespaceRoots...)
			test.mutate(&copyValue)

			err := validateCatalogAgainstReference(reference, &copyValue)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("validation error = %v, want %q", err, test.want)
			}
		})
	}
}
