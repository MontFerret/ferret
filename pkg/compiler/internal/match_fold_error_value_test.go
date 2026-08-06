package internal

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type matchFoldErrorValue struct {
	err    error
	called *int
}

func (v matchFoldErrorValue) String() string {
	return "match-fold-error"
}

func (v matchFoldErrorValue) Hash() uint64 {
	return 31
}

func (v matchFoldErrorValue) Copy() runtime.Value {
	return v
}

func (v matchFoldErrorValue) Type() runtime.Type {
	return runtime.TypeDuration
}

func (v matchFoldErrorValue) Equal(context.Context, runtime.Value) (bool, error) {
	if v.called != nil {
		*v.called++
	}

	return false, v.err
}
