package operationerror_test

import (
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/internal/operationerror"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestZeroDivisorErrorsPreserveCauseAndIdentity(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		identity error
		expected string
	}{
		{name: "division", err: operationerror.DivisionByZero(runtime.ErrInvalidOperation), identity: operationerror.ErrDivisionByZero, expected: "invalid operation: division by zero"},
		{name: "modulo", err: operationerror.ModuloByZero(runtime.ErrInvalidOperation), identity: operationerror.ErrModuloByZero, expected: "invalid operation: modulo by zero"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.err.Error() != test.expected {
				t.Fatalf("error = %q, want %q", test.err, test.expected)
			}
			if !errors.Is(test.err, runtime.ErrInvalidOperation) {
				t.Fatal("error lost ErrInvalidOperation identity")
			}
			if !errors.Is(test.err, test.identity) {
				t.Fatal("error lost zero-divisor identity")
			}
		})
	}
}
