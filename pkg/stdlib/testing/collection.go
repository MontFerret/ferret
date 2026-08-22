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
	// @return {None} No value is produced when the configured assertion succeeds.
	emptyAssertion = assertion{
		defaultMessage: func(_ context.Context, _ []runtime.Value) string {
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

	// Tests whether the actual container contains the expected value.
	// @param actual {String | Array | Object | Iterable} Haystack value.
	// @param expected {Any} Expected value.
	// @param message {String} Message to display on error.
	// @return {None} No value is produced when the configured assertion succeeds.
	containsAssertion = assertion{
		defaultMessage: func(ctx context.Context, args []runtime.Value) string {
			return fmt.Sprintf("contain %s", formatValue(ctx, args[1]))
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
	// @return {None} No value is produced when the configured assertion succeeds.
	lenAssertion = assertion{
		defaultMessage: func(_ context.Context, args []runtime.Value) string {
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
	// @return {None} No value is produced when the configured assertion succeeds.
	matchAssertion = assertion{
		defaultMessage: func(ctx context.Context, args []runtime.Value) string {
			return fmt.Sprintf("match regular expression %s", formatValue(ctx, args[1]))
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

	// Tests whether a map has one property or all properties in a list.
	// @param object {Map} Object or map to inspect.
	// @param keys {String | [String]} Property name or property names to require.
	// @param message {String} Message to display on error.
	// @return {None} No value is produced when the configured assertion succeeds.
	hasAssertion = assertion{
		defaultMessage: func(ctx context.Context, args []runtime.Value) string {
			if runtime.TypeString.Is(args[1]) {
				return fmt.Sprintf("have property %s", formatValue(ctx, args[1]))
			}

			return fmt.Sprintf("have all properties %s", formatValue(ctx, args[1]))
		},
		failureMessage: hasFailureMessage,
		args: assertionArgs{
			min: 2,
			max: 3,
		},
		fn: func(ctx context.Context, args []runtime.Value) (bool, error) {
			if err := runtime.ValidateArgType(args[0], 0, runtime.TypeMap); err != nil {
				return false, err
			}

			if err := runtime.ValidateArgType(args[1], 1, runtime.TypeString, runtime.TypeList); err != nil {
				return false, err
			}

			target := args[0].(runtime.Map)
			if key, ok := args[1].(runtime.String); ok {
				contains, err := target.ContainsKey(ctx, key)

				return bool(contains), err
			}

			keys := args[1].(runtime.List)
			err := keys.ForEach(ctx, func(ctx context.Context, key runtime.Value, index runtime.Int) (runtime.Boolean, error) {
				if err := runtime.ValidateType(key, runtime.TypeString); err != nil {
					return false, runtime.ArgError(runtime.Errorf(err, "key at index %d", index), 1)
				}

				return true, nil
			})
			if err != nil {
				return false, err
			}

			allPresent := true
			err = keys.ForEach(ctx, func(ctx context.Context, key runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
				contains, err := target.ContainsKey(ctx, key)
				if err != nil {
					return false, err
				}

				if !contains {
					allPresent = false

					return false, nil
				}

				return true, nil
			})
			if err != nil {
				return false, err
			}

			return allPresent, nil
		},
	}
)
