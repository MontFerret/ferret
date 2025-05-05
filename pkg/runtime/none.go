package runtime

import (
	"context"
)

type none struct{}

var (
	None = &none{}
)

func (n *none) Type() string {
	return TypeNone
}

func (n *none) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}

func (n *none) String() string {
	return ""
}

func (n *none) Unwrap() interface{} {
	return nil
}

func (n *none) Hash() uint64 {
	return 0
}

func (n *none) Copy() Value {
	return None
}

func (n *none) Clone(_ context.Context) (Cloneable, error) {
	return None, nil
}

func (n *none) Compare(_ context.Context, other Value) (int64, error) {
	if n == other {
		return 0, nil
	}

	return -1, nil
}
