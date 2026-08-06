package runtime

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"math"
	"strconv"
)

type Float float64

const (
	ZeroFloat = Float(0.0)
)

func NaN() Float {
	return Float(math.NaN())
}

func NewFloat(input float64) Float {
	return Float(input)
}

func ParseFloat(input any) (Float, error) {
	if IsNil(input) {
		return ZeroFloat, nil
	}

	i, ok := input.(float64)

	if ok {
		if i == 0 {
			return ZeroFloat, nil
		}

		return Float(i), nil
	}

	// try to cast
	str, ok := input.(string)

	if ok {
		i, err := strconv.Atoi(str)

		if err == nil {
			if i == 0 {
				return ZeroFloat, nil
			}

			return Float(i), nil
		}
	}

	return ZeroFloat, Error(ErrInvalidType, fmt.Sprintf("expected %s", TypeFloat))
}

func MustParseFloat(input any) Float {
	res, err := ParseFloat(input)

	if err != nil {
		panic(err)
	}

	return res
}

func ToFloat(ctx context.Context, input Value) (Float, error) {
	switch val := input.(type) {
	case Float:
		return val, nil
	case Int:
		return Float(val), nil
	case String:
		i, err := strconv.ParseFloat(string(val), 64)

		if err != nil {
			return ZeroFloat, newConversionError(TypeFloat, err)
		}

		return Float(i), nil
	case Boolean:
		if val {
			return Float(1), nil
		}

		return Float(0), nil
	case DateTime:
		dt := input.(DateTime)

		if dt.IsZero() {
			return ZeroFloat, nil
		}

		return NewFloat(float64(dt.Unix())), nil
	case List:
		iterator, err := val.Iterate(ctx)

		if err != nil {
			return ZeroFloat, err
		}

		res := ZeroFloat

		for {
			val, _, err := iterator.Next(ctx)
			if errors.Is(err, io.EOF) {
				break
			}

			if errors.Is(err, ErrTimeout) {
				break
			}

			if err != nil {
				continue
			}

			f, err := ToFloat(ctx, val)

			if err != nil {
				continue
			}

			res += f
		}

		return res, nil
	default:
		return ZeroFloat, newConversionError(TypeFloat, TypeErrorOf(input, TypeFloat))
	}
}

func IsNaN(input Float) Boolean {
	return NewBoolean(math.IsNaN(float64(input)))
}

func IsInf(input Float, sign Int) Boolean {
	return NewBoolean(math.IsInf(float64(input), int(sign)))
}

func (f Float) Type() Type {
	return TypeFloat
}

func (f Float) String() string {
	return fmt.Sprintf("%v", float64(f))
}

func (f Float) Hash() uint64 {
	if integer, ok := exactIntFromFloat(f); ok {
		return integer.Hash()
	}

	h := fnv.New64a()

	h.Write([]byte(TypeFloat.Name()))

	bytes := make([]byte, 8)
	bits := math.Float64bits(float64(f))
	if math.IsNaN(float64(f)) {
		bits = math.Float64bits(math.NaN())
	}

	binary.LittleEndian.PutUint64(bytes, bits)
	h.Write(bytes)

	return h.Sum64()
}

func (f Float) Copy() Value {
	return f
}

func (f Float) Unwrap() any {
	return float64(f)
}
