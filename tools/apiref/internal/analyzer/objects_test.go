package analyzer_test

import (
	"reflect"
	"testing"

	"github.com/MontFerret/specs/pkg/api"
	apicatalog "github.com/MontFerret/specs/pkg/api/catalog"
)

func assertObjectCompatibilityMetadata(t *testing.T, reference *api.Reference, catalog *apicatalog.Catalog) {
	t.Helper()

	signatures := make(map[apicatalog.FunctionRef][]api.Signature)
	for _, namespace := range reference.Namespaces {
		for _, function := range namespace.Functions {
			signatures[apicatalog.FunctionRef{Namespace: namespace.Name, Name: function.Name}] = function.Signatures
			if namespace.Name != "object" && namespace.Name != "object::mut" {
				continue
			}

			for _, signature := range function.Signatures {
				if signature.Deprecated != "" {
					t.Fatalf("canonical function %s::%s is deprecated", namespace.Name, function.Name)
				}
			}
		}
	}

	categories := make(map[apicatalog.FunctionRef]string)
	for _, category := range catalog.Categories {
		for _, function := range category.Functions {
			categories[function] = category.ID
		}
	}

	for legacy, canonical := range map[string]string{
		"keys": "keys", "values": "values", "has": "has_key", "keep_keys": "keep_keys",
		"merge": "merge", "merge_recursive": "merge_deep", "zip": "zip",
	} {
		global := apicatalog.FunctionRef{Name: legacy}
		target := apicatalog.FunctionRef{Namespace: "object", Name: canonical}
		old, current := signatures[global], signatures[target]
		if len(old) != 1 || len(current) != 1 {
			t.Fatalf("%s and object::%s must each have exactly one signature", legacy, canonical)
		}

		want := "Use object::" + canonical + " instead."
		if old[0].Deprecated != want {
			t.Fatalf("%s deprecation = %q, want %q", legacy, old[0].Deprecated, want)
		}

		comparable := old[0]
		comparable.Deprecated = ""
		comparable.Description = current[0].Description
		if !reflect.DeepEqual(comparable, current[0]) {
			t.Fatalf("%s signature differs from object::%s: %+v != %+v", legacy, canonical, comparable, current[0])
		}

		if categories[global] != "objects" || categories[target] != "objects" {
			t.Fatalf("%s and object::%s must belong to the Objects catalog", legacy, canonical)
		}
	}
}
