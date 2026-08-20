package runtime

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/internal/operator"
)

// Add adds native numeric or temporal values, concatenates when either operand
// is a String, or negotiates an unsupported pair through Addable host values.
func Add(ctx context.Context, left, right Value) (Value, error) {
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
		return dispatchAdd(ctx, left, right)
	case leftIsDuration && rightIsDuration:
		return addDurations(leftDuration, rightDuration)
	case leftIsDuration || rightIsDuration:
		return dispatchAdd(ctx, left, right)
	default:
		return dispatchAdd(ctx, left, right)
	}
}

// Subtract subtracts native numeric or compatible temporal values, or negotiates
// an unsupported pair through Subtractable host values.
func Subtract(ctx context.Context, left, right Value) (Value, error) {
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
		return dispatchSubtract(ctx, left, right)
	case leftIsDuration && rightIsDuration:
		return subtractDurations(leftDuration, rightDuration)
	case leftIsDuration || rightIsDuration:
		return dispatchSubtract(ctx, left, right)
	default:
		leftNumber, leftIsNumber := classifyNativeNumber(left)
		rightNumber, rightIsNumber := classifyNativeNumber(right)
		if leftIsNumber && rightIsNumber {
			return subtractNumbers(leftNumber, rightNumber)
		}

		return dispatchSubtract(ctx, left, right)
	}
}

// Multiply multiplies native numeric values, scales a Duration by a native number,
// or negotiates an unsupported pair through Multipliable host values.
func Multiply(ctx context.Context, left, right Value) (Value, error) {
	_, leftIsDateTime := left.(DateTime)
	_, rightIsDateTime := right.(DateTime)
	leftDuration, leftIsDuration := left.(Duration)
	rightDuration, rightIsDuration := right.(Duration)

	switch {
	case leftIsDateTime || rightIsDateTime:
		return dispatchMultiply(ctx, left, right)
	case leftIsDuration && rightIsDuration:
		return dispatchMultiply(ctx, left, right)
	case leftIsDuration:
		if !isDurationScalar(right) {
			return dispatchMultiply(ctx, left, right)
		}

		return multiplyDuration(leftDuration, right)
	case rightIsDuration:
		if !isDurationScalar(left) {
			return dispatchMultiply(ctx, left, right)
		}

		return multiplyDuration(rightDuration, left)
	default:
		leftNumber, leftIsNumber := classifyNativeNumber(left)
		rightNumber, rightIsNumber := classifyNativeNumber(right)
		if leftIsNumber && rightIsNumber {
			return multiplyNumbers(leftNumber, rightNumber)
		}

		return dispatchMultiply(ctx, left, right)
	}
}

// Divide divides native numeric values or a Duration by a compatible operand,
// or negotiates an unsupported pair through Dividable host values.
func Divide(ctx context.Context, left, right Value) (Value, error) {
	_, leftIsDateTime := left.(DateTime)
	_, rightIsDateTime := right.(DateTime)
	leftDuration, leftIsDuration := left.(Duration)
	rightDuration, rightIsDuration := right.(Duration)

	switch {
	case leftIsDateTime || rightIsDateTime:
		return dispatchDivide(ctx, left, right)
	case leftIsDuration && rightIsDuration:
		return divideDurations(leftDuration, rightDuration)
	case leftIsDuration:
		if !isDurationScalar(right) {
			return dispatchDivide(ctx, left, right)
		}

		return divideDuration(leftDuration, right)
	case rightIsDuration:
		return dispatchDivide(ctx, left, right)
	default:
		leftNumber, leftIsNumber := classifyNativeNumber(left)
		rightNumber, rightIsNumber := classifyNativeNumber(right)
		if leftIsNumber && rightIsNumber {
			return divideNumbers(leftNumber, rightNumber)
		}

		return dispatchDivide(ctx, left, right)
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

// Modulo computes the remainder of native numeric values, or negotiates an
// unsupported pair through Modulable host values.
func Modulo(ctx context.Context, left, right Value) (Value, error) {
	switch left.(type) {
	case DateTime, Duration:
		return dispatchMod(ctx, left, right)
	}

	switch right.(type) {
	case DateTime, Duration:
		return dispatchMod(ctx, left, right)
	}

	leftNumber, leftIsNumber := classifyNativeNumber(left)
	rightNumber, rightIsNumber := classifyNativeNumber(right)
	if leftIsNumber && rightIsNumber {
		return moduloNumbers(leftNumber, rightNumber)
	}

	return dispatchMod(ctx, left, right)
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
