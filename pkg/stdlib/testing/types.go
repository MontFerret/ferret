package testing

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

var (
	// Tests whether the value is true.
	// @param actual {Any} Value to test.
	// @param message {String} Message to display on error.
	// @return {Boolean} True when the configured assertion succeeds; otherwise an assertion error is returned.
	trueAssertion = assertion{
		defaultMessage: func(_ []runtime.Value) string {
			return fmt.Sprintf("be %s", formatValue(runtime.True))
		},
		args: assertionArgs{
			min: 1,
			max: 2,
		},
		fn: func(_ context.Context, args []runtime.Value) (bool, error) {
			return args[0] == runtime.True, nil
		},
	}

	// Tests whether the value is false.
	// @param actual {Any} Value to test.
	// @param message {String} Message to display on error.
	// @return {Boolean} True when the configured assertion succeeds; otherwise an assertion error is returned.
	falseAssertion = assertion{
		defaultMessage: func(_ []runtime.Value) string {
			return fmt.Sprintf("be %s", formatValue(runtime.False))
		},
		args: assertionArgs{
			min: 1,
			max: 2,
		},
		fn: func(_ context.Context, args []runtime.Value) (bool, error) {
			return args[0] == runtime.False, nil
		},
	}

	// Tests whether the value is None.
	// @param actual {Any} Value to test.
	// @param message {String} Message to display on error.
	// @return {Boolean} True when the configured assertion succeeds; otherwise an assertion error is returned.
	noneAssertion = newTypeAssertion(runtime.TypeNone)

	// Tests whether the value has the String type.
	// @param actual {Any} Value to test.
	// @param message {String} Message to display on error.
	// @return {Boolean} True when the configured assertion succeeds; otherwise an assertion error is returned.
	stringAssertion = newTypeAssertion(runtime.TypeString)

	// Tests whether the value has the Int type.
	// @param actual {Any} Actual value.
	// @param message {String} Message to display on error.
	// @return {Boolean} True when the configured assertion succeeds; otherwise an assertion error is returned.
	intAssertion = newTypeAssertion(runtime.TypeInt)

	// Tests whether the value has the Float type.
	// @param actual {Any} Value to test.
	// @param message {String} Message to display on error.
	// @return {Boolean} True when the configured assertion succeeds; otherwise an assertion error is returned.
	floatAssertion = newTypeAssertion(runtime.TypeFloat)

	// Tests whether the value has the DateTime type.
	// @param actual {Any} Value to test.
	// @param message {String} Message to display on error.
	// @return {Boolean} True when the configured assertion succeeds; otherwise an assertion error is returned.
	dateTimeAssertion = newTypeAssertion(runtime.TypeDateTime)

	// Tests whether the value has the array type.
	// @param actual {Any} Value to test.
	// @param message {String} Message to display on error.
	// @return {Boolean} True when the configured assertion succeeds; otherwise an assertion error is returned.
	arrayAssertion = newTypeAssertion(runtime.TypeArray)

	// Tests whether the value has the Object type.
	// @param actual {Any} Value to test.
	// @param message {String} Message to display on error.
	// @return {Boolean} True when the configured assertion succeeds; otherwise an assertion error is returned.
	objectAssertion = newTypeAssertion(runtime.TypeObject)

	// Tests whether the value has the binary type.
	// @param actual {Any} Value to test.
	// @param message {String} Message to display on error.
	// @return {Boolean} True when the configured assertion succeeds; otherwise an assertion error is returned.
	binaryAssertion = newTypeAssertion(runtime.TypeBinary)
)
