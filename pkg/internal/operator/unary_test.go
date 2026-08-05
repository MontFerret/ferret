package operator_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/internal/operator"
)

func TestUnaryOperatorContract(t *testing.T) {
	if operator.UnknownUnary != 0 {
		t.Fatalf("UnknownUnary = %d, want zero", operator.UnknownUnary)
	}

	tests := []struct {
		symbol string
		op     operator.Unary
		value  uint8
	}{
		{op: operator.Not, symbol: "!", value: 1},
		{op: operator.Positive, symbol: "+", value: 2},
		{op: operator.Negative, symbol: "-", value: 3},
		{op: operator.Increment, symbol: "++", value: 4},
		{op: operator.Decrement, symbol: "--", value: 5},
	}

	seen := map[operator.Unary]string{operator.UnknownUnary: "?"}
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

			actual, ok := operator.ParseUnary(test.symbol)
			if !ok || actual != test.op {
				t.Fatalf("ParseUnary(%q) = %v, %v; want %v, true", test.symbol, actual, ok, test.op)
			}
		})
	}
}

func TestParseUnaryRejectsUnsupportedOperators(t *testing.T) {
	for _, input := range []string{"", " ", "NOT", "~", "=", "+="} {
		t.Run(input, func(t *testing.T) {
			actual, ok := operator.ParseUnary(input)
			if ok || actual != operator.UnknownUnary {
				t.Fatalf("ParseUnary(%q) = %v, %v; want UnknownUnary, false", input, actual, ok)
			}
		})

	}

	if actual := operator.UnknownUnary.String(); actual != "?" {
		t.Fatalf("UnknownUnary.String() = %q, want %q", actual, "?")
	}
}

func TestCannotApplyUnary(t *testing.T) {
	const expected = "operator '!' cannot be applied to String"
	if actual := operator.CannotApplyUnary(operator.Not, "String"); actual != expected {
		t.Fatalf("CannotApplyUnary() = %q, want %q", actual, expected)
	}
}
