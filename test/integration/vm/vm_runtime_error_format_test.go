package vm_test

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestRuntimeErrorFormatting(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		spec.NewSpec(
			"LET numerator = 10\nRETURN numerator / 0",
			"script.fql",
		).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message: "Division by zero",
			Contains: []string{
				"DivideByZero: Division by zero",
				":2:8",
				"attempt to divide by zero",
				"Hint: Ensure the denominator is non-zero before division",
				"Note: Add a conditional check before dividing",
			},
			NotContains: []string{"~"},
		}),
		spec.NewSpec(
			"LET obj = {}\nRETURN obj.foo.bar",
			"obj.fql",
		).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message: "Cannot read property \"bar\" of None",
			Contains: []string{
				"TypeError: Cannot read property \"bar\" of None",
				"property access on None",
				"Hint: Use optional chaining (?.) or check for None before accessing a member",
			},
			NotContains: []string{"Caused by:"},
		}),
		spec.NewSpec(
			`
FUNC Inner() => FAIL()
FUNC Outer() (
  LET result = Inner()
  RETURN result
)
RETURN Outer()
`,
			"nested_udf_stack.fql",
		).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Contains: []string{
				"called from Inner (#1)",
				"called from Outer (#2)",
				"VM stack: Outer -> Inner",
			},
		}).Env(
			vm.WithFunction("FAIL", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
				return runtime.None, errors.New("boom")
			}),
		),
	})
}
