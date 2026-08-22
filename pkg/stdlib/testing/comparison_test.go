package testing

import (
	"context"
	"errors"
	"math"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestEqualAssertion(t *testing.T) {
	t.Parallel()

	one := runtime.NewArrayWith(runtime.NewInt(1))
	oneTwo := runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2))

	tests := []struct {
		name     string
		actual   runtime.Value
		expected runtime.Value
		failure  string
		negative string
	}{
		{
			name:     "string",
			actual:   runtime.NewString("Foo"),
			expected: runtime.NewString("Bar"),
			failure:  "assertion error: values are not equal\nexpected: String 'Bar'\nactual:   String 'Foo'",
			negative: "assertion error: expected values to differ\nboth: String 'Bar'",
		},
		{
			name:     "escaped string",
			actual:   runtime.NewString(`can't\stop`),
			expected: runtime.NewString("won't"),
			failure:  "assertion error: values are not equal\nexpected: String 'won\\'t'\nactual:   String 'can\\'t\\\\stop'",
			negative: "assertion error: expected values to differ\nboth: String 'won\\'t'",
		},
		{
			name:     "int",
			actual:   runtime.NewInt(1),
			expected: runtime.NewInt(2),
			failure:  "assertion error: values are not equal\nexpected: Int '2'\nactual:   Int '1'",
			negative: "assertion error: expected values to differ\nboth: Int '2'",
		},
		{
			name:     "boolean",
			actual:   runtime.False,
			expected: runtime.True,
			failure:  "assertion error: values are not equal\nexpected: Boolean 'true'\nactual:   Boolean 'false'",
			negative: "assertion error: expected values to differ\nboth: Boolean 'true'",
		},
		{
			name:     "array",
			actual:   one,
			expected: oneTwo,
			failure:  "assertion error: values are not equal\n$[1]\n  expected: Int '2'\n  actual:   <missing>",
			negative: "assertion error: expected values to differ\nboth: Array '[1,2]'",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			requireAssertionFailure(t, equalAssertion, true, test.failure, test.actual, test.expected)
			requireAssertionSuccess(t, equalAssertion, false, test.actual, test.expected)
			requireAssertionSuccess(t, equalAssertion, true, test.expected, test.expected)
			requireAssertionFailure(t, equalAssertion, false, test.negative, test.expected, test.expected)
		})
	}
}

func TestOrderingAssertions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		failure    string
		descriptor assertion
		actual     runtime.Int
		expected   runtime.Int
		positive   bool
	}{
		{name: "gt less", descriptor: gtAssertion, actual: 1, expected: 2, positive: true, failure: "assertion error: expected Int '1' to be greater than Int '2'"},
		{name: "gt equal", descriptor: gtAssertion, actual: 1, expected: 1, positive: true, failure: "assertion error: expected Int '1' to be greater than Int '1'"},
		{name: "not gt greater", descriptor: gtAssertion, actual: 2, expected: 1, positive: false, failure: "assertion error: expected Int '2' not to be greater than Int '1'"},
		{name: "gte less", descriptor: gteAssertion, actual: 1, expected: 2, positive: true, failure: "assertion error: expected Int '1' to be greater than or equal to Int '2'"},
		{name: "not gte equal", descriptor: gteAssertion, actual: 1, expected: 1, positive: false, failure: "assertion error: expected Int '1' not to be greater than or equal to Int '1'"},
		{name: "not gte greater", descriptor: gteAssertion, actual: 2, expected: 1, positive: false, failure: "assertion error: expected Int '2' not to be greater than or equal to Int '1'"},
		{name: "lt greater", descriptor: ltAssertion, actual: 2, expected: 1, positive: true, failure: "assertion error: expected Int '2' to be less than Int '1'"},
		{name: "lt equal", descriptor: ltAssertion, actual: 1, expected: 1, positive: true, failure: "assertion error: expected Int '1' to be less than Int '1'"},
		{name: "not lt less", descriptor: ltAssertion, actual: 1, expected: 2, positive: false, failure: "assertion error: expected Int '1' not to be less than Int '2'"},
		{name: "lte greater", descriptor: lteAssertion, actual: 2, expected: 1, positive: true, failure: "assertion error: expected Int '2' to be less than or equal to Int '1'"},
		{name: "not lte less", descriptor: lteAssertion, actual: 1, expected: 2, positive: false, failure: "assertion error: expected Int '1' not to be less than or equal to Int '2'"},
		{name: "not lte equal", descriptor: lteAssertion, actual: 1, expected: 1, positive: false, failure: "assertion error: expected Int '1' not to be less than or equal to Int '1'"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			requireAssertionFailure(t, test.descriptor, test.positive, test.failure, test.actual, test.expected)
			requireAssertionSuccess(t, test.descriptor, !test.positive, test.actual, test.expected)
		})
	}
}

func TestApproxAssertion(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		actual    runtime.Value
		expected  runtime.Value
		tolerance runtime.Value
		name      string
	}{
		{name: "exact integers", actual: runtime.NewInt(10), expected: runtime.NewInt(10), tolerance: runtime.ZeroInt},
		{name: "within float tolerance", actual: runtime.NewFloat(10), expected: runtime.NewFloat(10.001), tolerance: runtime.NewFloat(0.01)},
		{name: "inclusive mixed boundary", actual: runtime.NewInt(10), expected: runtime.NewFloat(10.5), tolerance: runtime.NewFloat(0.5)},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			requireAssertionSuccess(t, approxAssertion, true, test.actual, test.expected, test.tolerance)
		})
	}

	requireAssertionFailure(
		t,
		approxAssertion,
		true,
		"assertion error: expected Float '10' to approximately equal Float '11' within Float '0.01'",
		runtime.NewFloat(10),
		runtime.NewFloat(11),
		runtime.NewFloat(0.01),
	)
	requireAssertionSuccess(
		t,
		approxAssertion,
		false,
		runtime.NewFloat(10),
		runtime.NewFloat(11),
		runtime.NewFloat(0.01),
	)
	requireAssertionFailure(
		t,
		approxAssertion,
		false,
		"assertion error: expected Int '10' not to approximately equal Int '10' within Int '0'",
		runtime.NewInt(10),
		runtime.NewInt(10),
		runtime.ZeroInt,
	)
	requireAssertionFailure(
		t,
		approxAssertion,
		true,
		"assertion error: expected Int '10' to approximately equal Int '11' within Int '0'",
		runtime.NewInt(10),
		runtime.NewInt(11),
		runtime.ZeroInt,
	)
	requireAssertionFailure(
		t,
		approxAssertion,
		true,
		"assertion error: values are not close enough",
		runtime.NewInt(10),
		runtime.NewInt(12),
		runtime.NewInt(1),
		runtime.NewString("values are not close enough"),
	)
}

func TestApproxAssertionRejectsInvalidUsage(t *testing.T) {
	t.Parallel()

	for _, positive := range []bool{true, false} {
		fn := approxAssertion.negative()
		if positive {
			fn = approxAssertion.positive()
		}

		out, err := fn(context.Background(), runtime.NewInt(10), runtime.NewInt(10), runtime.NewInt(-1))
		if out != runtime.None {
			t.Fatalf("negative tolerance output = %v, want None", out)
		}
		if !errors.Is(err, runtime.ErrInvalidArgument) || errors.Is(err, errAssertion) {
			t.Fatalf("negative tolerance error = %v, want invalid argument without assertion failure", err)
		}
	}

	for pos := 0; pos < 3; pos++ {
		args := []runtime.Value{runtime.NewInt(10), runtime.NewInt(10), runtime.NewInt(1)}
		args[pos] = runtime.NewString("10")

		out, err := approxAssertion.positive()(context.Background(), args...)
		if out != runtime.None {
			t.Fatalf("invalid argument %d output = %v, want None", pos, out)
		}
		if !errors.Is(err, runtime.ErrInvalidArgument) || !errors.Is(err, runtime.ErrInvalidType) {
			t.Fatalf("invalid argument %d error = %v, want invalid argument type", pos, err)
		}
	}
}

func TestApproxAssertionPropagatesNumericErrors(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		actual    runtime.Value
		expected  runtime.Value
		tolerance runtime.Value
		name      string
	}{
		{name: "NaN actual", actual: runtime.NewFloat(math.NaN()), expected: runtime.ZeroFloat, tolerance: runtime.NewFloat(1)},
		{name: "infinite expected", actual: runtime.ZeroFloat, expected: runtime.NewFloat(math.Inf(1)), tolerance: runtime.NewFloat(1)},
		{name: "NaN tolerance", actual: runtime.ZeroFloat, expected: runtime.ZeroFloat, tolerance: runtime.NewFloat(math.NaN())},
		{name: "infinite tolerance", actual: runtime.ZeroFloat, expected: runtime.ZeroFloat, tolerance: runtime.NewFloat(math.Inf(1))},
		{name: "integer subtraction overflow", actual: runtime.NewInt64(math.MaxInt64), expected: runtime.NewInt(-1), tolerance: runtime.NewInt(1)},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			out, err := approxAssertion.positive()(context.Background(), test.actual, test.expected, test.tolerance)
			if out != runtime.None {
				t.Fatalf("output = %v, want None", out)
			}
			if !errors.Is(err, runtime.ErrRange) || errors.Is(err, errAssertion) {
				t.Fatalf("error = %v, want propagated range error", err)
			}
		})
	}
}

func TestBetweenAssertion(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		value runtime.Value
		min   runtime.Value
		max   runtime.Value
		name  string
	}{
		{name: "strictly inside", value: runtime.NewInt(15), min: runtime.NewInt(10), max: runtime.NewInt(20)},
		{name: "equal minimum", value: runtime.NewInt(10), min: runtime.NewInt(10), max: runtime.NewInt(20)},
		{name: "equal maximum", value: runtime.NewInt(20), min: runtime.NewInt(10), max: runtime.NewInt(20)},
		{name: "mixed numeric representations", value: runtime.NewFloat(15.5), min: runtime.NewInt(10), max: runtime.NewFloat(20)},
		{name: "canonical string ordering", value: runtime.NewString("m"), min: runtime.NewString("a"), max: runtime.NewString("z")},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			requireAssertionSuccess(t, betweenAssertion, true, test.value, test.min, test.max)
		})
	}

	for _, value := range []runtime.Value{runtime.NewInt(9), runtime.NewInt(21)} {
		requireAssertionSuccess(t, betweenAssertion, false, value, runtime.NewInt(10), runtime.NewInt(20))
	}

	requireAssertionFailure(
		t,
		betweenAssertion,
		true,
		"assertion error: expected Int '9' to be between Int '10' and Int '20'",
		runtime.NewInt(9),
		runtime.NewInt(10),
		runtime.NewInt(20),
	)
	requireAssertionFailure(
		t,
		betweenAssertion,
		false,
		"assertion error: expected Int '15' not to be between Int '10' and Int '20'",
		runtime.NewInt(15),
		runtime.NewInt(10),
		runtime.NewInt(20),
	)
	requireAssertionFailure(
		t,
		betweenAssertion,
		true,
		"assertion error: status outside range",
		runtime.NewInt(500),
		runtime.NewInt(200),
		runtime.NewInt(299),
		runtime.NewString("status outside range"),
	)
}

func TestBetweenAssertionRejectsInvalidUsage(t *testing.T) {
	t.Parallel()

	for _, positive := range []bool{true, false} {
		fn := betweenAssertion.negative()
		if positive {
			fn = betweenAssertion.positive()
		}

		out, err := fn(context.Background(), runtime.NewInt(15), runtime.NewInt(20), runtime.NewInt(10))
		if out != runtime.None {
			t.Fatalf("reversed bounds output = %v, want None", out)
		}
		if !errors.Is(err, runtime.ErrInvalidArgument) || errors.Is(err, errAssertion) {
			t.Fatalf("reversed bounds error = %v, want invalid argument without assertion failure", err)
		}
	}

	duration := runtime.NewDuration(time.Second)
	out, err := betweenAssertion.positive()(
		context.Background(),
		runtime.NewString("1s"),
		duration,
		runtime.NewDuration(2*time.Second),
	)
	if out != runtime.None {
		t.Fatalf("incomparable value output = %v, want None", out)
	}
	if !errors.Is(err, runtime.ErrInvalidOperation) || errors.Is(err, errAssertion) {
		t.Fatalf("incomparable value error = %v, want invalid operation", err)
	}
}
