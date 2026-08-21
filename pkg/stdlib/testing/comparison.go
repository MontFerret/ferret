package testing

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type comparisonOperator int

const (
	comparisonEqual comparisonOperator = iota + 1
	comparisonLess
	comparisonLessOrEqual
	comparisonGreater
	comparisonGreaterOrEqual
)

var (
	// Tests equality of the actual and expected values.
	// @param actual {Any} Actual value.
	// @param expected {Any} Expected value.
	// @param message {String} Message to display on error.
	// @return {Boolean} True when the configured assertion succeeds; otherwise an assertion error is returned.
	equalAssertion = newComparisonAssertion(comparisonEqual)

	// Tests whether the actual value is greater than the expected value.
	// @param actual {Any} Actual value.
	// @param expected {Any} Expected value.
	// @param message {String} Message to display on error.
	// @return {Boolean} True when the configured assertion succeeds; otherwise an assertion error is returned.
	gtAssertion = newComparisonAssertion(comparisonGreater)

	// Tests whether the actual value is greater than or equal to the expected value.
	// @param actual {Any} Actual value.
	// @param expected {Any} Expected value.
	// @param message {String} Message to display on error.
	// @return {Boolean} True when the configured assertion succeeds; otherwise an assertion error is returned.
	gteAssertion = newComparisonAssertion(comparisonGreaterOrEqual)

	// Tests whether the actual value is less than the expected value.
	// @param actual {Any} Actual value.
	// @param expected {Any} Expected value.
	// @param message {String} Message to display on error.
	// @return {Boolean} True when the configured assertion succeeds; otherwise an assertion error is returned.
	ltAssertion = newComparisonAssertion(comparisonLess)

	// Tests whether the actual value is less than or equal to the expected value.
	// @param actual {Any} Actual value.
	// @param expected {Any} Expected value.
	// @param message {String} Message to display on error.
	// @return {Boolean} True when the configured assertion succeeds; otherwise an assertion error is returned.
	lteAssertion = newComparisonAssertion(comparisonLessOrEqual)
)

func (op comparisonOperator) String() string {
	switch op {
	case comparisonEqual:
		return "equal to"
	case comparisonLess:
		return "less than"
	case comparisonLessOrEqual:
		return "less than or equal to"
	case comparisonGreater:
		return "greater than"
	default:
		return "greater than or equal to"
	}
}

func (op comparisonOperator) compare(ctx context.Context, args []runtime.Value) (bool, error) {
	err := runtime.ValidateArgs(args, 2, 3)
	if err != nil {
		return false, err
	}

	actual := args[0]
	expected := args[1]

	switch op {
	case comparisonEqual:
		equal, err := runtime.EqualValues(ctx, actual, expected)

		return bool(equal), err
	case comparisonLess:
		result, err := runtime.CompareValues(ctx, actual, expected)

		return result < 0, err
	case comparisonLessOrEqual:
		result, err := runtime.CompareValues(ctx, actual, expected)

		return result <= 0, err
	case comparisonGreater:
		result, err := runtime.CompareValues(ctx, actual, expected)

		return result > 0, err
	default:
		result, err := runtime.CompareValues(ctx, actual, expected)

		return result >= 0, err
	}
}
