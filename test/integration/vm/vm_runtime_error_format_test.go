package vm_test

import "testing"

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
				Message: "Cannot read property \"bar\" of none",
				Contains: []string{
					"TypeError: Cannot read property \"bar\" of none",
					"property access on none",
					"Hint: Use optional chaining (?.) or check for none before accessing a member",
				},
				NotContains: []string{"Caused by:"},
			},
			"obj.fql",
		),
	})
}
