package analyzer_test

import (
	"reflect"
	"testing"

	"github.com/MontFerret/specs/pkg/api"
	apicatalog "github.com/MontFerret/specs/pkg/api/catalog"
)

func assertCollectionMetadata(t *testing.T, reference *api.Reference, catalog *apicatalog.Catalog) {
	t.Helper()

	want := map[string][]string{
		"count":          {"Iterable"},
		"count_distinct": {"Iterable"},
		"includes":       {"String", "Containable", "Iterable"},
		"reverse":        {"String", "List"},
	}

	signatures := make(map[apicatalog.FunctionRef][]api.Signature)
	for _, namespace := range reference.Namespaces {
		for _, function := range namespace.Functions {
			signatures[apicatalog.FunctionRef{Namespace: namespace.Name, Name: function.Name}] = function.Signatures
		}
	}

	for _, namespace := range reference.Namespaces {
		if namespace.Name != "collections" {
			continue
		}

		for _, function := range namespace.Functions {
			expected, ok := want[function.Name]
			if !ok {
				continue
			}

			if len(function.Signatures) != 1 {
				t.Fatalf("%s overloads: %v", function.Name, function.Signatures)
			}

			signature := function.Signatures[0]
			if signature.Deprecated != "" {
				t.Fatalf("canonical collections::%s is deprecated", function.Name)
			}

			legacy := signatures[apicatalog.FunctionRef{Name: function.Name}]
			if len(legacy) != 1 || legacy[0].Deprecated != "Use collections::"+function.Name+" instead." {
				t.Fatalf("incorrect global compatibility metadata for %s: %+v", function.Name, legacy)
			}

			comparable := legacy[0]
			comparable.Deprecated = ""
			comparable.Description = signature.Description
			if !reflect.DeepEqual(comparable, signature) {
				t.Fatalf("global and canonical signatures differ for %s", function.Name)
			}

			arity := 1
			if function.Name == "includes" {
				arity = 2
			}

			if len(signature.Parameters) != arity || signature.Variadic {
				t.Fatalf("%s arity changed", function.Name)
			}

			parameter := signature.Parameters[0].Type
			actual := []string{parameter.Name}
			if parameter.Kind == api.TypeKindUnion {
				actual = nil
				for _, member := range parameter.Types {
					actual = append(actual, member.Name)
				}
			}

			if !reflect.DeepEqual(actual, expected) {
				t.Fatalf("%s accepts %v; want %v", function.Name, actual, expected)
			}

			if function.Name == "reverse" {
				actual = nil
				for _, member := range signature.Return.Type.Types {
					actual = append(actual, member.Name)
				}

				if !reflect.DeepEqual(actual, expected) {
					t.Fatalf("reverse return type: %v", actual)
				}
			} else {
				result := "Int"
				if function.Name == "includes" {
					result = "Boolean"
				}

				if signature.Return.Type.Name != result {
					t.Fatalf("%s return type changed", function.Name)
				}
			}

			delete(want, function.Name)
		}
	}

	if len(want) != 0 {
		t.Fatalf("missing functions: %v", want)
	}

	for _, category := range catalog.Categories {
		if category.ID != "collections" {
			continue
		}

		if category.Title != "Collections" || category.Description != "Functions for working with collections and collection values." {
			t.Fatalf("incorrect Collections category: %+v", category)
		}

		expected := []apicatalog.FunctionRef{
			{Name: "count"}, {Name: "count_distinct"}, {Name: "includes"}, {Name: "reverse"},
			{Namespace: "collections", Name: "count"}, {Namespace: "collections", Name: "count_distinct"},
			{Namespace: "collections", Name: "includes"}, {Namespace: "collections", Name: "reverse"},
		}
		if !reflect.DeepEqual(category.Functions, expected) {
			t.Fatalf("Collections catalog = %v, want %v", category.Functions, expected)
		}

		return
	}

	t.Fatal("missing Collections catalog category")
}
