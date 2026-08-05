package operator_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/internal/operator"
)

func TestBinaryOperatorContract(t *testing.T) {
	if operator.Unknown != 0 {
		t.Fatalf("Unknown = %d, want zero", operator.Unknown)
	}

	tests := []struct {
		symbol     string
		op         operator.Binary
		relational bool
	}{
		{op: operator.Add, symbol: "+"},
		{op: operator.Subtract, symbol: "-"},
		{op: operator.Multiply, symbol: "*"},
		{op: operator.Divide, symbol: "/"},
		{op: operator.Modulus, symbol: "%"},
		{op: operator.Less, symbol: "<", relational: true},
		{op: operator.LessOrEqual, symbol: "<=", relational: true},
		{op: operator.Greater, symbol: ">", relational: true},
		{op: operator.GreaterOrEqual, symbol: ">=", relational: true},
	}

	for _, test := range tests {
		t.Run(test.symbol, func(t *testing.T) {
			if actual := test.op.String(); actual != test.symbol {
				t.Fatalf("String() = %q, want %q", actual, test.symbol)
			}

			actual, ok := operator.ParseBinary(test.symbol)
			if !ok || actual != test.op {
				t.Fatalf("ParseBinary(%q) = %v, %v; want %v, true", test.symbol, actual, ok, test.op)
			}

			if actual := test.op.IsRelational(); actual != test.relational {
				t.Fatalf("IsRelational() = %v, want %v", actual, test.relational)
			}
		})
	}
}

func TestParseBinaryRejectsNonDiagnosticOperators(t *testing.T) {
	for _, input := range []string{"", " ", "==", "!=", "IN", "+=", "?"} {
		t.Run(input, func(t *testing.T) {
			actual, ok := operator.ParseBinary(input)
			if ok || actual != operator.Unknown {
				t.Fatalf("ParseBinary(%q) = %v, %v; want Unknown, false", input, actual, ok)
			}
		})
	}

	if actual := operator.Unknown.String(); actual != "?" {
		t.Fatalf("Unknown.String() = %q, want %q", actual, "?")
	}
	if operator.Unknown.IsRelational() {
		t.Fatal("Unknown must not be relational")
	}
}

func TestCannotApplyPreservesOperandOrder(t *testing.T) {
	const expected = "operator '>' cannot be applied to String and Duration"
	if actual := operator.CannotApply(operator.Greater, "String", "Duration"); actual != expected {
		t.Fatalf("CannotApply() = %q, want %q", actual, expected)
	}
}
