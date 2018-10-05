package values

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type none struct{}

var None = &none{}

func (t *none) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}

func (t *none) Type() core.Type {
	return core.NoneType
}

func (t *none) String() string {
	return ""
}

func (t *none) Compare(other core.Value) int {
	switch other.Type() {
	case core.NoneType:
		return 0
	default:
		return -1
	}
}

func (t *none) Unwrap() interface{} {
	return nil
}

func (t *none) Hash() uint64 {
	return 0
}

func (t *none) Clone() core.Value {
	return None
}
