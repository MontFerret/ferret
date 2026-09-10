package analyzer_test

import (
	"testing"

	"github.com/MontFerret/specs/pkg/api"
	apicatalog "github.com/MontFerret/specs/pkg/api/catalog"
)

func assertArrayMetadata(t *testing.T, reference *api.Reference, catalog *apicatalog.Catalog) {
	t.Helper()

	want := map[string][]int{
		"first": {1}, "last": {1}, "at": {2}, "append": {2}, "concat": {-1},
		"flatten": {1, 2}, "slice": {2, 3}, "contains": {2}, "index_of": {2},
		"unique": {1}, "sorted": {1}, "remove": {2}, "remove_at": {2}, "remove_any": {2},
		"union": {-1}, "intersection": {-1}, "difference": {-1}, "symmetric_difference": {-1},
	}

	seen := make(map[string]bool)
	deprecated := make(map[string]bool)
	for _, namespace := range reference.Namespaces {
		if namespace.Name == "arrays::mut" {
			t.Fatal("array mutation must not be published")
		}

		for _, function := range namespace.Functions {
			if namespace.Name == "" {
				for _, signature := range function.Signatures {
					if signature.Deprecated != "" {
						deprecated[function.Name] = true
					}
				}
			}

			if namespace.Name != "arrays" {
				continue
			}

			arities, ok := want[function.Name]
			if !ok || len(arities) != len(function.Signatures) {
				t.Fatalf("unexpected canonical array signatures: %s %+v", function.Name, function.Signatures)
			}

			seen[function.Name] = true
			for index, signature := range function.Signatures {
				if !hasSignature([]api.Signature{signature}, arities[index], arities[index] == -1) || signature.Deprecated != "" {
					t.Fatalf("invalid canonical array signature: %s %+v", function.Name, signature)
				}

				resultType := signature.Return.Type
				switch function.Name {
				case "first", "last", "at":
					if resultType.Kind != api.TypeKindNamed || resultType.Name != "Any" {
						t.Fatalf("incorrect element return type: %s %+v", function.Name, resultType)
					}
				case "contains", "index_of":
					name := "Boolean"
					if function.Name == "index_of" {
						name = "Int"
					}

					if resultType.Kind != api.TypeKindNamed || resultType.Name != name {
						t.Fatalf("incorrect search return type: %s %+v", function.Name, resultType)
					}
				default:
					if resultType.Kind != api.TypeKindNamed || resultType.Name != "Any[]" {
						t.Fatalf("incorrect array return type: %s %+v", function.Name, resultType)
					}
				}
			}
		}
	}

	if len(seen) != len(want) {
		t.Fatalf("canonical array metadata incomplete: %v", seen)
	}

	for _, name := range []string{"append", "push", "position", "remove_value", "remove_nth", "outersection", "sorted_unique", "pop", "shift", "unshift"} {
		if !deprecated[name] {
			t.Errorf("legacy adapter %s lacks deprecation metadata", name)
		}
	}

	categorized := 0
	for _, category := range catalog.Categories {
		for _, function := range category.Functions {
			if function.Namespace == "arrays" {
				if category.ID != "arrays" {
					t.Fatalf("array function categorized under %s", category.ID)
				}

				categorized++
			}
		}
	}

	if categorized != len(want) {
		t.Fatalf("categorized %d canonical arrays, want %d", categorized, len(want))
	}
}
