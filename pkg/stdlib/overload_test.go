package stdlib_test

import (
	"slices"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib"
)

func TestBoundedFunctionsUseFixedArityRegistrations(t *testing.T) {
	t.Parallel()

	functions := buildFunctions(t, stdlib.Full())
	expected := boundedFunctionArities()

	if len(expected) != 81 {
		t.Fatalf("bounded function matrix has %d names, want 81", len(expected))
	}

	for name, arities := range expected {
		name, arities := name, arities
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if functions.Var().Has(name) {
				t.Fatalf("%s unexpectedly has a variadic fallback", name)
			}

			for arity := 0; arity <= 4; arity++ {
				if got, want := hasFixedArity(functions, name, arity), slices.Contains(arities, arity); got != want {
					t.Fatalf("%s fixed arity %d registered = %t, want %t; expected arities %v", name, arity, got, want, arities)
				}
			}
		})
	}
}

func TestRepeatedArgumentFunctionsRemainVariadicOnly(t *testing.T) {
	t.Parallel()

	functions := buildFunctions(t, stdlib.Full())
	names := []string{
		"INTERSECTION",
		"MINUS",
		"OUTERSECTION",
		"UNION",
		"UNION_DISTINCT",
		"KEEP_KEYS",
		"MERGE",
		"MERGE_RECURSIVE",
		"JOIN",
		"CONCAT",
		"CONCAT_SEPARATOR",
		"FMT",
		"PRINT",
	}

	for _, name := range names {
		if !functions.Var().Has(name) {
			t.Fatalf("%s is not registered as variadic", name)
		}

		for arity := 0; arity <= 4; arity++ {
			if hasFixedArity(functions, name, arity) {
				t.Fatalf("%s unexpectedly has fixed arity %d", name, arity)
			}
		}
	}
}

func boundedFunctionArities() map[string][]int {
	functions := map[string][]int{
		"APPEND":        {2, 3},
		"FLATTEN":       {1, 2},
		"POSITION":      {2, 3},
		"PUSH":          {2, 3},
		"REMOVE_VALUE":  {2, 3},
		"SLICE":         {2, 3},
		"UNSHIFT":       {2, 3},
		"DATE":          {1, 2},
		"DATE_COMPARE":  {3, 4},
		"DATE_DIFF":     {3, 4},
		"IO::FS::WRITE": {2, 3},
		"PERCENTILE":    {2, 3},
		"RAND":          {0, 1, 2},
		"RANGE":         {2, 3},
		"KEYS":          {1, 2},
		"TO_DATETIME":   {1, 2},
		"CONTAINS":      {2, 3},
		"FIND_FIRST":    {2, 3, 4},
		"FIND_LAST":     {2, 3, 4},
		"LIKE":          {2, 3},
		"LTRIM":         {1, 2},
		"REGEX_MATCH":   {2, 3},
		"REGEX_SPLIT":   {2, 3, 4},
		"REGEX_TEST":    {2, 3},
		"REGEX_REPLACE": {3, 4},
		"RTRIM":         {1, 2},
		"SPLIT":         {2, 3},
		"SUBSTITUTE":    {2, 3, 4},
		"SUBSTRING":     {2, 3},
		"TRIM":          {1, 2},
		"T::FAIL":       {0, 1},
	}

	unaryAssertions := []string{
		"EMPTY",
		"FALSE",
		"NONE",
		"TRUE",
		"BOOL",
		"STRING",
		"INT",
		"FLOAT",
		"NUMBER",
		"DURATION",
		"DATETIME",
		"ARRAY",
		"OBJECT",
		"BINARY",
	}
	binaryAssertions := []string{
		"EQ",
		"GT",
		"GTE",
		"CONTAINS",
		"HAS",
		"LEN",
		"MATCH",
		"LT",
		"LTE",
	}
	ternaryAssertions := []string{
		"APPROX",
		"BETWEEN",
	}

	for _, name := range unaryAssertions {
		functions["T::"+name] = []int{1, 2}
		functions["T::NOT::"+name] = []int{1, 2}
	}

	for _, name := range binaryAssertions {
		functions["T::"+name] = []int{2, 3}
		functions["T::NOT::"+name] = []int{2, 3}
	}

	for _, name := range ternaryAssertions {
		functions["T::"+name] = []int{3, 4}
		functions["T::NOT::"+name] = []int{3, 4}
	}

	return functions
}

func hasFixedArity(functions *runtime.Functions, name string, arity int) bool {
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
