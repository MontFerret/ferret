package values

import (
	"hash/fnv"
	"strings"

	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

type Boolean bool

const (
	False = Boolean(false)
	True  = Boolean(true)
)

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
		return Boolean(strings.ToLower(s) == "true"), nil
	}

	return False, core.Error(core.ErrInvalidType, "expected 'bool'")
}

func MustParseBoolean(input interface{}) Boolean {
	res, err := ParseBoolean(input)

	if err != nil {
		panic(err)
	}

	return res
}

func (t Boolean) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(bool(t), jettison.NoHTMLEscaping())
}

func (t Boolean) Type() core.Type {
	return types.Boolean
}

func (t Boolean) String() string {
	if t {
		return "true"
	}

	return "false"
}

func (t Boolean) Compare(other core.Value) int64 {
	raw := bool(t)

	if types.Boolean.Equals(other.Type()) {
		i := other.Unwrap().(bool)

		if raw == i {
			return 0
		}

		if !raw && i {
			return -1
		}

		return +1
	}

	return types.Compare(types.Boolean, other.Type())
}

func (t Boolean) Unwrap() interface{} {
	return bool(t)
}

func (t Boolean) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(t.Type().String()))
	h.Write([]byte(":"))
	h.Write([]byte(t.String()))

	return h.Sum64()
}

func (t Boolean) Copy() core.Value {
	return t
}
