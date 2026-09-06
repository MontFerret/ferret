package stdlib_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/stdlib"
)

func TestPathNamespaceSurface(t *testing.T) {
	want := []string{"path::base", "path::clean", "path::dir", "path::ext", "path::is_abs", "path::join", "path::match", "path::separate"}
	for _, test := range []struct {
		name string
		set  stdlib.Set
	}{
		{name: "only", set: stdlib.Only(stdlib.Path)},
		{name: "full", set: stdlib.Full()},
		{name: "safe", set: stdlib.Safe()},
	} {
		t.Run(test.name, func(t *testing.T) {
			functions := buildFunctions(t, test.set)
			without := buildFunctions(t, test.set.Without(stdlib.Path))
			var got []string
			for _, name := range functionNames(functions) {
				if strings.HasPrefix(name, "path::") {
					got = append(got, name)
				}
			}

			if !reflect.DeepEqual(got, want) {
				t.Fatalf("path functions = %v, want %v", got, want)
			}

			if test.name == "only" && !reflect.DeepEqual(functionNames(functions), want) {
				t.Fatalf("PATH-only functions = %v, want %v", functionNames(functions), want)
			}

			for _, name := range want {
				if without.Has(name) {
					t.Fatalf("%s remains after excluding PATH", name)
				}

				arity := 1
				switch name {
				case "path::match":
					arity = 2
				case "path::join":
					arity = -1
				}

				if functions.Var().Has(name) != (arity == -1) {
					t.Fatalf("%s has incorrect variadic registration", name)
				}

				for count := 0; count <= 4; count++ {
					if hasFixedArity(functions, name, count) != (count == arity) {
						t.Fatalf("%s has incorrect fixed arity %d", name, count)
					}
				}

				global := strings.TrimPrefix(name, "path::")
				if global != "join" && (functions.Has(global) || without.Has(global)) {
					t.Fatalf("obsolete global PATH function %s remains", global)
				}
			}

			if functions.Var().Has("join") || without.Var().Has("join") {
				t.Fatal("global path join remains registered")
			}

			for count := 0; count <= 4; count++ {
				wantStringJoin := test.name != "only" && count == 2
				if hasFixedArity(functions, "join", count) != wantStringJoin || hasFixedArity(without, "join", count) != wantStringJoin {
					t.Fatalf("string join has incorrect fixed arity %d", count)
				}
			}
		})
	}
}
