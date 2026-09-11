package runtime_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type typedList struct {
	*runtime.Array
}

func (t typedList) Concat(ctx context.Context, other runtime.List) error {
	return t.Array.Concat(ctx, other)
}

func (t typedList) Type() runtime.Type {
	return typeTypedList
}

func (t typedList) New(ctx context.Context) (runtime.List, error) {
	base, err := t.Array.New(ctx)
	if err != nil {
		return nil, err
	}

	result := t
	result.Array = base.(*runtime.Array)

	return result, nil
}

var typeTypedList = runtime.NewType("test", "typedList", func(value runtime.Value) bool {
	_, ok := value.(typedList)

	return ok
})
