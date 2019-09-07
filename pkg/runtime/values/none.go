package values

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

type none struct{}

var None = &none{}

func (t *none) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}

func (t *none) Type() core.Type {
	return types.None
}

func (t *none) String() string {
	return ""
}

func (t *none) Compare(other core.Value) int64 {
	if other.Type() == types.None {
		return 0
	}

	return -1
}

func (t *none) Unwrap() interface{} {
	return nil
}

func (t *none) Hash() uint64 {
	return 0
}

func (t *none) Copy() core.Value {
	return None
}

func (t *none) Clone() core.Cloneable {
	return None
}
