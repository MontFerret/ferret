package stdlib_test

import (
	"reflect"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/stdlib"
)

func TestObjectSurface(t *testing.T) {
	want := []string{"object::entries", "object::from_entries", "object::has_key", "object::keep_keys", "object::keys", "object::merge", "object::merge_deep", "object::mut::keep_keys", "object::mut::merge", "object::mut::merge_deep", "object::mut::omit_keys", "object::omit_keys", "object::values", "object::zip"}
	got := functionNames(buildFunctions(t, stdlib.Only(stdlib.Objects)))
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("object surface = %v, want %v", got, want)
	}

	for _, set := range []stdlib.Set{stdlib.Full(), stdlib.Safe()} {
		functions := buildFunctions(t, set)
		without := buildFunctions(t, set.Without(stdlib.Objects))
		for _, name := range []string{"object::keys", "object::values", "object::entries", "object::from_entries"} {
			for arity := 0; arity <= 4; arity++ {
				if hasFixedArity(functions, name, arity) != (arity == 1) || functions.Var().Has(name) {
					t.Fatalf("%s must have exactly one unary registration", name)
				}
			}
		}

		for _, name := range []string{"object::has_key", "object::zip"} {
			for arity := 0; arity <= 4; arity++ {
				if hasFixedArity(functions, name, arity) != (arity == 2) || functions.Var().Has(name) {
					t.Fatalf("%s must have exactly one binary registration", name)
				}
			}
		}

		for _, name := range []string{"object::mut::merge", "object::mut::merge_deep", "object::mut::keep_keys", "object::mut::omit_keys"} {
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
			if !functions.Has(name) || without.Has(name) {
				t.Fatalf("capability mismatch: %s", name)
			}
		}

		for _, name := range []string{"keys", "values", "has", "keep_keys", "merge", "merge_recursive", "zip", "object::has", "object::merge_recursive", "object::from_keys_values", "object::mut_merge", "object::merge_mut", "object::merge_in_place", "object::mut::keys", "object::mut::set", "object::mut::delete"} {
			if functions.Has(name) {
				t.Fatalf("unexpected alias: %s", name)
			}
		}
	}
}
