package runtime_test

import (
	"context"
	"strconv"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type equatableOnlyValue int

func (v equatableOnlyValue) String() string {
	return strconv.Itoa(int(v))
}

func (v equatableOnlyValue) Hash() uint64 {
	return uint64(v)
}

func (v equatableOnlyValue) Copy() runtime.Value {
	return v
}

func (v equatableOnlyValue) Equal(ctx context.Context, other runtime.Value) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}

	otherValue, ok := other.(equatableOnlyValue)
	return ok && v == otherValue, nil
}
