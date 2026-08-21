package testing

import (
	"context"
	"slices"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type expectedAssertionRegistration struct {
	name      string
	minimum   int
	maximum   int
	negatable bool
}

func TestAssertionCatalogDefinesExactRuntimeTopology(t *testing.T) {
	t.Parallel()

	expected := []expectedAssertionRegistration{
		{name: "approx", minimum: 3, maximum: 4, negatable: true},
		{name: "between", minimum: 3, maximum: 4, negatable: true},
		{name: "bool", minimum: 1, maximum: 2, negatable: true},
		{name: "contains", minimum: 2, maximum: 3, negatable: true},
		{name: "duration", minimum: 1, maximum: 2, negatable: true},
		{name: "empty", minimum: 1, maximum: 2, negatable: true},
		{name: "eq", minimum: 2, maximum: 3, negatable: true},
		{name: "fail", minimum: 0, maximum: 1, negatable: false},
		{name: "false", minimum: 1, maximum: 2, negatable: true},
		{name: "gt", minimum: 2, maximum: 3, negatable: true},
		{name: "gte", minimum: 2, maximum: 3, negatable: true},
		{name: "has", minimum: 2, maximum: 3, negatable: true},
		{name: "len", minimum: 2, maximum: 3, negatable: true},
		{name: "match", minimum: 2, maximum: 3, negatable: true},
		{name: "lt", minimum: 2, maximum: 3, negatable: true},
		{name: "lte", minimum: 2, maximum: 3, negatable: true},
		{name: "none", minimum: 1, maximum: 2, negatable: true},
		{name: "number", minimum: 1, maximum: 2, negatable: true},
		{name: "true", minimum: 1, maximum: 2, negatable: true},
		{name: "string", minimum: 1, maximum: 2, negatable: true},
		{name: "int", minimum: 1, maximum: 2, negatable: true},
		{name: "float", minimum: 1, maximum: 2, negatable: true},
		{name: "datetime", minimum: 1, maximum: 2, negatable: true},
		{name: "array", minimum: 1, maximum: 2, negatable: true},
		{name: "object", minimum: 1, maximum: 2, negatable: true},
		{name: "binary", minimum: 1, maximum: 2, negatable: true},
	}

	if len(assertionCatalog) != len(expected) {
		t.Fatalf("assertionCatalog has %d entries, want %d", len(assertionCatalog), len(expected))
	}

	library := runtime.NewLibrary()
	RegisterLib(library)

	functions, err := library.Build()
	if err != nil {
		t.Fatal(err)
	}

	expectedNames := make([]string, 0, 51)
	expectedDefinitions := 0
	for index, want := range expected {
		registration := assertionCatalog[index]
		if registration.name != want.name ||
			registration.descriptor.args.min != want.minimum ||
			registration.descriptor.args.max != want.maximum ||
			registration.negatable != want.negatable {
			t.Errorf("assertionCatalog[%d] = %#v, want %#v", index, registration, want)
		}

		positiveName := "t::" + want.name
		expectedNames = append(expectedNames, positiveName)
		expectedDefinitions += want.maximum - want.minimum + 1
		assertArities(t, functions, positiveName, want.minimum, want.maximum)

		if want.negatable {
			negativeName := "t::not::" + want.name
			expectedNames = append(expectedNames, negativeName)
			expectedDefinitions += want.maximum - want.minimum + 1
			assertArities(t, functions, negativeName, want.minimum, want.maximum)
		} else if functions.Has("t::not::" + want.name) {
			t.Errorf("non-negatable assertion %s is registered under t::not", want.name)
		}
	}

	slices.Sort(expectedNames)
	if actual := functions.List(); !slices.Equal(actual, expectedNames) {
		t.Fatalf("registered assertion names = %v, want %v", actual, expectedNames)
	}
	if functions.Size() != expectedDefinitions {
		t.Fatalf("registered definitions = %d, want %d", functions.Size(), expectedDefinitions)
	}

	for _, removed := range []string{"t::include", "t::not::include"} {
		if functions.Has(removed) {
			t.Errorf("removed assertion %s remains registered", removed)
		}
	}
}

func TestRegisteredAssertionFunctionsPreservePolarity(t *testing.T) {
	t.Parallel()

	library := runtime.NewLibrary()
	RegisterLib(library)

	functions, err := library.Build()
	if err != nil {
		t.Fatal(err)
	}

	positive, ok := functions.A2().Get("t::eq")
	if !ok {
		t.Fatal("t::eq/2 is not registered")
	}
	negative, ok := functions.A2().Get("t::not::eq")
	if !ok {
		t.Fatal("t::not::eq/2 is not registered")
	}

	out, err := positive(context.Background(), runtime.NewInt(1), runtime.NewInt(1))
	if err != nil || out != runtime.None {
		t.Fatalf("t::eq success = (%v, %v), want (None, nil)", out, err)
	}

	out, err = negative(context.Background(), runtime.NewInt(1), runtime.NewInt(2))
	if err != nil || out != runtime.None {
		t.Fatalf("t::not::eq success = (%v, %v), want (None, nil)", out, err)
	}
}

func assertArities(t *testing.T, functions *runtime.Functions, name string, minimum, maximum int) {
	t.Helper()

	if functions.Var().Has(name) {
		t.Errorf("%s unexpectedly has a variadic registration", name)
	}

	for arity := 0; arity <= 4; arity++ {
		want := arity >= minimum && arity <= maximum
		if actual := hasAssertionArity(functions, name, arity); actual != want {
			t.Errorf("%s/%d registered = %t, want %t", name, arity, actual, want)
		}
	}
}

func hasAssertionArity(functions *runtime.Functions, name string, arity int) bool {
	switch arity {
	case 0:
		return functions.A0().Has(name)
	case 1:
		return functions.A1().Has(name)
	case 2:
		return functions.A2().Has(name)
	case 3:
		return functions.A3().Has(name)
	case 4:
		return functions.A4().Has(name)
	default:
		return false
	}
}
