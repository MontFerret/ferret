package runtime

import (
	"context"
	"math"
	"math/big"
)

// AddChecked applies addition while enforcing native duration arithmetic rules.
func AddChecked(ctx context.Context, left, right Value) (Value, error) {
	leftDuration, leftIsDuration := left.(Duration)
	rightDuration, rightIsDuration := right.(Duration)

	if !leftIsDuration && !rightIsDuration {
		return Add(ctx, left, right), nil
	}

	if !leftIsDuration || !rightIsDuration {
		return None, durationBinaryTypeError("+", left, right)
	}

	result, ok := addDurationInt64(int64(leftDuration), int64(rightDuration))
	if !ok {
		return None, durationRangeError("addition")
	}

	return Duration(result), nil
}

// SubtractChecked applies subtraction while enforcing native duration arithmetic rules.
func SubtractChecked(ctx context.Context, left, right Value) (Value, error) {
	leftDuration, leftIsDuration := left.(Duration)
	rightDuration, rightIsDuration := right.(Duration)

	if !leftIsDuration && !rightIsDuration {
		return Subtract(ctx, left, right), nil
	}

	if !leftIsDuration || !rightIsDuration {
		return None, durationBinaryTypeError("-", left, right)
	}

	result, ok := subtractDurationInt64(int64(leftDuration), int64(rightDuration))
	if !ok {
		return None, durationRangeError("subtraction")
	}

	return Duration(result), nil
}

// MultiplyChecked applies multiplication while enforcing native duration arithmetic rules.
func MultiplyChecked(ctx context.Context, left, right Value) (Value, error) {
	leftDuration, leftIsDuration := left.(Duration)
	rightDuration, rightIsDuration := right.(Duration)

	if !leftIsDuration && !rightIsDuration {
		return Multiply(ctx, left, right), nil
	}

	if leftIsDuration && rightIsDuration {
		return None, durationBinaryTypeError("*", left, right)
	}

	if leftIsDuration {
		return scaleDuration(leftDuration, right, false)
	}

	return scaleDuration(rightDuration, left, false)
}

// DivideChecked applies division while enforcing native duration arithmetic rules.
func DivideChecked(ctx context.Context, left, right Value) (Value, error) {
	leftDuration, leftIsDuration := left.(Duration)
	rightDuration, rightIsDuration := right.(Duration)

	if !leftIsDuration && !rightIsDuration {
		return Divide(ctx, left, right), nil
	}

	if !leftIsDuration {
		return None, durationBinaryTypeError("/", left, right)
	}

	if rightIsDuration {
		if rightDuration == 0 {
			return None, Error(ErrInvalidOperation, "division by zero")
		}

		leftNanos := int64(leftDuration)
		rightNanos := int64(rightDuration)

		if leftNanos == math.MinInt64 && rightNanos == -1 {
			return Float(float64(leftNanos) / float64(rightNanos)), nil
		}

		if leftNanos%rightNanos == 0 {
			return Int(leftNanos / rightNanos), nil
		}

		return Float(float64(leftNanos) / float64(rightNanos)), nil
	}

	return scaleDuration(leftDuration, right, true)
}

// ModulusChecked rejects duration modulus and preserves legacy numeric behavior.
func ModulusChecked(ctx context.Context, left, right Value) (Value, error) {
	if _, ok := left.(Duration); ok {
		return None, durationBinaryTypeError("%", left, right)
	}
	if _, ok := right.(Duration); ok {
		return None, durationBinaryTypeError("%", left, right)
	}

	return Modulus(ctx, left, right), nil
}

// IncrementChecked rejects incrementing a duration.
func IncrementChecked(ctx context.Context, value Value) (Value, error) {
	if _, ok := value.(Duration); ok {
		return None, Error(ErrInvalidOperation, "increment is not supported for Duration")
	}

	return Increment(ctx, value), nil
}

// DecrementChecked rejects decrementing a duration.
func DecrementChecked(ctx context.Context, value Value) (Value, error) {
	if _, ok := value.(Duration); ok {
		return None, Error(ErrInvalidOperation, "decrement is not supported for Duration")
	}

	return Decrement(ctx, value), nil
}

// NegateChecked applies unary negation while detecting duration overflow.
func NegateChecked(value Value) (Value, error) {
	duration, ok := value.(Duration)
	if !ok {
		switch value := value.(type) {
		case Int:
			return -value, nil
		case Float:
			return -value, nil
		case Boolean:
			return !value, nil
		default:
			return None, nil
		}
	}

	if duration == Duration(math.MinInt64) {
		return None, durationRangeError("negation")
	}

	return -duration, nil
}

// NegativeChecked applies unary minus while detecting duration overflow.
func NegativeChecked(value Value) (Value, error) {
	duration, ok := value.(Duration)
	if ok {
		if duration == Duration(math.MinInt64) {
			return None, durationRangeError("negation")
		}

		return -duration, nil
	}

	switch value := value.(type) {
	case Int:
		return -value, nil
	case Float:
		return -value, nil
	default:
		return None, nil
	}
}

// PositiveChecked applies unary positive to numeric and duration values.
func PositiveChecked(value Value) (Value, error) {
	switch value := value.(type) {
	case Duration:
		return value, nil
	case Int:
		return +value, nil
	case Float:
		return +value, nil
	default:
		return None, nil
	}
}

func scaleDuration(duration Duration, scalar Value, divide bool) (Value, error) {
	switch value := scalar.(type) {
	case Int:
		if divide {
			if value == 0 {
				return None, Error(ErrInvalidOperation, "division by zero")
			}

			if duration == Duration(math.MinInt64) && value == -1 {
				return None, durationRangeError("division")
			}

			return Duration(roundIntegerQuotient(int64(duration), int64(value))), nil
		}

		result, ok := multiplyDurationInt64(int64(duration), int64(value))
		if !ok {
			return None, durationRangeError("multiplication")
		}

		return Duration(result), nil
	case Float:
		scalarFloat := float64(value)

		if math.IsNaN(scalarFloat) || math.IsInf(scalarFloat, 0) {
			return None, Error(ErrInvalidOperation, "duration scale must be finite")
		}

		if divide && scalarFloat == 0 {
			return None, Error(ErrInvalidOperation, "division by zero")
		}

		result := new(big.Rat).SetInt64(int64(duration))
		scalarRat := new(big.Rat).SetFloat64(scalarFloat)

		if divide {
			result.Quo(result, scalarRat)
		} else {
			result.Mul(result, scalarRat)
		}

		nanos, err := roundDurationRat(result)
		if err != nil {
			return None, err
		}

		return Duration(nanos), nil
	default:
		operator := "*"
		if divide {
			operator = "/"
		}

		return None, durationBinaryTypeError(operator, duration, scalar)
	}
}
func roundIntegerQuotient(left, right int64) int64 {
	quotient := left / right
	remainder := left % right

	if remainder == 0 {
		return quotient
	}

	absRemainder := durationMagnitude(remainder)
	absRight := durationMagnitude(right)

	if absRemainder*2 >= absRight {
		if (left < 0) != (right < 0) {
			quotient--
		} else {
			quotient++
		}
	}

	return quotient
}

func durationMagnitude(value int64) uint64 {
	if value >= 0 {
		return uint64(value)
	}

	return uint64(-(value + 1)) + 1
}

func addDurationInt64(left, right int64) (int64, bool) {
	result := left + right
	if (right > 0 && result < left) || (right < 0 && result > left) {
		return 0, false
	}

	return result, true
}

func subtractDurationInt64(left, right int64) (int64, bool) {
	result := left - right
	if (right > 0 && result > left) || (right < 0 && result < left) {
		return 0, false
	}

	return result, true
}

func multiplyDurationInt64(left, right int64) (int64, bool) {
	if left == 0 || right == 0 {
		return 0, true
	}
	if (left == math.MinInt64 && right == -1) || (right == math.MinInt64 && left == -1) {
		return 0, false
	}

	result := left * right
	if result/right != left {
		return 0, false
	}

	return result, true
}

func durationBinaryTypeError(operator string, left, right Value) error {
	return Errorf(
		ErrInvalidOperation,
		"operator %s is not supported for %s and %s",
		operator,
		TypeName(TypeOf(left)),
		TypeName(TypeOf(right)),
	)
}

func durationRangeError(operation string) error {
	return Errorf(ErrRange, "duration %s exceeds the supported range", operation)
}
