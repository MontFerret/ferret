package vm_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestTemplateLiteral(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		S("RETURN ``", "", "Should return empty string for empty template literal"),
		S("RETURN `${@foo}`", "bar", "Should interpolate single param").Env(vm.WithParam("foo", runtime.NewString("bar"))),
		S("RETURN `${@a}${@b}`", "xy", "Should interpolate adjacent params").Env(vm.WithParam("a", runtime.NewString("x")), vm.WithParam("b", runtime.NewString("y"))),
		S("RETURN `pre-${@foo}`", "pre-bar", "Should interpolate param with prefix literal").Env(vm.WithParam("foo", runtime.NewString("bar"))),
		S("RETURN `pre-${@foo}-post`", "pre-bar-post", "Should interpolate param with prefix and suffix literals").Env(vm.WithParam("foo", runtime.NewString("bar"))),
		S("RETURN `x${@foo}y`", "xy", "Should coerce NONE param to empty string").Env(vm.WithParam("foo", runtime.None)),
		S("RETURN `cost=\\${1}`", "cost=${1}", "Should escape interpolation marker and constant fold"),
		S("RETURN `line\\nend`", "line\nend", "Should handle newline escapes in template literals"),
		S("RETURN `tab\\tend`", "tab\tend", "Should handle tab escapes in template literals"),
		S("RETURN `use \\`backtick\\``", "use `backtick`", "Should allow escaped backticks in template literals"),
		S("RETURN `slash\\\\test`", "slash\\test", "Should allow escaped backslash in template literals"),
		S("RETURN `foo-${FN()}`", "foo-bar", "Should interpolate function call").Env(vm.WithFunction("FN", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.NewString("bar"), nil
		})),
	})
}
