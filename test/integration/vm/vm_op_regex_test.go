package vm_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/file"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRegexpOperator(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case(`RETURN "foo" =~ "^f[o].$" `, true),
		Case(`RETURN "foo" !~ "[a-z]+bar$"`, true),
		Case(`RETURN "foo" !~ T::REGEXP()`, true),
	}, vm.WithFunction("T::REGEXP", func(_ context.Context, _ ...runtime.Value) (value runtime.Value, e error) {
		return runtime.NewString("[a-z]+bar$"), nil
	}))

	// TODO: Fix
	SkipConvey("Should return an error during compilation when a regexp string invalid", t, func() {
		_, err := compiler.New(compiler.WithOptimizationLevel(compiler.O0)).
			Compile(file.NewAnonymousSource(`
			RETURN "foo" !~ "[ ]\K(?<!\d )(?=(?: ?\d){8})(?!(?: ?\d){9})\d[ \d]+\d" 
		`))

		So(err, ShouldBeError)
	})

	// TODO: Fix
	SkipConvey("Should return an error during compilation when a regexp is not a string", t, func() {
		right := []string{
			"[]",
			"{}",
			"1",
			"1.1",
			"TRUE",
		}

		for _, r := range right {
			_, err := compiler.New(compiler.WithOptimizationLevel(compiler.O0)).
				Compile(file.NewAnonymousSource(fmt.Sprintf(`
			RETURN "foo" !~ %s 
		`, r)))

			So(err, ShouldBeError)
		}
	})
}
