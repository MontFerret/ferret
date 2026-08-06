package runtime

import (
	"context"
)

type none struct{}

var (
	None = &none{}
)

func (n *none) Type() Type {
	return TypeNone
}

func (n *none) String() string {
	return ""
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
