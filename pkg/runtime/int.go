package runtime

import (
	"context"
	"encoding/binary"
	"errors"
	"hash/fnv"
	"io"
	"math"
	"strconv"
)

type Int int64

const (
	ZeroInt = Int(0)
	MaxInt  = Int(math.MaxInt64)
)

func NewInt(input int) Int {
	return Int(int64(input))
}

func NewInt64(input int64) Int {
	return Int(input)
}

func ParseInt(input any) (Int, error) {
	if IsNil(input) {
		return ZeroInt, nil
	}

	switch val := input.(type) {
	case int:
		return Int(val), nil
	case int64:
		return Int(val), nil
	case int32:
		return Int(val), nil
	case int16:
		return Int(val), nil
	case int8:
		return Int(val), nil
	case string:
		i, err := strconv.Atoi(val)

		if err == nil {
			if i == 0 {
				return ZeroInt, nil
			}

			return Int(i), nil
		}

		return ZeroInt, err
	default:
		return ZeroInt, Error(ErrInvalidType, "expected 'int'")
	}
}

func ToInt(ctx context.Context, input Value) (Int, error) {
	switch val := input.(type) {
	case Int:
		return val, nil
	case Float:
		return Int(val), nil
	case String:
		i, err := strconv.ParseInt(string(val), 10, 64)

		if err != nil {
			return ZeroInt, newConversionError(TypeInt, err)
		}

		return Int(i), nil
	case Boolean:
		if val {
			return Int(1), nil
		}

		return Int(0), nil
	case DateTime:
		dt := input.(DateTime)

		if dt.IsZero() {
			return ZeroInt, nil
		}

		return NewInt(int(dt.Unix())), nil
	case List:
		iterator, err := val.Iterate(ctx)

		if err != nil {
			return ZeroInt, err
		}

		res := ZeroInt

		for {
			item, _, err := iterator.Next(ctx)
			if errors.Is(err, io.EOF) {
				break
			}

			if errors.Is(err, ErrTimeout) {
				break
			}

			if err != nil {
				continue
			}

			i, err := ToInt(ctx, item)

			if err != nil {
				continue
			}

			res += i
		}

		return res, nil
	default:
		return ZeroInt, newConversionError(TypeInt, TypeErrorOf(input, TypeInt))
	}
}

func ToIntSafe(ctx context.Context, input Value) Int {
	result, err := ToInt(ctx, input)

	if err != nil {
		return ZeroInt
	}

	if result > 0 {
		return result
	}

	return ZeroInt
}

// ToIntDefault attempts to convert an arbitrary Value into an Int.
// If the conversion fails or if the resulting Int is not greater than zero, it returns the provided defaultValue.
// This function is useful for safely converting values to Int while providing a fallback option in case of errors or non-positive results.
func ToIntDefault(ctx context.Context, input Value, defaultValue Int) (Int, error) {
	result, err := ToInt(ctx, input)

	if err != nil {
		return defaultValue, err
	}

	if result > 0 {
		return result, nil
	}

	return defaultValue, nil
}

func (i Int) Type() Type {
	return TypeInt
}

func (i Int) String() string {
	return strconv.FormatInt(int64(i), 10)
}

func (i Int) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(TypeInt.Name()))

	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(i))
	h.Write(bytes)

	return h.Sum64()
}

func (i Int) Copy() Value {
	return i
}

func (i Int) Unwrap() any {
	return int(i)
}
