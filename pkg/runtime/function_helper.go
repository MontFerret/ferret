package runtime

import "fmt"

// ArgError creates an error for an invalid argument at the specified position.
// The position is 0-based internally but reported as 1-based to users.
func ArgError(err error, pos int) error {
	return newInvalidArgumentError(err, pos)
}

// ArgTypeError returns an error indicating that the argument at the specified position does not match the expected type.
// The position is 0-based internally but reported as 1-based to users.
func ArgTypeError(arg Value, pos int, expected ...Type) error {
	return ArgError(TypeErrorOf(arg, expected...), pos)
}

// ArityError returns an error if the number of arguments is outside the [minimum, maximum] range.
// The minimum and maximum values are inclusive.
func ArityError(count, minimum, maximum int) error {
	if count < minimum || count > maximum {
		var num string

		if minimum == maximum {
			num = fmt.Sprintf("%d", minimum)
		} else {
			num = fmt.Sprintf("%d-%d", minimum, maximum)
		}

		return Error(
			ErrInvalidArgumentNumber,
			fmt.Sprintf(
				"expected number of arguments %s, but got %d",
				num,
				count))
	}

	return nil
}

// ValidateArgs validates that the number of arguments is within the specified range.
// It returns an error if the argument count is outside the [minimum, maximum] range.
func ValidateArgs(args []Value, minimum, maximum int) error {
	return ArityError(len(args), minimum, maximum)
}

// ValidateArgsType validates that each argument in the provided slice matches at least one of the expected types.
// It returns an error if any argument does not match any of the expected types, indicating the position of the invalid argument.
// The position is 0-based internally but reported as 1-based to users.
func ValidateArgsType(args []Value, expected ...Type) error {
	for pos, arg := range args {
		var matched bool

		for _, targetType := range expected {
			if targetType.Is(arg) {
				matched = true
				break
			}
		}

		if !matched {
			return ArgError(TypeErrorOf(arg, expected...), pos)
		}
	}

	return nil
}

// ValidateArgTypeAt checks if the argument at the specified position matches one of the expected types.
// Returns an error if the position is out of bounds or if the argument does not match any of the expected types.
func ValidateArgTypeAt(args []Value, pos int, expected ...Type) error {
	if pos >= len(args) {
		return errArgOutOfBounds(args, pos)
	}

	arg := args[pos]

	for _, targetType := range expected {
		if targetType.Is(arg) {
			return nil
		}
	}

	return ArgError(TypeErrorOf(arg, expected...), pos)
}

// ValidateArgType validates that the provided argument matches one of the expected types.
// It returns an error if the argument does not match any of the expected types, indicating the position of the invalid argument.
func ValidateArgType(arg Value, pos int, expected ...Type) error {
	for _, targetType := range expected {
		if targetType.Is(arg) {
			return nil
		}
	}

	return ArgError(TypeErrorOf(arg, expected...), pos)
}

// ValidateArgsValue validates that each argument in the provided slice satisfies the given value assertion.
// It returns an error if any argument does not satisfy the value assertion, indicating the position of the invalid argument.
func ValidateArgsValue(args []Value, assertion ValueAssertion) error {
	for pos, arg := range args {
		err := assertion(arg)

		if err != nil {
			return ArgError(err, pos)
		}
	}

	return nil
}

// ValidateArgValueAt checks if the argument at the specified position satisfies the given assertion function.
// Returns an error if the position is out of bounds or if the assertion fails.
func ValidateArgValueAt(args []Value, pos int, assertion ValueAssertion) error {
	if pos >= len(args) {
		return errArgOutOfBounds(args, pos)
	}

	arg := args[pos]
	err := assertion(arg)

	if err == nil {
		return nil
	}

	return ArgError(err, pos)
}

// ValidateArgValue asserts that each argument in the provided slice satisfies the given value assertion.
// It returns an error if any argument does not satisfy the value assertion, indicating the position of the invalid argument.
func ValidateArgValue(arg Value, pos int, assertion ValueAssertion) error {
	err := assertion(arg)

	if err == nil {
		return nil
	}

	return ArgError(err, pos)
}

// CastArg attempts to cast a given Value to the specified type T.
// Returns the casted value and nil error if the type assertion succeeds, otherwise the zero value of T and an error.
func CastArg[T Value](arg Value, index int) (T, error) {
	val, ok := arg.(T)

	if ok {
		return val, nil
	}

	var zero T

	return zero, ArgError(TypeErrorOf(arg, expectedTypeOf[T]()), index)
}

// CastArgAt attempts to cast an argument from the args slice at the specified index to type T.
// Returns the casted value and nil error if successful, otherwise the zero value of type T and an error.
func CastArgAt[T Value](args []Value, index int) (T, error) {
	if index >= len(args) {
		var zero T

		return zero, errArgOutOfBounds(args, index)
	}

	return CastArg[T](args[index], index)
}

// CastArgAtOr attempts to cast the argument at the specified index to type T.
// Returns the casted value if successful or the provided fallback value if the index is out of range.
func CastArgAtOr[T Value](args []Value, index int, fallback T) (T, error) {
	if index >= len(args) {
		return fallback, nil
	}

	return CastArg[T](args[index], index)
}

// CastArgs attempts to cast a slice of Value to a slice of a specified type T.
// Returns the slice of casted values and nil error if all casts succeed, otherwise nil and an error.
func CastArgs[T Value](arg []Value) ([]T, error) {
	casted := make([]T, len(arg))

	for i, a := range arg {
		val, err := CastArg[T](a, i)

		if err != nil {
			return nil, err
		}

		casted[i] = val
	}

	return casted, nil
}

// CastArgs2 attempts to cast two Values to the specified types T1 and T2.
// Returns the casted values and nil error if the type assertions succeed;
// otherwise, returns zero values of T1 and T2 along with an error.
func CastArgs2[T1, T2 Value](arg1, arg2 Value) (T1, T2, error) {
	val1, err := CastArg[T1](arg1, 0)

	if err != nil {
		var zero1 T1
		var zero2 T2

		return zero1, zero2, err
	}

	val2, err := CastArg[T2](arg2, 1)

	if err != nil {
		var zero1 T1
		var zero2 T2

		return zero1, zero2, err
	}

	return val1, val2, nil
}

// CastArgs3 attempts to cast three Values to the specified types T1, T2, and T3.
// Returns the casted values and nil error if the type assertions succeed;
// otherwise, returns zero values of T1, T2, and T3 along with an error.
func CastArgs3[T1, T2, T3 Value](arg1, arg2, arg3 Value) (T1, T2, T3, error) {
	val1, err := CastArg[T1](arg1, 0)

	if err != nil {
		var zero1 T1
		var zero2 T2
		var zero3 T3

		return zero1, zero2, zero3, err
	}

	val2, err := CastArg[T2](arg2, 1)

	if err != nil {
		var zero1 T1
		var zero2 T2
		var zero3 T3

		return zero1, zero2, zero3, err
	}

	val3, err := CastArg[T3](arg3, 2)

	if err != nil {
		var zero1 T1
		var zero2 T2
		var zero3 T3

		return zero1, zero2, zero3, err
	}

	return val1, val2, val3, nil
}

// CastArgs4 attempts to cast four Values to the specified types T1, T2, T3, and T4.
// Returns the casted values and nil error if the type assertions succeed;
// otherwise, returns zero values of T1, T2, T3, and T4 along with an error.
func CastArgs4[T1, T2, T3, T4 Value](arg1, arg2, arg3, arg4 Value) (T1, T2, T3, T4, error) {
	val1, err := CastArg[T1](arg1, 0)

	if err != nil {
		var zero1 T1
		var zero2 T2
		var zero3 T3
		var zero4 T4

		return zero1, zero2, zero3, zero4, err
	}

	val2, err := CastArg[T2](arg2, 1)

	if err != nil {
		var zero1 T1
		var zero2 T2
		var zero3 T3
		var zero4 T4

		return zero1, zero2, zero3, zero4, err
	}

	val3, err := CastArg[T3](arg3, 2)

	if err != nil {
		var zero1 T1
		var zero2 T2
		var zero3 T3
		var zero4 T4

		return zero1, zero2, zero3, zero4, err
	}

	val4, err := CastArg[T4](arg4, 3)

	if err != nil {
		var zero1 T1
		var zero2 T2
		var zero3 T3
		var zero4 T4

		return zero1, zero2, zero3, zero4, err
	}

	return val1, val2, val3, val4, nil
}

// CastVarArgs2 attempts to cast the first two elements of the args slice to the specified types T1 and T2.
// It first validates that the number of arguments is between 2 and MaxArgs. If validation fails, it returns zero values of T1 and T2 along with the error.
// If validation succeeds, it attempts to cast the first two arguments to T1 and T2 respectively, returning any errors encountered during casting.
func CastVarArgs2[T1, T2 Value](args []Value) (T1, T2, error) {
	if err := ValidateArgs(args, 2, MaxArgs); err != nil {
		var zero1 T1
		var zero2 T2

		return zero1, zero2, err
	}

	val1, err := CastArgAt[T1](args, 0)

	if err != nil {
		var zero1 T1
		var zero2 T2

		return zero1, zero2, err
	}

	val2, err := CastArgAt[T2](args, 1)

	if err != nil {
		var zero1 T1
		var zero2 T2

		return zero1, zero2, err
	}

	return val1, val2, nil
}

// CastVarArgs3 attempts to cast the first three elements of the args slice to the specified types T1, T2, and T3.
// It first validates that the number of arguments is between 3 and MaxArgs. If validation fails, it returns zero values of T1, T2, and T3 along with the error.
// If validation succeeds, it attempts to cast the first three arguments to T1, T2, and T3 respectively, returning any errors encountered during casting.
func CastVarArgs3[T1, T2, T3 Value](args []Value) (T1, T2, T3, error) {
	if err := ValidateArgs(args, 3, MaxArgs); err != nil {
		var zero1 T1
		var zero2 T2
		var zero3 T3

		return zero1, zero2, zero3, err
	}

	val1, err := CastArgAt[T1](args, 0)

	if err != nil {
		var zero1 T1
		var zero2 T2
		var zero3 T3

		return zero1, zero2, zero3, err
	}

	val2, err := CastArgAt[T2](args, 1)

	if err != nil {
		var zero1 T1
		var zero2 T2
		var zero3 T3

		return zero1, zero2, zero3, err
	}

	val3, err := CastArgAt[T3](args, 2)

	if err != nil {
		var zero1 T1
		var zero2 T2
		var zero3 T3

		return zero1, zero2, zero3, err
	}

	return val1, val2, val3, nil
}

// CastVarArgs4 attempts to cast the first four elements of the args slice to the specified types T1, T2, T3, and T4.
// It first validates that the number of arguments is between 4 and MaxArgs. If validation fails, it returns zero values of T1, T2, T3, and T4 along with the error.
// If validation succeeds, it attempts to cast the first four arguments to T1, T2, T3, and T4 respectively, returning any errors encountered during casting.
func CastVarArgs4[T1, T2, T3, T4 Value](args []Value) (T1, T2, T3, T4, error) {
	if err := ValidateArgs(args, 4, MaxArgs); err != nil {
		var zero1 T1
		var zero2 T2
		var zero3 T3
		var zero4 T4

		return zero1, zero2, zero3, zero4, err
	}

	val1, err := CastArgAt[T1](args, 0)

	if err != nil {
		var zero1 T1
		var zero2 T2
		var zero3 T3
		var zero4 T4

		return zero1, zero2, zero3, zero4, err
	}

	val2, err := CastArgAt[T2](args, 1)

	if err != nil {
		var zero1 T1
		var zero2 T2
		var zero3 T3
		var zero4 T4

		return zero1, zero2, zero3, zero4, err
	}

	val3, err := CastArgAt[T3](args, 2)

	if err != nil {
		var zero1 T1
		var zero2 T2
		var zero3 T3
		var zero4 T4

		return zero1, zero2, zero3, zero4, err
	}

	val4, err := CastArgAt[T4](args, 3)

	if err != nil {
		var zero1 T1
		var zero2 T2
		var zero3 T3
		var zero4 T4

		return zero1, zero2, zero3, zero4, err
	}

	return val1, val2, val3, val4, nil
}

func errArgOutOfBounds(args []Value, pos int) error {
	return Error(ErrInvalidArgumentNumber, fmt.Sprintf("expected at least %d arguments, but got %d", pos+1, len(args)))
}
