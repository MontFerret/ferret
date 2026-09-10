package valueset_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	collisionValue struct {
		err   error
		label string
	}

	observedCollisionValue struct {
		err   error
		calls *[]string
		label string
	}
)

func (v collisionValue) String() string {
	return v.label
}

func (v collisionValue) Hash() uint64 {
	return 7
}

func (v collisionValue) Copy() runtime.Value {
	return v
}

func (v collisionValue) Equal(ctx context.Context, other runtime.Value) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}

	if v.err != nil {
		return false, v.err
	}

	o, ok := other.(collisionValue)
	if !ok {
		return false, nil
	}

	return v.label == o.label, nil
}

func (v *observedCollisionValue) String() string {
	return v.label
}

func (v *observedCollisionValue) Hash() uint64 {
	return 7
}

func (v *observedCollisionValue) Copy() runtime.Value {
	return v
}

func (v *observedCollisionValue) Equal(ctx context.Context, other runtime.Value) (bool, error) {
	*v.calls = append(*v.calls, v.label)
	if err := ctx.Err(); err != nil {
		return false, err
	}

	o, ok := other.(*observedCollisionValue)

	return ok && v.label == o.label, v.err
}
