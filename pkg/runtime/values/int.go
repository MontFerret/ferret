package values

import (
	"encoding/json"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/pkg/errors"
	"strconv"
)

type Int int

var ZeroInt = Int(0)

func NewInt(input int) Int {
	return Int(input)
}

func ParseInt(input interface{}) (Int, error) {
	if core.IsNil(input) {
		return ZeroInt, nil
	}

	i, ok := input.(int)

	if ok {
		if i == 0 {
			return ZeroInt, nil
		}

		return Int(i), nil
	}

	// try to cast
	str, ok := input.(string)

	if ok {
		i, err := strconv.Atoi(str)

		if err == nil {
			if i == 0 {
				return ZeroInt, nil
			}

			return Int(i), nil
		}
	}

	return ZeroInt, errors.Wrap(core.ErrInvalidType, "expected 'int'")
}

func ParseIntP(input interface{}) Int {
	res, err := ParseInt(input)

	if err != nil {
		panic(err)
	}

	return res
}

func (t Int) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(t))
}

func (t Int) Type() core.Type {
	return core.IntType
}

func (t Int) String() string {
	return strconv.Itoa(int(t))
}

func (t Int) Compare(other core.Value) int {
	raw := int(t)

	switch other.Type() {
	case core.IntType:
		i := other.Unwrap().(int)

		if raw == i {
			return 0
		}

		if raw < i {
			return -1
		}

		return +1
	case core.FloatType:
		f := other.Unwrap().(float64)
		i := int(f)

		if raw == i {
			return 0
		}

		if raw < i {
			return -1
		}

		return +1
	case core.BooleanType, core.NoneType:
		return 1
	default:
		return -1
	}
}

func (t Int) Unwrap() interface{} {
	return int(t)
}

func (t Int) Hash() int {
	return int(t)
}

func (t Int) Clone() core.Value {
	return t
}
