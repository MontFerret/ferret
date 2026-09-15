package analyzer_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/MontFerret/specs/pkg/api"
	apicatalog "github.com/MontFerret/specs/pkg/api/catalog"
)

func assertRandomMetadata(t *testing.T, reference *api.Reference, catalog *apicatalog.Catalog) {
	t.Helper()

	want := map[string][]int{"float": {0, 2}, "int": {2}, "bool": {0}, "choice": {1}, "shuffle": {1}}
	seen := 0
	for _, namespace := range reference.Namespaces {
		for _, function := range namespace.Functions {
			if namespace.Name == "" && function.Name == "rand" {
				if len(function.Signatures) != 3 {
					t.Fatal("legacy rand lost fixed overloads")
				}

				for arity, signature := range function.Signatures {
					if len(signature.Parameters) != arity || !strings.Contains(signature.Deprecated, "random::") {
						t.Fatalf("legacy rand signature = %+v", signature)
					}

					if arity == 2 && (signature.Parameters[0].Name != "max" || signature.Parameters[1].Name != "min") {
						t.Fatal("legacy rand metadata changed argument order")
					}
				}
			}

			if namespace.Name != "random" {
				continue
			}

			seen++
			arities, ok := want[function.Name]
			if !ok || len(function.Signatures) != len(arities) {
				t.Fatalf("unexpected random function %+v", function)
			}

			for _, arity := range arities {
				if !hasSignature(function.Signatures, arity, false) {
					t.Fatalf("missing random::%s/%d", function.Name, arity)
				}
			}

			for _, signature := range function.Signatures {
				wantType := map[string]string{"float": "Float", "int": "Int", "bool": "Boolean", "choice": "Any", "shuffle": "List"}[function.Name]
				if signature.Return.Type.Name != wantType {
					t.Fatalf("random::%s return type = %+v, want %s", function.Name, signature.Return.Type, wantType)
				}

				if signature.Deprecated != "" {
					t.Fatalf("canonical random::%s is deprecated", function.Name)
				}

				if function.Name == "choice" || function.Name == "shuffle" {
					parameter := signature.Parameters[0]
					wantInput := "Any[]"
					if function.Name == "shuffle" {
						wantInput = "List"
					}

					if parameter.Name != "values" || parameter.Type.Name != wantInput {
						t.Fatalf("invalid random::%s input: %+v", function.Name, parameter)
					}
				}

				if len(signature.Parameters) == 2 && (signature.Parameters[0].Name != "min" || signature.Parameters[1].Name != "max") {
					t.Fatalf("canonical order is not min, max: %+v", signature)
				}
			}
		}
	}

	if seen != len(want) {
		t.Fatalf("random namespace has %d functions", seen)
	}

	found := false
	for _, category := range catalog.Categories {
		if category.ID == "random" {
			found = true
			wantFunctions := []apicatalog.FunctionRef{
				{Namespace: "random", Name: "bool"}, {Namespace: "random", Name: "choice"},
				{Namespace: "random", Name: "float"}, {Namespace: "random", Name: "int"},
				{Namespace: "random", Name: "shuffle"},
			}
			if category.Title != "Random" || category.Description != "Pseudo-random value generation functions in the random namespace." || !reflect.DeepEqual(category.Functions, wantFunctions) {
				t.Fatalf("random category = %+v", category)
			}
		}

		for _, function := range category.Functions {
			if function.Namespace == "" && function.Name == "rand" && category.ID != "math" {
				t.Fatal("legacy rand left the Math category")
			}
		}
	}

	if !found {
		t.Fatal("missing Random category")
	}
}
