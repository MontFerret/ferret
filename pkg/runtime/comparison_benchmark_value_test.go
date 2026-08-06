package runtime_test

import (
	"context"
	"strconv"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

var comparisonBenchmarkType = runtime.NewType(
	"github.com/MontFerret/ferret/v2/pkg/runtime_test",
	"comparisonBenchmarkValue",
	func(value runtime.Value) bool {
		_, ok := value.(comparisonBenchmarkValue)
		return ok
	},
)

type comparisonBenchmarkValue struct {
	value int
}

func (v comparisonBenchmarkValue) Type() runtime.Type {
	return comparisonBenchmarkType
}

func (v comparisonBenchmarkValue) String() string {
	return strconv.Itoa(v.value)
}

func (v comparisonBenchmarkValue) Hash() uint64 {
	return uint64(v.value)
}

func (v comparisonBenchmarkValue) Copy() runtime.Value {
	return v
}

func (v comparisonBenchmarkValue) Equal(ctx context.Context, other runtime.Value) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}

	otherValue, ok := other.(comparisonBenchmarkValue)
	return ok && v.value == otherValue.value, nil
}

func (v comparisonBenchmarkValue) Compare(
	ctx context.Context,
	other runtime.Value,
) (runtime.Ordering, error) {
	if err := ctx.Err(); err != nil {
		return runtime.Equal, err
	}

	otherValue, ok := other.(comparisonBenchmarkValue)
	if !ok {
		return runtime.Equal, runtime.Error(runtime.ErrInvalidOperation, "incompatible benchmark value")
	}

	switch {
	case v.value < otherValue.value:
		return runtime.Less, nil
	case v.value > otherValue.value:
		return runtime.Greater, nil
	default:
		return runtime.Equal, nil
	}
}
