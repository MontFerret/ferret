package runtime

import (
	"math/big"
	"strings"
	"time"
)

// ParseDuration parses a Ferret duration literal or a normalized Go duration.
func ParseDuration(input string) (Duration, error) {
	raw := strings.TrimSpace(input)
	if raw == "" {
		return ZeroDuration, Error(ErrInvalidArgument, "empty duration")
	}

	number, multiplier, ok := splitDuration(raw)
	if ok {
		value, parsed := new(big.Rat).SetString(number)
		if parsed {
			value.Mul(value, new(big.Rat).SetInt64(multiplier))
			nanos, err := roundDurationRat(value)
			if err != nil {
				return ZeroDuration, err
			}

			return Duration(nanos), nil
		}
	}

	parsed, err := time.ParseDuration(raw)
	if err != nil {
		return ZeroDuration, err
	}

	return Duration(parsed), nil
}

func splitDuration(input string) (string, int64, bool) {
	lower := strings.ToLower(input)

	switch {
	case strings.HasSuffix(lower, "ms"):
		return input[:len(input)-2], int64(time.Millisecond), true
	case strings.HasSuffix(lower, "s"):
		return input[:len(input)-1], int64(time.Second), true
	case strings.HasSuffix(lower, "m"):
		return input[:len(input)-1], int64(time.Minute), true
	case strings.HasSuffix(lower, "h"):
		return input[:len(input)-1], int64(time.Hour), true
	case strings.HasSuffix(lower, "d"):
		return input[:len(input)-1], int64(24 * time.Hour), true
	default:
		return "", 0, false
	}
}

func roundDurationRat(value *big.Rat) (int64, error) {
	quotient := new(big.Int)
	remainder := new(big.Int)
	quotient.QuoRem(value.Num(), value.Denom(), remainder)

	if remainder.Sign() != 0 {
		doubled := new(big.Int).Lsh(new(big.Int).Abs(remainder), 1)
		if doubled.Cmp(value.Denom()) >= 0 {
			if value.Sign() < 0 {
				quotient.Sub(quotient, big.NewInt(1))
			} else {
				quotient.Add(quotient, big.NewInt(1))
			}
		}
	}

	if !quotient.IsInt64() {
		return 0, Error(ErrRange, "duration exceeds the supported range")
	}

	return quotient.Int64(), nil
}
