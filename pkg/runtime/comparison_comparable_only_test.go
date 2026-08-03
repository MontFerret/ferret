package runtime_test

import (
	"context"
	"strconv"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type comparableOnlyValue int

func (v comparableOnlyValue) String() string {
	return strconv.Itoa(int(v))
}

func (v comparableOnlyValue) Hash() uint64 {
	return uint64(v)
}

func (v comparableOnlyValue) Copy() runtime.Value {
	return v
}

func (v comparableOnlyValue) Compare(
	ctx context.Context,
	other runtime.Value,
) (runtime.Ordering, error) {
	if err := ctx.Err(); err != nil {
		return runtime.Equal, err
	}

	otherValue, ok := other.(comparableOnlyValue)
	if !ok {
		return runtime.Equal, runtime.Error(runtime.ErrInvalidOperation, "incompatible ordering-only value")
	}

	switch {
	case v < otherValue:
		return runtime.Less, nil
	case v > otherValue:
		return runtime.Greater, nil
	default:
		return runtime.Equal, nil
	}
}
