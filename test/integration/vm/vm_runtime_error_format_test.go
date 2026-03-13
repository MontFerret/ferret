package vm_test

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestRuntimeErrorFormatting(t *testing.T) {
	RunUseCases(t, []UseCase{
		RuntimeErrorCase(
			"LET numerator = 10\nRETURN numerator / 0",
			ExpectedRuntimeError{
				Message: "Division by zero",
				Contains: []string{
					"DivideByZero: Division by zero",
					":2:8",
					"attempt to divide by zero",
					"Hint: Ensure the denominator is non-zero before division",
					"Note: Add a conditional check before dividing",
				},
				NotContains: []string{"~"},
			},
			"script.fql",
		),
		RuntimeErrorCase(
			"LET obj = {}\nRETURN obj.foo.bar",
			ExpectedRuntimeError{
				Message: "Cannot read property \"bar\" of None",
				Contains: []string{
					"TypeError: Cannot read property \"bar\" of None",
					"property access on None",
					"Hint: Use optional chaining (?.) or check for None before accessing a member",
				},
				NotContains: []string{"Caused by:"},
			},
			"obj.fql",
		),
		Options(
			RuntimeErrorCase(
				`
FUNC Inner() => FAIL()
FUNC Outer() (
  LET result = Inner()
  RETURN result
)
RETURN Outer()
`,
				ExpectedRuntimeError{
					Contains: []string{
						"called from Inner (#1)",
						"called from Outer (#2)",
						"VM stack: Outer -> Inner",
					},
				},
				"nested_udf_stack.fql",
			),
			vm.WithFunction("FAIL", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
				return runtime.None, errors.New("boom")
			}),
		),
	})
}
