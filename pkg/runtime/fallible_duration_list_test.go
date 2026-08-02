package runtime_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type fallibleDurationList struct {
	*runtime.Array
	lengthErr error
	atErr     error
}

func newFallibleDurationList(lengthErr, atErr error, values ...runtime.Value) *fallibleDurationList {
	return &fallibleDurationList{
		Array:     runtime.NewArrayWith(values...),
		lengthErr: lengthErr,
		atErr:     atErr,
	}
}

func (l *fallibleDurationList) Length(ctx context.Context) (runtime.Int, error) {
	if l.lengthErr != nil {
		return 0, l.lengthErr
	}

	return l.Array.Length(ctx)
}

func (l *fallibleDurationList) At(ctx context.Context, idx runtime.Int) (runtime.Value, error) {
	if l.atErr != nil {
		return runtime.None, l.atErr
	}

	return l.Array.At(ctx, idx)
}
