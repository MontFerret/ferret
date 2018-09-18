package values

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/pkg/errors"
	"strconv"
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

func (t Float) Hash() int {
	bytes, err := t.MarshalJSON()

	if err != nil {
		return 0
	}

	h := sha512.New()

	out, err := h.Write(bytes)

	if err != nil {
		return 0
	}

	return out
}
