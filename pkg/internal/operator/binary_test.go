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
		value      uint8
		relational bool
		equality   bool
	}{
		{op: operator.Add, symbol: "+", value: 1},
		{op: operator.Subtract, symbol: "-", value: 2},
		{op: operator.Multiply, symbol: "*", value: 3},
		{op: operator.Divide, symbol: "/", value: 4},
		{op: operator.Modulus, symbol: "%", value: 5},
		{op: operator.Less, symbol: "<", value: 6, relational: true},
		{op: operator.LessOrEqual, symbol: "<=", value: 7, relational: true},
		{op: operator.Greater, symbol: ">", value: 8, relational: true},
		{op: operator.GreaterOrEqual, symbol: ">=", value: 9, relational: true},
		{op: operator.Equal, symbol: "==", value: 10, equality: true},
		{op: operator.NotEqual, symbol: "!=", value: 11, equality: true},
	}

	seen := map[operator.Binary]string{operator.Unknown: "?"}
	for _, test := range tests {
		t.Run(test.symbol, func(t *testing.T) {
			if actual := uint8(test.op); actual != test.value {
				t.Fatalf("numeric value = %d, want %d", actual, test.value)
			}
			if previous, exists := seen[test.op]; exists {
				t.Fatalf("numeric value %d is shared by %q and %q", test.op, previous, test.symbol)
			}
			seen[test.op] = test.symbol

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

			if actual := test.op.IsEquality(); actual != test.equality {
				t.Fatalf("IsEquality() = %v, want %v", actual, test.equality)
			}
		})
	}
}

func TestParseBinaryRejectsNonDiagnosticOperators(t *testing.T) {
	for _, input := range []string{"", " ", "=", "IN", "+=", "?"} {
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
	if operator.Unknown.IsEquality() {
		t.Fatal("Unknown must not be equality")
	}
}

func TestCannotApplyPreservesOperandOrder(t *testing.T) {
	const expected = "operator '>' cannot be applied to String and Duration"
	if actual := operator.CannotApply(operator.Greater, "String", "Duration"); actual != expected {
		t.Fatalf("CannotApply() = %q, want %q", actual, expected)
	}
}
