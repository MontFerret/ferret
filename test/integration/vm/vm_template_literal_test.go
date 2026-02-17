package vm_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestTemplateLiteral(t *testing.T) {
	RunUseCases(t, []UseCase{
		Options(Case("RETURN `${@foo}`", "bar", "Should interpolate single param"), vm.WithParam("foo", runtime.NewString("bar"))),
		Options(Case("RETURN `${@a}${@b}`", "xy", "Should interpolate adjacent params"), vm.WithParam("a", runtime.NewString("x")), vm.WithParam("b", runtime.NewString("y"))),
		Options(Case("RETURN `pre-${@foo}`", "pre-bar", "Should interpolate param with prefix literal"), vm.WithParam("foo", runtime.NewString("bar"))),
		Case("RETURN `cost=\\${1}`", "cost=${1}", "Should escape interpolation marker and constant fold"),
		Options(Case("RETURN `foo-${FN()}`", "foo-bar", "Should interpolate function call"), vm.WithFunction("FN", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.NewString("bar"), nil
		})),
	})
}
