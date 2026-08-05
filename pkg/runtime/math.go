package runtime

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/internal/operator"
)

// Add applies Ferret addition without implicitly converting temporal operands.
func Add(ctx context.Context, left, right Value) (Value, error) {
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
		switch left := left.(type) {
		case Int:
			switch right := right.(type) {
			case Int:
				return addInts(left, right)
			case Float:
				return addFloats(Float(left), right, "addition")
			}
		case Float:
			switch right := right.(type) {
			case Int:
				return addFloats(left, Float(right), "addition")
			case Float:
				return addFloats(left, right, "addition")
			}
		}

		return addNonTemporal(ctx, left, right)
	}
}

// Subtract applies Ferret subtraction without implicitly converting temporal operands.
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
		return None, binaryOperatorTypeError(operator.Subtract, left, right)
	case leftIsDuration && rightIsDuration:
		return subtractDurations(leftDuration, rightDuration)
	case leftIsDuration || rightIsDuration:
		return None, binaryOperatorTypeError(operator.Subtract, left, right)
	default:
		return subtractNonTemporal(ctx, left, right)
	}
}

// Multiply applies Ferret multiplication without implicitly converting temporal operands.
func Multiply(ctx context.Context, left, right Value) (Value, error) {
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
		return multiplyNonTemporal(ctx, left, right)
	}
}

// Divide applies Ferret division without implicitly converting temporal operands.
func Divide(ctx context.Context, left, right Value) (Value, error) {
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
		return divideNonTemporal(ctx, left, right)
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

// Modulo applies Ferret modulo and rejects temporal operands.
func Modulo(ctx context.Context, left, right Value) (Value, error) {
	switch left.(type) {
	case DateTime, Duration:
		return None, binaryOperatorTypeError(operator.Modulus, left, right)
	}

	switch right.(type) {
	case DateTime, Duration:
		return None, binaryOperatorTypeError(operator.Modulus, left, right)
	}

	return moduloNonTemporal(ctx, left, right)
}

// Increment applies Ferret numeric increment and rejects temporal operands.
func Increment(ctx context.Context, value Value) (Value, error) {
	switch value.(type) {
	case DateTime, Duration:
		return None, unaryOperatorTypeError(operator.Increment, value)
	}

	return incrementNonTemporal(ctx, value)
}

// Decrement applies Ferret numeric decrement and rejects temporal operands.
func Decrement(ctx context.Context, value Value) (Value, error) {
	switch value.(type) {
	case DateTime, Duration:
		return None, unaryOperatorTypeError(operator.Decrement, value)
	}

	return decrementNonTemporal(ctx, value)
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
