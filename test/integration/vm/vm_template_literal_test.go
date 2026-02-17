package vm_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestTemplateLiteral(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case("RETURN ``", "", "Should return empty string for empty template literal"),
		Options(Case("RETURN `${@foo}`", "bar", "Should interpolate single param"), vm.WithParam("foo", runtime.NewString("bar"))),
		Options(Case("RETURN `${@a}${@b}`", "xy", "Should interpolate adjacent params"), vm.WithParam("a", runtime.NewString("x")), vm.WithParam("b", runtime.NewString("y"))),
		Options(Case("RETURN `pre-${@foo}`", "pre-bar", "Should interpolate param with prefix literal"), vm.WithParam("foo", runtime.NewString("bar"))),
		Options(Case("RETURN `pre-${@foo}-post`", "pre-bar-post", "Should interpolate param with prefix and suffix literals"), vm.WithParam("foo", runtime.NewString("bar"))),
		Options(Case("RETURN `x${@foo}y`", "xy", "Should coerce NONE param to empty string"), vm.WithParam("foo", runtime.None)),
		Case("RETURN `cost=\\${1}`", "cost=${1}", "Should escape interpolation marker and constant fold"),
		Case("RETURN `line\\nend`", "line\nend", "Should handle newline escapes in template literals"),
		Case("RETURN `tab\\tend`", "tab\tend", "Should handle tab escapes in template literals"),
		Case("RETURN `use \\`backtick\\``", "use `backtick`", "Should allow escaped backticks in template literals"),
		Case("RETURN `slash\\\\test`", "slash\\test", "Should allow escaped backslash in template literals"),
		Options(Case("RETURN `foo-${FN()}`", "foo-bar", "Should interpolate function call"), vm.WithFunction("FN", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.NewString("bar"), nil
		})),
	})
}
