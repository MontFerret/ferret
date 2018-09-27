package values

import (
	"crypto/sha512"
	"encoding/json"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/pkg/errors"
	"strings"
)

type Boolean bool

var False = Boolean(false)
var True = Boolean(true)

func NewBoolean(input bool) Boolean {
	return Boolean(input)
}

func ParseBoolean(input interface{}) (Boolean, error) {
	b, ok := input.(bool)

	if ok {
		if b {
			return True, nil
		}

		return False, nil
	}

	s, ok := input.(string)

	if ok {
		s := strings.ToLower(s)

		if s == "true" {
			return True, nil
		}

		if s == "false" {
			return False, nil
		}
	}

	return False, errors.Wrap(core.ErrInvalidType, "expected 'bool'")
}

func ParseBooleanP(input interface{}) Boolean {
	res, err := ParseBoolean(input)

	if err != nil {
		panic(err)
	}

	return res
}

func (t Boolean) MarshalJSON() ([]byte, error) {
	return json.Marshal(bool(t))
}

func (t Boolean) Type() core.Type {
	return core.BooleanType
}

func (t Boolean) String() string {
	if t {
		return "true"
	}

	return "false"
}

func (t Boolean) Compare(other core.Value) int {
	raw := bool(t)

	switch other.Type() {
	case core.BooleanType:
		i := other.Unwrap().(bool)

		if raw == i {
			return 0
		}

		if raw == false && i == true {
			return -1
		}

		return +1
	case core.NoneType:
		return 1
	default:
		return -1
	}
}

func (t Boolean) Unwrap() interface{} {
	return bool(t)
}

func (t Boolean) Hash() int {
	h := sha512.New()

	out, err := h.Write([]byte(t.String()))

	if err != nil {
		return 0
	}

	return out
}

func (t Boolean) Clone() core.Value {
	return t
}
