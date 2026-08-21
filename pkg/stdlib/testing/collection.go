package testing

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/collections"
	stdlibstrings "github.com/MontFerret/ferret/v2/pkg/stdlib/strings"
)

var (
	// Tests whether the target is empty.
	// @param actual {Measurable | Binary | Object | Any[] | String} Value to test.
	// @param message {String} Message to display on error.
	// @return {Boolean} True when the configured assertion succeeds; otherwise an assertion error is returned.
	emptyAssertion = assertion{
		defaultMessage: func(_ []runtime.Value) string {
			return "be empty"
		},
		args: assertionArgs{
			min: 1,
			max: 2,
		},
		fn: func(ctx context.Context, args []runtime.Value) (bool, error) {
			value := args[0]
			if err := runtime.AssertMeasurable(value); err != nil {
				return false, err
			}

			return runtime.IsEmpty(ctx, value)
		},
	}

	// Tests whether the actual container includes the expected value.
	// @param actual {String | Array | Object | Iterable} Haystack value.
	// @param expected {Any} Expected value.
	// @param message {String} Message to display on error.
	// @return {Boolean} True when the configured assertion succeeds; otherwise an assertion error is returned.
	includeAssertion = assertion{
		defaultMessage: func(args []runtime.Value) string {
			return fmt.Sprintf("include %s", formatValue(args[1]))
		},
		args: assertionArgs{
			min: 2,
			max: 3,
		},
		fn: func(ctx context.Context, args []runtime.Value) (bool, error) {
			haystack := args[0]
			needle := args[1]

			out, err := collections.Includes(ctx, haystack, needle)
			if err != nil {
				return false, err
			}

			equal, err := runtime.EqualValues(ctx, out, runtime.True)

			return bool(equal), err
		},
	}

	// Tests whether a measurable value has the expected length or size.
	// @param actual {Measurable} Measurable value.
	// @param length {Int} Target length.
	// @param message {String} Message to display on error.
	// @return {Boolean} True when the configured assertion succeeds; otherwise an assertion error is returned.
	lenAssertion = assertion{
		defaultMessage: func(args []runtime.Value) string {
			return fmt.Sprintf("has size %s", args[1])
		},
		args: assertionArgs{
			min: 2,
			max: 3,
		},
		fn: func(ctx context.Context, args []runtime.Value) (bool, error) {
			collection := args[0]
			size := args[1]

			if err := runtime.AssertMeasurable(collection); err != nil {
				return false, err
			}

			out, err := runtime.Length(ctx, collection)
			if err != nil {
				return false, err
			}

			equal, err := runtime.EqualValues(ctx, out, size)

			return bool(equal), err
		},
	}

	// Tests whether the value matches the regular expression.
	// @param actual {Any} Actual value.
	// @param expression {String} Regular expression.
	// @param message {String} Message to display on error.
	// @return {Boolean} True when the configured assertion succeeds; otherwise an assertion error is returned.
	matchAssertion = assertion{
		defaultMessage: func(_ []runtime.Value) string {
			return "match regular expression"
		},
		args: assertionArgs{
			min: 2,
			max: 3,
		},
		fn: func(ctx context.Context, args []runtime.Value) (bool, error) {
			value := args[0]
			expression := args[1]

			out, err := stdlibstrings.RegexTest(ctx, value, expression)
			if err != nil {
				return false, err
			}

			equal, err := runtime.EqualValues(ctx, out, runtime.True)

			return bool(equal), err
		},
	}
)
