package runtime

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/internal/operator"
)

// Add adds native numeric or temporal values, or concatenates when either operand is a String.
func Add(_ context.Context, left, right Value) (Value, error) {
	leftNumber, leftIsNumber := classifyNativeNumber(left)
	rightNumber, rightIsNumber := classifyNativeNumber(right)
	if leftIsNumber && rightIsNumber {
		return addNumbers(leftNumber, rightNumber)
	}

	if _, ok := left.(String); ok {
		return concatValues(left, right), nil
	}
	if _, ok := right.(String); ok {
		return concatValues(left, right), nil
	}

	leftDateTime, leftIsDateTime := left.(DateTime)
	rightDateTime, rightIsDateTime := right.(DateTime)
	leftDuration, leftIsDuration := left.(Duration)
	rightDuration, rightIsDuration := right.(Duration)

	switch {
	case leftIsDateTime && rightIsDuration:
		return addDateTimeDuration(leftDateTime, rightDuration)
	case leftIsDuration && rightIsDateTime:
		return addDateTimeDuration(rightDateTime, leftDuration)
	case leftIsDateTime || rightIsDateTime:
		return None, binaryOperatorTypeError(operator.Add, left, right)
	case leftIsDuration && rightIsDuration:
		return addDurations(leftDuration, rightDuration)
	case leftIsDuration || rightIsDuration:
		return None, binaryOperatorTypeError(operator.Add, left, right)
	default:
		return None, binaryOperatorTypeError(operator.Add, left, right)
	}
}

// Subtract subtracts native numeric or compatible temporal values.
func Subtract(_ context.Context, left, right Value) (Value, error) {
	leftDateTime, leftIsDateTime := left.(DateTime)
	rightDateTime, rightIsDateTime := right.(DateTime)
	leftDuration, leftIsDuration := left.(Duration)
	rightDuration, rightIsDuration := right.(Duration)

	switch {
	case leftIsDateTime && rightIsDateTime:
		return subtractDateTimes(leftDateTime, rightDateTime)
	case leftIsDateTime && rightIsDuration:
		return subtractDateTimeDuration(leftDateTime, rightDuration)
	case leftIsDateTime || rightIsDateTime:
		return None, binaryOperatorTypeError(operator.Subtract, left, right)
	case leftIsDuration && rightIsDuration:
		return subtractDurations(leftDuration, rightDuration)
	case leftIsDuration || rightIsDuration:
		return None, binaryOperatorTypeError(operator.Subtract, left, right)
	default:
		leftNumber, leftIsNumber := classifyNativeNumber(left)
		rightNumber, rightIsNumber := classifyNativeNumber(right)
		if leftIsNumber && rightIsNumber {
			return subtractNumbers(leftNumber, rightNumber)
		}

		return None, binaryOperatorTypeError(operator.Subtract, left, right)
	}
}

// Multiply multiplies native numeric values or scales a Duration by a native number.
func Multiply(_ context.Context, left, right Value) (Value, error) {
	_, leftIsDateTime := left.(DateTime)
	_, rightIsDateTime := right.(DateTime)
	leftDuration, leftIsDuration := left.(Duration)
	rightDuration, rightIsDuration := right.(Duration)

	switch {
	case leftIsDateTime || rightIsDateTime:
		return None, binaryOperatorTypeError(operator.Multiply, left, right)
	case leftIsDuration && rightIsDuration:
		return None, binaryOperatorTypeError(operator.Multiply, left, right)
	case leftIsDuration:
		if !isDurationScalar(right) {
			return None, binaryOperatorTypeError(operator.Multiply, left, right)
		}

		return multiplyDuration(leftDuration, right)
	case rightIsDuration:
		if !isDurationScalar(left) {
			return None, binaryOperatorTypeError(operator.Multiply, left, right)
		}

		return multiplyDuration(rightDuration, left)
	default:
		leftNumber, leftIsNumber := classifyNativeNumber(left)
		rightNumber, rightIsNumber := classifyNativeNumber(right)
		if leftIsNumber && rightIsNumber {
			return multiplyNumbers(leftNumber, rightNumber)
		}

		return None, binaryOperatorTypeError(operator.Multiply, left, right)
	}
}

// Divide divides native numeric values or a Duration by a compatible operand.
func Divide(_ context.Context, left, right Value) (Value, error) {
	_, leftIsDateTime := left.(DateTime)
	_, rightIsDateTime := right.(DateTime)
	leftDuration, leftIsDuration := left.(Duration)
	rightDuration, rightIsDuration := right.(Duration)

	switch {
	case leftIsDateTime || rightIsDateTime:
		return None, binaryOperatorTypeError(operator.Divide, left, right)
	case leftIsDuration && rightIsDuration:
		return divideDurations(leftDuration, rightDuration)
	case leftIsDuration:
		if !isDurationScalar(right) {
			return None, binaryOperatorTypeError(operator.Divide, left, right)
		}

		return divideDuration(leftDuration, right)
	case rightIsDuration:
		return None, binaryOperatorTypeError(operator.Divide, left, right)
	default:
		leftNumber, leftIsNumber := classifyNativeNumber(left)
		rightNumber, rightIsNumber := classifyNativeNumber(right)
		if leftIsNumber && rightIsNumber {
			return divideNumbers(leftNumber, rightNumber)
		}

		return None, binaryOperatorTypeError(operator.Divide, left, right)
	}
}

func isDurationScalar(value Value) bool {
	switch value.(type) {
	case Int, Float:
		return true
	default:
		return false
	}
}

// Modulo computes the remainder of native numeric values and rejects temporal operands.
func Modulo(_ context.Context, left, right Value) (Value, error) {
	switch left.(type) {
	case DateTime, Duration:
		return None, binaryOperatorTypeError(operator.Modulus, left, right)
	}

	switch right.(type) {
	case DateTime, Duration:
		return None, binaryOperatorTypeError(operator.Modulus, left, right)
	}

	leftNumber, leftIsNumber := classifyNativeNumber(left)
	rightNumber, rightIsNumber := classifyNativeNumber(right)
	if leftIsNumber && rightIsNumber {
		return moduloNumbers(leftNumber, rightNumber)
	}

	return None, binaryOperatorTypeError(operator.Modulus, left, right)
}

// Increment increments a native numeric value.
func Increment(_ context.Context, value Value) (Value, error) {
	switch value := value.(type) {
	case Int:
		return addInts(value, 1)
	case Float:
		return addFloats(value, 1, "increment")
	default:
		return None, unaryOperatorTypeError(operator.Increment, value)
	}
}

// Decrement decrements a native numeric value.
func Decrement(_ context.Context, value Value) (Value, error) {
	switch value := value.(type) {
	case Int:
		return subtractInts(value, 1)
	case Float:
		return addFloats(value, -1, "decrement")
	default:
		return None, unaryOperatorTypeError(operator.Decrement, value)
	}
}

// Not applies logical negation to a Boolean.
func Not(value Value) (Value, error) {
	boolean, ok := value.(Boolean)
	if !ok {
		return None, unaryOperatorTypeError(operator.Not, value)
	}

	return !boolean, nil
}

// Positive applies unary plus to numeric and Duration values.
func Positive(value Value) (Value, error) {
	switch value := value.(type) {
	case Int, Float, Duration:
		return value, nil
	default:
		return None, unaryOperatorTypeError(operator.Positive, value)
	}
}

// Negative applies unary minus to numeric and Duration values.
func Negative(value Value) (Value, error) {
	switch value := value.(type) {
	case Int:
		return negateInt(value)
	case Float:
		return negateFloat(value)
	case Duration:
		return negateDuration(value)
	default:
		return None, unaryOperatorTypeError(operator.Negative, value)
	}
}
