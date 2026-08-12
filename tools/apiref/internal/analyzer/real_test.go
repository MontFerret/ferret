package analyzer_test

import (
	"context"
	"encoding/json"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/MontFerret/specs/pkg/api"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib"
	"github.com/MontFerret/ferret/v2/tools/apiref/internal/analyzer"
)

func TestGenerateMatchesFullRuntimeRegistryAndIsDeterministic(t *testing.T) {
	functions := fullFunctions(t)
	root, err := filepath.Abs("../../../..")
	if err != nil {
		t.Fatal(err)
	}

	options := analyzer.Options{Root: root, Version: "2.0.0-alpha.45", Functions: functions}
	first, err := analyzer.Generate(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}

	second, err := analyzer.Generate(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}

	firstBytes, err := json.MarshalIndent(first, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	secondBytes, err := json.MarshalIndent(second, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(firstBytes, secondBytes) {
		t.Fatal("repeated generation produced different bytes")
	}

	if first.ID != "montferret/core" || first.Version != options.Version {
		t.Fatalf("identity = %s@%s", first.ID, first.Version)
	}

	if err := api.Validate(first); err != nil {
		t.Fatalf("Specs validation: %v", err)
	}

	assertExactTopology(t, functions, first)
	assertSorted(t, first)
}

func fullFunctions(t *testing.T) *runtime.Functions {
	t.Helper()

	library := runtime.NewLibrary()
	if err := stdlib.Full().Register(library); err != nil {
		t.Fatal(err)
	}

	functions, err := library.Build()
	if err != nil {
		t.Fatal(err)
	}

	return functions
}

func assertExactTopology(t *testing.T, functions *runtime.Functions, reference *api.Reference) {
	t.Helper()

	topology := make(map[string][]api.Signature)
	count := 0
	for _, namespace := range reference.Namespaces {
		for _, function := range namespace.Functions {
			if function.Name != strings.ToLower(function.Name) {
				t.Fatalf("reference function %q in namespace %q is not canonical lowercase", function.Name, namespace.Name)
			}

			name := function.Name
			if namespace.Name != "" {
				name = namespace.Name + runtime.NamespaceSeparator + function.Name
			}

			topology[name] = function.Signatures
			count += len(function.Signatures)
		}
	}

	if count != functions.Size() {
		t.Fatalf("reference signatures = %d, runtime signatures = %d", count, functions.Size())
	}

	assertFixedCollection(t, topology, 0, functions.A0().Names())
	assertFixedCollection(t, topology, 1, functions.A1().Names())
	assertFixedCollection(t, topology, 2, functions.A2().Names())
	assertFixedCollection(t, topology, 3, functions.A3().Names())
	assertFixedCollection(t, topology, 4, functions.A4().Names())

	for _, name := range functions.Var().Names() {
		if !hasSignature(topology[name], 0, true) {
			t.Fatalf("reference is missing variadic runtime signature %s", name)
		}
	}
}

func assertFixedCollection(t *testing.T, topology map[string][]api.Signature, arity int, names []string) {
	t.Helper()

	for _, name := range names {
		if !hasSignature(topology[name], arity, false) {
			t.Fatalf("reference is missing runtime signature %s/%d", name, arity)
		}
	}
}

func hasSignature(signatures []api.Signature, arity int, variadic bool) bool {
	for _, signature := range signatures {
		if signature.Variadic == variadic && (variadic || len(signature.Parameters) == arity) {
			return true
		}
	}

	return false
}

func assertSorted(t *testing.T, reference *api.Reference) {
	t.Helper()

	namespaceNames := make([]string, 0, len(reference.Namespaces))
	for _, namespace := range reference.Namespaces {
		namespaceNames = append(namespaceNames, namespace.Name)
		functionNames := make([]string, 0, len(namespace.Functions))
		for _, function := range namespace.Functions {
			functionNames = append(functionNames, function.Name)

			for index := 1; index < len(function.Signatures); index++ {
				previous := function.Signatures[index-1]
				current := function.Signatures[index]
				if previous.Variadic || !current.Variadic && len(previous.Parameters) > len(current.Parameters) {
					t.Fatalf("signatures for %s are not fixed-arity ascending then variadic", function.Name)
				}
			}
		}

		if !sort.StringsAreSorted(functionNames) {
			t.Fatalf("functions in namespace %q are not sorted: %s", namespace.Name, strings.Join(functionNames, ", "))
		}
	}

	if !sort.StringsAreSorted(namespaceNames) {
		t.Fatalf("namespaces are not sorted: %s", strings.Join(namespaceNames, ", "))
	}
}
