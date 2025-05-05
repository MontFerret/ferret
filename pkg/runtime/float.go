package runtime

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"math"
	"strconv"

	"github.com/wI2L/jettison"
)

type Float float64

var (
	NaN = Float(math.NaN())
)

const (
	ZeroFloat = Float(0.0)
)

func NewFloat(input float64) Float {
	return Float(input)
}

func ParseFloat(input interface{}) (Float, error) {
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

	return ZeroFloat, Error(ErrInvalidType, "expected "+"'"+TypeFloat+"'")
}

func MustParseFloat(input interface{}) Float {
	res, err := ParseFloat(input)

	if err != nil {
		panic(err)
	}

	return res
}

func IsNaN(input Float) Boolean {
	return NewBoolean(math.IsNaN(float64(input)))
}

func IsInf(input Float, sign Int) Boolean {
	return NewBoolean(math.IsInf(float64(input), int(sign)))
}

func (f Float) Type() string {
	return TypeFloat
}

func (f Float) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(float64(f), jettison.NoHTMLEscaping())
}

func (f Float) String() string {
	return fmt.Sprintf("%v", float64(f))
}

func (f Float) Unwrap() interface{} {
	return float64(f)
}

func (f Float) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(TypeFloat))

	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, math.Float64bits(float64(f)))
	h.Write(bytes)

	return h.Sum64()
}

func (f Float) Copy() Value {
	return f
}

func (f Float) Compare(other Value) int64 {
	switch otherVal := other.(type) {
	case Float:
		if f == otherVal {
			return 0
		}

		if f < otherVal {
			return -1
		}

		return +1
	case Int:
		f := Float(otherVal)

		if f == f {
			return 0
		}

		if f < f {
			return -1
		}

		return +1
	default:
		return CompareTypes(f, other)
	}
}
