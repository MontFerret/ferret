package testing

import (
	"context"
	"fmt"

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
	// @return {None} No value is produced when the configured assertion succeeds.
	equalAssertion = newComparisonAssertion(comparisonEqual, equalityFailureMessage)

	// Tests whether the actual value is greater than the expected value.
	// @param actual {Any} Actual value.
	// @param expected {Any} Expected value.
	// @param message {String} Message to display on error.
	// @return {None} No value is produced when the configured assertion succeeds.
	gtAssertion = newComparisonAssertion(comparisonGreater, nil)

	// Tests whether the actual value is greater than or equal to the expected value.
	// @param actual {Any} Actual value.
	// @param expected {Any} Expected value.
	// @param message {String} Message to display on error.
	// @return {None} No value is produced when the configured assertion succeeds.
	gteAssertion = newComparisonAssertion(comparisonGreaterOrEqual, nil)

	// Tests whether the actual value is less than the expected value.
	// @param actual {Any} Actual value.
	// @param expected {Any} Expected value.
	// @param message {String} Message to display on error.
	// @return {None} No value is produced when the configured assertion succeeds.
	ltAssertion = newComparisonAssertion(comparisonLess, nil)

	// Tests whether the actual value is less than or equal to the expected value.
	// @param actual {Any} Actual value.
	// @param expected {Any} Expected value.
	// @param message {String} Message to display on error.
	// @return {None} No value is produced when the configured assertion succeeds.
	lteAssertion = newComparisonAssertion(comparisonLessOrEqual, nil)

	// Tests whether two numeric values are within an inclusive tolerance.
	// @param actual {Int | Float} Actual numeric value.
	// @param expected {Int | Float} Expected numeric value.
	// @param tolerance {Int | Float} Maximum non-negative difference.
	// @param message {String} Message to display on error.
	// @return {None} No value is produced when the configured assertion succeeds.
	approxAssertion = assertion{
		defaultMessage: func(ctx context.Context, args []runtime.Value) string {
			return fmt.Sprintf(
				"approximately equal %s within %s",
				formatValue(ctx, args[1]),
				formatValue(ctx, args[2]),
			)
		},
		args: assertionArgs{
			min: 3,
			max: 4,
		},
		fn: func(ctx context.Context, args []runtime.Value) (bool, error) {
			for pos, arg := range args[:3] {
				if !runtime.IsNumber(arg) {
					return false, runtime.ArgTypeError(arg, pos, runtime.TypeInt, runtime.TypeFloat)
				}
			}

			tolerance := args[2]
			comparison, err := runtime.CompareValues(ctx, tolerance, runtime.ZeroInt)
			if err != nil {
				return false, err
			}

			if comparison < runtime.Equal {
				return false, runtime.ArgError(
					runtime.Error(runtime.ErrInvalidArgument, "tolerance must be non-negative"),
					2,
				)
			}

			negativeTolerance, err := runtime.Negative(tolerance)
			if err != nil {
				return false, err
			}

			difference, err := runtime.Subtract(ctx, args[0], args[1])
			if err != nil {
				return false, err
			}

			lower, err := runtime.CompareValues(ctx, difference, negativeTolerance)
			if err != nil {
				return false, err
			}

			if lower < runtime.Equal {
				return false, nil
			}

			upper, err := runtime.CompareValues(ctx, difference, tolerance)

			return upper <= runtime.Equal, err
		},
	}

	// Tests whether a value is within inclusive minimum and maximum boundaries.
	// @param value {Any} Value to test.
	// @param min {Any} Inclusive minimum boundary.
	// @param max {Any} Inclusive maximum boundary.
	// @param message {String} Message to display on error.
	// @return {None} No value is produced when the configured assertion succeeds.
	betweenAssertion = assertion{
		defaultMessage: func(ctx context.Context, args []runtime.Value) string {
			return fmt.Sprintf(
				"be between %s and %s",
				formatValue(ctx, args[1]),
				formatValue(ctx, args[2]),
			)
		},
		args: assertionArgs{
			min: 3,
			max: 4,
		},
		fn: func(ctx context.Context, args []runtime.Value) (bool, error) {
			bounds, err := runtime.CompareValues(ctx, args[1], args[2])
			if err != nil {
				return false, err
			}

			if bounds > runtime.Equal {
				return false, runtime.Error(
					runtime.ErrInvalidArgument,
					"minimum boundary must not exceed maximum boundary",
				)
			}

			lower, err := runtime.CompareValues(ctx, args[0], args[1])
			if err != nil {
				return false, err
			}

			if lower < runtime.Equal {
				return false, nil
			}

			upper, err := runtime.CompareValues(ctx, args[0], args[2])

			return upper <= runtime.Equal, err
		},
	}
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
