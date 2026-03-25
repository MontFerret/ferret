package vm_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestRegexpOperator(t *testing.T) {
	RunSpecs(t, []Spec{
		S(`RETURN "foo" =~ "^f[o].$" `, true),
		S(`RETURN "foo" !~ "[a-z]+bar$"`, true),
		S(`RETURN "foo" !~ T::REGEXP()`, true),
		Array(`FOR p IN ["^f..$", "^b..$"] RETURN "foo" =~ p`, []any{true, false}),
		Array(`FOR p IN [1, 2] RETURN "foo" =~ T::REGEXP_DYNAMIC(p)`, []any{true, false}),
	}, vm.WithFunction("T::REGEXP", func(_ context.Context, _ ...runtime.Value) (value runtime.Value, e error) {
		return runtime.NewString("[a-z]+bar$"), nil
	}), vm.WithFunction("T::REGEXP_DYNAMIC", func(_ context.Context, args ...runtime.Value) (value runtime.Value, e error) {
		if len(args) > 0 && args[0].String() == "1" {
			return runtime.NewString("^f..$"), nil
		}

		return runtime.NewString("^b..$"), nil
	}))

	t.Run("Should return an error during compilation when a regexp string invalid", func(t *testing.T) {
		_, err := compiler.New(compiler.WithOptimizationLevel(compiler.O0)).
			Compile(file.NewAnonymousSource(`
			RETURN "foo" !~ "[ ]\K(?<!\d )(?=(?: ?\d){8})(?!(?: ?\d){9})\d[ \d]+\d" 
		`))

		if err == nil {
			t.Fatal("expected compilation error")
		}
	})

	t.Run("Should return an error during compilation when a regexp is not a string", func(t *testing.T) {
		right := []string{
			"[]",
			"{}",
			"1",
			"1.1",
			"TRUE",
		}

		for _, r := range right {
			r := r

			t.Run(r, func(t *testing.T) {
				_, err := compiler.New(compiler.WithOptimizationLevel(compiler.O0)).
					Compile(file.NewAnonymousSource(fmt.Sprintf(`
			RETURN "foo" !~ %s 
		`, r)))

				if err == nil {
					t.Fatal("expected compilation error")
				}
			})
		}
	})
}
