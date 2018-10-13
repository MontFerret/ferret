package values

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"math"
	"strconv"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/pkg/errors"
)

type Float float64

var ZeroFloat = Float(0.0)

func NewFloat(input float64) Float {
	return Float(input)
}

func ParseFloat(input interface{}) (Float, error) {
	if core.IsNil(input) {
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

	return ZeroFloat, errors.Wrap(core.ErrInvalidType, "expected 'float'")
}

func ParseFloatP(input interface{}) Float {
	res, err := ParseFloat(input)

	if err != nil {
		panic(err)
	}

	return res
}

func IsNaN(input Float) Boolean {
	return NewBoolean(math.IsNaN(float64(input)))
}

func (t Float) MarshalJSON() ([]byte, error) {
	return json.Marshal(float64(t))
}

func (t Float) Type() core.Type {
	return core.FloatType
}

func (t Float) String() string {
	return fmt.Sprintf("%f", t)
}

func (t Float) Compare(other core.Value) int {
	raw := float64(t)

	switch other.Type() {
	case core.FloatType:
		f := other.Unwrap().(float64)

		if raw == f {
			return 0
		}

		if raw < f {
			return -1
		}

		return +1
	case core.IntType:
		i := other.Unwrap().(int)
		f := float64(i)

		if raw == f {
			return 0
		}

		if raw < f {
			return -1
		}

		return +1
	case core.BooleanType, core.NoneType:
		return 1
	default:
		return -1
	}
}

func (t Float) Unwrap() interface{} {
	return float64(t)
}

func (t Float) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(t.Type().String()))
	h.Write([]byte(":"))

	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, math.Float64bits(float64(t)))
	h.Write(bytes)

	return h.Sum64()
}

func (t Float) Copy() core.Value {
	return t
}
