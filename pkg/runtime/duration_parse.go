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
			nanos, err := truncateDurationRat(value)
			if err != nil {
				return ZeroDuration, err
			}

			return Duration(nanos), nil
		}
	}

	parsed, ok, err := parseCompoundDuration(raw)
	if !ok {
		return ZeroDuration, Errorf(ErrInvalidArgument, "invalid duration %q", raw)
	}

	if err != nil {
		return ZeroDuration, err
	}

	return parsed, nil
}

func parseCompoundDuration(input string) (Duration, bool, error) {
	sign := int64(1)
	if len(input) > 0 {
		switch input[0] {
		case '-':
			sign = -1
			input = input[1:]
		case '+':
			input = input[1:]
		}
	}

	if input == "0" {
		return ZeroDuration, true, nil
	}
	if input == "" {
		return ZeroDuration, false, nil
	}

	total := new(big.Rat)
	for input != "" {
		numberEnd := durationNumberEnd(input)
		if numberEnd == 0 {
			return ZeroDuration, false, nil
		}

		number, ok := new(big.Rat).SetString(input[:numberEnd])
		if !ok || number.Sign() < 0 {
			return ZeroDuration, false, nil
		}

		multiplier, unitLength, ok := durationUnit(input[numberEnd:])
		if !ok {
			return ZeroDuration, false, nil
		}

		number.Mul(number, new(big.Rat).SetInt64(multiplier))
		total.Add(total, number)
		input = input[numberEnd+unitLength:]
	}

	if sign < 0 {
		total.Neg(total)
	}

	nanos, err := truncateDurationRat(total)
	if err != nil {
		return ZeroDuration, true, err
	}

	return Duration(nanos), true, nil
}

func durationNumberEnd(input string) int {
	digits := 0
	dots := 0

	for index := 0; index < len(input); index++ {
		switch input[index] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			digits++
		case '.':
			dots++
			if dots > 1 {
				return 0
			}
		default:
			if digits == 0 {
				return 0
			}
			return index
		}
	}

	return 0
}

func durationUnit(input string) (int64, int, bool) {
	lower := strings.ToLower(input)

	switch {
	case strings.HasPrefix(lower, "ms"):
		return int64(time.Millisecond), len("ms"), true
	case strings.HasPrefix(lower, "us"):
		return int64(time.Microsecond), len("us"), true
	case strings.HasPrefix(lower, "µs"):
		return int64(time.Microsecond), len("µs"), true
	case strings.HasPrefix(lower, "μs"):
		return int64(time.Microsecond), len("μs"), true
	case strings.HasPrefix(lower, "ns"):
		return int64(time.Nanosecond), len("ns"), true
	case strings.HasPrefix(lower, "d"):
		return int64(24 * time.Hour), len("d"), true
	case strings.HasPrefix(lower, "h"):
		return int64(time.Hour), len("h"), true
	case strings.HasPrefix(lower, "m"):
		return int64(time.Minute), len("m"), true
	case strings.HasPrefix(lower, "s"):
		return int64(time.Second), len("s"), true
	default:
		return 0, 0, false
	}
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

func truncateDurationRat(value *big.Rat) (int64, error) {
	quotient := new(big.Int)
	quotient.Quo(value.Num(), value.Denom())

	if !quotient.IsInt64() {
		return 0, Error(ErrRange, "duration exceeds the supported range")
	}

	return quotient.Int64(), nil
}
