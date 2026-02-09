package vm_test

import "testing"

func TestRuntimeErrorFormatting(t *testing.T) {
	useCases := []UseCase{
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
	}

	RunUseCases(t, useCases)
}
