package stdlib_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/stdlib"
)

func TestObjectSurface(t *testing.T) {
	want := []string{"has", "keep_keys", "keys", "merge", "merge_recursive", "object::entries", "object::from_entries", "object::has_key", "object::keep_keys", "object::keys", "object::merge", "object::merge_deep", "object::mut::keep_keys", "object::mut::merge", "object::mut::merge_deep", "object::mut::omit_keys", "object::omit_keys", "object::values", "object::zip", "values", "zip"}
	got := functionNames(buildFunctions(t, stdlib.Only(stdlib.Objects)))
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("object surface = %v, want %v", got, want)
	}

	for _, set := range []stdlib.Set{stdlib.Only(stdlib.Objects), stdlib.Full(), stdlib.Safe()} {
		functions := buildFunctions(t, set)
		without := buildFunctions(t, set.Without(stdlib.Objects))
		for _, name := range []string{"keys", "values", "object::keys", "object::values", "object::entries", "object::from_entries"} {
			for arity := 0; arity <= 4; arity++ {
				if hasFixedArity(functions, name, arity) != (arity == 1) || functions.Var().Has(name) {
					t.Fatalf("%s must have exactly one unary registration", name)
				}
			}
		}

		for _, name := range []string{"has", "zip", "object::has_key", "object::zip"} {
			for arity := 0; arity <= 4; arity++ {
				if hasFixedArity(functions, name, arity) != (arity == 2) || functions.Var().Has(name) {
					t.Fatalf("%s must have exactly one binary registration", name)
				}
			}
		}

		for _, name := range []string{"keep_keys", "merge", "merge_recursive", "object::keep_keys", "object::omit_keys", "object::merge", "object::merge_deep", "object::mut::merge", "object::mut::merge_deep", "object::mut::keep_keys", "object::mut::omit_keys"} {
			if !functions.Var().Has(name) {
				t.Fatalf("%s must have a variadic registration", name)
			}

			for arity := 0; arity <= 4; arity++ {
				if hasFixedArity(functions, name, arity) {
					t.Fatalf("unexpected fixed registration for %s", name)
				}
			}
		}

		for _, name := range want {
			for _, spelling := range []string{name, strings.ToUpper(name)} {
				if !functions.Has(spelling) || without.Has(spelling) {
					t.Fatalf("capability mismatch: %s", spelling)
				}
			}
		}

		for _, name := range []string{"entries", "from_entries", "omit_keys", "has_key", "merge_deep", "object::has", "object::merge_recursive", "object::from_keys_values", "mut_merge", "merge_mut", "merge_in_place", "object::mut_merge", "object::merge_mut", "object::merge_in_place", "object::mut::keys", "object::mut::set", "object::mut::delete"} {
			if functions.Has(name) {
				t.Fatalf("unexpected alias: %s", name)
			}
		}

		for _, name := range []string{"merge", "merge_deep", "keep_keys", "omit_keys"} {
			for _, alias := range []string{"mut::" + name, "mut_" + name, name + "_mut", name + "_in_place"} {
				if functions.Has(alias) {
					t.Fatalf("unexpected mutable alias: %s", alias)
				}
			}
		}
	}
}
