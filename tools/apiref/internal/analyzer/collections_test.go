package analyzer_test

import (
	"reflect"
	"testing"

	"github.com/MontFerret/specs/pkg/api"
)

func assertCollectionMetadata(t *testing.T, reference *api.Reference) {
	t.Helper()

	want := map[string][]string{
		"count":          {"Iterable"},
		"count_distinct": {"Iterable"},
		"includes":       {"String", "Containable", "Iterable"},
		"reverse":        {"String", "List"},
	}

	for _, namespace := range reference.Namespaces {
		if namespace.Name != "" {
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
}
