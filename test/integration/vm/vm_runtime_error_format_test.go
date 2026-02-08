package vm_test

import "testing"

func TestRuntimeErrorFormatting(t *testing.T) {
	useCases := []UseCase{
		RuntimeErrorCase(
			"LET numerator = 10\nRETURN numerator / 0",
			ExpectedRuntimeError{
				Message: "division by zero",
				Contains: []string{
					"error: division by zero",
					":2:8",
					"attempt to divide by zero",
					"= help: ensure the denominator is non-zero before division",
					"= note: add a conditional check before dividing",
				},
				NotContains: []string{"~"},
			},
			"script.fql",
		),
	}

	RunUseCases(t, useCases)
}
