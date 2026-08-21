package testing

import (
	"testing"

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
			failure:  "assertion error: expected String 'Foo' to be equal to String 'Bar'",
			negative: "assertion error: expected String 'Bar' not to be equal to String 'Bar'",
		},
		{
			name:     "escaped string",
			actual:   runtime.NewString(`can't\stop`),
			expected: runtime.NewString("won't"),
			failure:  `assertion error: expected String 'can\'t\\stop' to be equal to String 'won\'t'`,
			negative: `assertion error: expected String 'won\'t' not to be equal to String 'won\'t'`,
		},
		{
			name:     "int",
			actual:   runtime.NewInt(1),
			expected: runtime.NewInt(2),
			failure:  "assertion error: expected Int '1' to be equal to Int '2'",
			negative: "assertion error: expected Int '2' not to be equal to Int '2'",
		},
		{
			name:     "boolean",
			actual:   runtime.False,
			expected: runtime.True,
			failure:  "assertion error: expected Boolean 'false' to be equal to Boolean 'true'",
			negative: "assertion error: expected Boolean 'true' not to be equal to Boolean 'true'",
		},
		{
			name:     "array",
			actual:   one,
			expected: oneTwo,
			failure:  "assertion error: expected Array '[1]' to be equal to Array '[1,2]'",
			negative: "assertion error: expected Array '[1,2]' not to be equal to Array '[1,2]'",
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
