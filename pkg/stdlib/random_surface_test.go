package stdlib_test

import (
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/stdlib"
)

func TestRandomSurface(t *testing.T) {
	want := []string{"random::bool", "random::float", "random::int"}
	if got := functionNames(buildFunctions(t, stdlib.Only(stdlib.Random))); !reflect.DeepEqual(got, want) {
		t.Fatalf("Random functions = %v, want %v", got, want)
	}

	for _, selection := range []struct {
		set            stdlib.Set
		random, legacy bool
	}{
		{stdlib.Only(stdlib.Random), true, false}, {stdlib.Only(stdlib.Math), false, true},
		{stdlib.Full(), true, true}, {stdlib.Safe(), true, true},
		{stdlib.Full().Without(stdlib.Random), false, true},
	} {
		functions := buildFunctions(t, selection.set)
		for name, arities := range map[string][]int{"random::float": {0, 2}, "random::int": {2}, "random::bool": {0}, "rand": {0, 1, 2}} {
			present := selection.random
			if name == "rand" {
				present = selection.legacy
			}

			for _, spelling := range []string{name, strings.ToUpper(name)} {
				for arity := 0; arity <= 4; arity++ {
					if got, want := hasFixedArity(functions, spelling, arity), present && slices.Contains(arities, arity); got != want {
						t.Fatalf("%s/%d availability %t, want %t", spelling, arity, got, want)
					}
				}

				if functions.Var().Has(spelling) {
					t.Fatalf("unexpected variadic %s", spelling)
				}
			}
		}

		for _, name := range []string{"float", "int", "bool", "math::rand", "math::random", "random::seed", "random::choice", "random::shuffle"} {
			if functions.Has(name) {
				t.Fatalf("unexpected function %s", name)
			}
		}
	}
}
