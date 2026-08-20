package runtime_test

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type exampleAddableValue struct {
	label string
}

func (v exampleAddableValue) String() string {
	return v.label
}

func (v exampleAddableValue) Hash() uint64 {
	return 1
}

func (v exampleAddableValue) Copy() runtime.Value {
	return v
}

func (v exampleAddableValue) Add(ctx context.Context, right runtime.Value) (runtime.Value, error) {
	if err := ctx.Err(); err != nil {
		return runtime.None, err
	}

	integer, ok := right.(runtime.Int)
	if !ok {
		return runtime.None, runtime.ErrUnsupportedOperands
	}

	return runtime.NewString(fmt.Sprintf("%s+%d", v.label, integer)), nil
}

func (v exampleAddableValue) RightAdd(ctx context.Context, left runtime.Value) (runtime.Value, error) {
	if err := ctx.Err(); err != nil {
		return runtime.None, err
	}

	integer, ok := left.(runtime.Int)
	if !ok {
		return runtime.None, runtime.ErrUnsupportedOperands
	}

	return runtime.NewString(fmt.Sprintf("%d+%s", integer, v.label)), nil
}
