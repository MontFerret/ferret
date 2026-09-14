package stdlib_test

import (
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/stdlib"
)

func TestMathSurface(t *testing.T) {
	canonical := map[string][]int{"e": {0}, "pi": {0}, "atan2": {2}, "hypot": {2}, "pow": {2}, "percentile": {2}, "clamp": {3}}
	for _, name := range strings.Fields("abs acos asin atan cbrt ceil cos degrees exp exp2 expm1 floor log log1p log2 log10 max mean median min radians round sign sin sqrt stddev stddev_sample sum tan trunc variance variance_sample") {
		canonical[name] = []int{1}
	}

	legacy := map[string][]int{"pi": {0}, "atan2": {2}, "pow": {2}, "percentile": {2, 3}, "rand": {0, 1, 2}, "range": {2, 3}}
	for _, name := range strings.Fields("abs acos asin atan average ceil cos degrees exp exp2 floor log log2 log10 max median min radians round sin sqrt stddev_population stddev_sample sum tan variance_population variance_sample") {
		legacy[name] = []int{1}
	}

	want := []string{}
	for name := range canonical {
		want = append(want, "math::"+name)
	}

	for name := range legacy {
		want = append(want, name)
	}

	slices.Sort(want)
	if got := functionNames(buildFunctions(t, stdlib.Only(stdlib.Math))); !reflect.DeepEqual(got, want) {
		t.Fatalf("Math functions = %v, want %v", got, want)
	}

	for _, selection := range []struct {
		set    stdlib.Set
		math   bool
		arrays bool
	}{
		{stdlib.Only(stdlib.Math), true, false},
		{stdlib.Only(stdlib.Arrays), false, true},
		{stdlib.Full(), true, true},
		{stdlib.Safe(), true, true},
	} {
		functions := buildFunctions(t, selection.set)
		for _, surface := range []struct {
			entries map[string][]int
			prefix  string
		}{{prefix: "math::", entries: canonical}, {entries: legacy}, {prefix: "arrays::", entries: map[string][]int{"range": {2, 3}}}} {
			for name, arities := range surface.entries {
				name = surface.prefix + name
				present := selection.math
				if surface.prefix == "arrays::" {
					present = selection.arrays
				}

				for _, spelling := range []string{name, strings.ToUpper(name)} {
					for arity := 0; arity <= 4; arity++ {
						if got, want := hasFixedArity(functions, spelling, arity), present && slices.Contains(arities, arity); got != want {
							t.Fatalf("%s/%d availability=%t, want %t", spelling, arity, got, want)
						}
					}

					if functions.Var().Has(spelling) {
						t.Fatalf("unexpected variadic %s", spelling)
					}
				}
			}
		}

		for _, name := range strings.Fields("math::average math::variance_population math::stddev_population math::range math::rand clamp sign trunc cbrt hypot log1p expm1 e random::float") {
			for _, spelling := range []string{name, strings.ToUpper(name)} {
				if functions.Has(spelling) {
					t.Fatalf("unexpected function %s", spelling)
				}
			}
		}
	}
}
