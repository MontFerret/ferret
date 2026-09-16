package stdlib_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/stdlib"
)

func TestCollectionsSurface(t *testing.T) {
	want := []string{
		"collections::count", "collections::count_distinct", "collections::includes", "collections::reverse",
		"count", "count_distinct", "includes", "reverse",
	}
	got := functionNames(buildFunctions(t, stdlib.Only(stdlib.Collections)))
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("collections surface = %v, want %v", got, want)
	}

	for _, set := range []stdlib.Set{stdlib.Only(stdlib.Collections), stdlib.Full(), stdlib.Safe()} {
		functions := buildFunctions(t, set)
		without := buildFunctions(t, set.Without(stdlib.Collections))
		for _, name := range want {
			expectedArity := 1
			if name == "includes" || name == "collections::includes" {
				expectedArity = 2
			}

			for _, spelling := range []string{name, strings.ToUpper(name)} {
				for arity := 0; arity <= 4; arity++ {
					if hasFixedArity(functions, spelling, arity) != (arity == expectedArity) || functions.Var().Has(spelling) {
						t.Fatalf("wrong arity %d for %s", arity, spelling)
					}
				}

				if without.Has(spelling) {
					t.Fatalf("collections group removal retained %s", spelling)
				}
			}
		}

		for _, prefix := range []string{"collection::", "coll::", "col::", "arrays::"} {
			for _, name := range []string{"count", "count_distinct", "includes", "reverse"} {
				if functions.Has(prefix + name) {
					t.Fatalf("unexpected alias %s%s", prefix, name)
				}
			}
		}
	}
}
