package runtime_test

import (
	"math"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestMixedNumericComparisonIsExact(t *testing.T) {
	tests := []struct {
		name     string
		left     runtime.Int
		right    runtime.Float
		expected runtime.Ordering
	}{
		{name: "equal", left: 1, right: 1, expected: runtime.Equal},
		{name: "largest consecutive integer", left: 1 << 53, right: 1 << 53, expected: runtime.Equal},
		{name: "integer beyond float precision", left: 1<<53 + 1, right: 1 << 53, expected: runtime.Greater},
		{name: "maximum integer", left: math.MaxInt64, right: 1 << 53, expected: runtime.Greater},
		{name: "minimum integer", left: math.MinInt64, right: -1 << 63, expected: runtime.Equal},
		{name: "float above integer range", left: math.MaxInt64, right: 1 << 63, expected: runtime.Less},
		{name: "positive fraction", left: 1, right: 1.5, expected: runtime.Less},
		{name: "negative fraction", left: -1, right: -1.5, expected: runtime.Greater},
		{name: "negative infinity", left: 0, right: runtime.Float(math.Inf(-1)), expected: runtime.Greater},
		{name: "positive infinity", left: 0, right: runtime.Float(math.Inf(1)), expected: runtime.Less},
		{name: "signed zero", left: 0, right: runtime.Float(math.Copysign(0, -1)), expected: runtime.Equal},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertNumericComparison(t, test.left, test.right, test.expected)
			assertNumericComparison(t, test.right, test.left, reverseOrdering(test.expected))
		})
	}
}

func TestNaNIsReflexiveAndSortsLast(t *testing.T) {
	first := runtime.NewFloat(math.Float64frombits(0x7ff8000000000001))
	second := runtime.NewFloat(math.Float64frombits(0xfff8000000000002))

	assertNumericComparison(t, first, first, runtime.Equal)
	assertNumericComparison(t, first, second, runtime.Equal)
	assertNumericComparison(t, runtime.NewFloat(math.Inf(1)), first, runtime.Less)
	assertNumericComparison(t, first, runtime.NewFloat(math.Inf(1)), runtime.Greater)
	assertNumericComparison(t, runtime.NewInt(1), first, runtime.Less)
	assertNumericComparison(t, first, runtime.NewInt(1), runtime.Greater)

	if first.Hash() != second.Hash() {
		t.Fatal("equal NaN values must have equal hashes")
	}
}

func TestEqualNumericValuesHaveEqualHashes(t *testing.T) {
	negativeZero := runtime.NewFloat(math.Copysign(0, -1))

	tests := []struct {
		left  runtime.Value
		right runtime.Value
		name  string
	}{
		{name: "integer and float", left: runtime.NewInt(1), right: runtime.NewFloat(1)},
		{name: "large exact integer and float", left: runtime.NewInt64(1 << 53), right: runtime.NewFloat(1 << 53)},
		{name: "signed zero", left: runtime.ZeroInt, right: negativeZero},
		{
			name:  "nested array",
			left:  runtime.NewArrayWith(runtime.NewInt(1), runtime.ZeroInt),
			right: runtime.NewArrayWith(runtime.NewFloat(1), negativeZero),
		},
		{
			name: "nested object",
			left: runtime.NewObjectWith(map[string]runtime.Value{
				"value": runtime.NewInt(1),
			}),
			right: runtime.NewObjectWith(map[string]runtime.Value{
				"value": runtime.NewFloat(1),
			}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			equal, err := runtime.EqualValues(t.Context(), test.left, test.right)
			if err != nil || !equal {
				t.Fatalf("EqualValues() = %v, %v, want true", equal, err)
			}

			if test.left.Hash() != test.right.Hash() {
				t.Fatalf("equal values have different hashes: %d != %d", test.left.Hash(), test.right.Hash())
			}
		})
	}
}

func assertNumericComparison(t *testing.T, left, right runtime.Value, expected runtime.Ordering) {
	t.Helper()

	comparison, err := runtime.CompareValues(t.Context(), left, right)
	if err != nil {
		t.Fatalf("CompareValues(%v, %v): %v", left, right, err)
	}
	if comparison != expected {
		t.Fatalf("CompareValues(%v, %v) = %v, want %v", left, right, comparison, expected)
	}

	equal, err := runtime.EqualValues(t.Context(), left, right)
	if err != nil {
		t.Fatalf("EqualValues(%v, %v): %v", left, right, err)
	}
	if bool(equal) != (expected == runtime.Equal) {
		t.Fatalf("EqualValues(%v, %v) = %v, want %v", left, right, equal, expected == runtime.Equal)
	}
}

func reverseOrdering(ordering runtime.Ordering) runtime.Ordering {
	switch ordering {
	case runtime.Less:
		return runtime.Greater
	case runtime.Greater:
		return runtime.Less
	default:
		return runtime.Equal
	}
}
