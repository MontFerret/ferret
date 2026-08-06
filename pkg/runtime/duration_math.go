package runtime

import (
	"math"
	"math/big"

	"github.com/MontFerret/ferret/v2/pkg/internal/operationerror"
)

func addDurations(left, right Duration) (Value, error) {
	result, ok := addDurationInt64(int64(left), int64(right))
	if !ok {
		return None, durationRangeError("addition")
	}

	return Duration(result), nil
}

func subtractDurations(left, right Duration) (Value, error) {
	result, ok := subtractDurationInt64(int64(left), int64(right))
	if !ok {
		return None, durationRangeError("subtraction")
	}

	return Duration(result), nil
}

func multiplyDuration(duration Duration, scalar Value) (Value, error) {
	switch value := scalar.(type) {
	case Int:
		result, ok := multiplyDurationInt64(int64(duration), int64(value))
		if !ok {
			return None, durationRangeError("multiplication")
		}

		return Duration(result), nil
	case Float:
		return scaleDurationFloat(duration, value, false)
	default:
		return None, Error(ErrInvalidOperation, "invalid Duration multiplier")
	}
}

func divideDuration(duration Duration, scalar Value) (Value, error) {
	switch value := scalar.(type) {
	case Int:
		if value == 0 {
			return None, operationerror.DivisionByZero(ErrInvalidOperation)
		}
		if duration == Duration(math.MinInt64) && value == -1 {
			return None, durationRangeError("division")
		}

		return Duration(int64(duration) / int64(value)), nil
	case Float:
		return scaleDurationFloat(duration, value, true)
	default:
		return None, Error(ErrInvalidOperation, "invalid Duration divisor")
	}
}

func scaleDurationFloat(duration Duration, scalar Float, divide bool) (Value, error) {
	scalarFloat := float64(scalar)
	if math.IsNaN(scalarFloat) || math.IsInf(scalarFloat, 0) {
		return None, Error(ErrInvalidOperation, "duration scale must be finite")
	}
	if divide && scalarFloat == 0 {
		return None, operationerror.DivisionByZero(ErrInvalidOperation)
	}

	result := new(big.Rat).SetInt64(int64(duration))
	scalarRat := new(big.Rat).SetFloat64(scalarFloat)

	if divide {
		result.Quo(result, scalarRat)
	} else {
		result.Mul(result, scalarRat)
	}

	nanos, err := truncateDurationRat(result)
	if err != nil {
		return None, err
	}

	return Duration(nanos), nil
}

func divideDurations(left, right Duration) (Value, error) {
	if right == 0 {
		return None, operationerror.DivisionByZero(ErrInvalidOperation)
	}

	leftNanos := int64(left)
	rightNanos := int64(right)

	if leftNanos == math.MinInt64 && rightNanos == -1 {
		return Float(float64(leftNanos) / float64(rightNanos)), nil
	}
	if leftNanos%rightNanos == 0 {
		return Int(leftNanos / rightNanos), nil
	}

	return Float(float64(leftNanos) / float64(rightNanos)), nil
}

func negateDuration(duration Duration) (Value, error) {
	if duration == Duration(math.MinInt64) {
		return None, durationRangeError("negation")
	}

	return -duration, nil
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

func durationRangeError(operation string) error {
	return Errorf(ErrRange, "Duration %s exceeds the supported range", operation)
}
