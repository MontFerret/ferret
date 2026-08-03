package runtime_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

var comparisonContractType = runtime.NewType("runtime_test", "ComparisonContract", func(runtime.Value) bool {
	return true
})

var comparisonOtherType = runtime.NewType("runtime_test", "OtherComparisonContract", func(runtime.Value) bool {
	return true
})

type contractHostValue struct {
	typ           runtime.Type
	equalityErr   error
	comparisonErr error
	equalityCalls *int
	compareCalls  *int
	hash          uint64
	unknownType   bool
	equal         bool
	ordering      runtime.Ordering
}

func (v *contractHostValue) String() string {
	panic("comparison dispatch inspected String")
}

func (v *contractHostValue) Hash() uint64 {
	return v.hash
}

func (v *contractHostValue) Copy() runtime.Value {
	panic("comparison dispatch inspected Copy")
}

func (v *contractHostValue) Type() runtime.Type {
	if v.unknownType {
		return nil
	}
	if v.typ == nil {
		return comparisonContractType
	}

	return v.typ
}

func (v *contractHostValue) Unwrap() any {
	panic("comparison dispatch inspected Unwrap")
}

func (v *contractHostValue) MarshalJSON() ([]byte, error) {
	panic("comparison dispatch inspected MarshalJSON")
}

func (v *contractHostValue) Iterate(context.Context) (runtime.Iterator, error) {
	panic("comparison dispatch inspected Iterate")
}

func (v *contractHostValue) Equal(ctx context.Context, _ runtime.Value) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	if v.equalityCalls != nil {
		(*v.equalityCalls)++
	}
	if v.equalityErr != nil {
		return false, v.equalityErr
	}

	return v.equal, nil
}

func (v *contractHostValue) Compare(ctx context.Context, _ runtime.Value) (runtime.Ordering, error) {
	if err := ctx.Err(); err != nil {
		return runtime.Equal, err
	}
	if v.compareCalls != nil {
		(*v.compareCalls)++
	}
	if v.comparisonErr != nil {
		return runtime.Equal, v.comparisonErr
	}

	return v.ordering, nil
}
