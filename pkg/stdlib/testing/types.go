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
	// @return {None} No value is produced when the configured assertion succeeds.
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
	// @return {None} No value is produced when the configured assertion succeeds.
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
	// @return {None} No value is produced when the configured assertion succeeds.
	noneAssertion = newTypeAssertion(runtime.TypeNone)

	// Tests whether the value has the Boolean type.
	// @param actual {Any} Value to test.
	// @param message {String} Message to display on error.
	// @return {None} No value is produced when the configured assertion succeeds.
	boolAssertion = newTypeAssertion(runtime.TypeBoolean, "a boolean")

	// Tests whether the value has the String type.
	// @param actual {Any} Value to test.
	// @param message {String} Message to display on error.
	// @return {None} No value is produced when the configured assertion succeeds.
	stringAssertion = newTypeAssertion(runtime.TypeString)

	// Tests whether the value has the Int type.
	// @param actual {Any} Actual value.
	// @param message {String} Message to display on error.
	// @return {None} No value is produced when the configured assertion succeeds.
	intAssertion = newTypeAssertion(runtime.TypeInt)

	// Tests whether the value has the Float type.
	// @param actual {Any} Value to test.
	// @param message {String} Message to display on error.
	// @return {None} No value is produced when the configured assertion succeeds.
	floatAssertion = newTypeAssertion(runtime.TypeFloat)

	// Tests whether the value belongs to the numeric family.
	// @param actual {Any} Value to test.
	// @param message {String} Message to display on error.
	// @return {None} No value is produced when the configured assertion succeeds.
	numberAssertion = assertion{
		defaultMessage: func(_ []runtime.Value) string {
			return "be a number"
		},
		args: assertionArgs{
			min: 1,
			max: 2,
		},
		fn: func(_ context.Context, args []runtime.Value) (bool, error) {
			return bool(runtime.IsNumber(args[0])), nil
		},
	}

	// Tests whether the value has the Duration type.
	// @param actual {Any} Value to test.
	// @param message {String} Message to display on error.
	// @return {None} No value is produced when the configured assertion succeeds.
	durationAssertion = newTypeAssertion(runtime.TypeDuration, "a duration")

	// Tests whether the value has the DateTime type.
	// @param actual {Any} Value to test.
	// @param message {String} Message to display on error.
	// @return {None} No value is produced when the configured assertion succeeds.
	dateTimeAssertion = newTypeAssertion(runtime.TypeDateTime)

	// Tests whether the value has the array type.
	// @param actual {Any} Value to test.
	// @param message {String} Message to display on error.
	// @return {None} No value is produced when the configured assertion succeeds.
	arrayAssertion = newTypeAssertion(runtime.TypeArray)

	// Tests whether the value has the Object type.
	// @param actual {Any} Value to test.
	// @param message {String} Message to display on error.
	// @return {None} No value is produced when the configured assertion succeeds.
	objectAssertion = newTypeAssertion(runtime.TypeObject)

	// Tests whether the value has the binary type.
	// @param actual {Any} Value to test.
	// @param message {String} Message to display on error.
	// @return {None} No value is produced when the configured assertion succeeds.
	binaryAssertion = newTypeAssertion(runtime.TypeBinary)
)
