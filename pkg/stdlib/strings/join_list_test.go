package strings_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type joinTestList struct {
	runtime.List
	failure error
	cancel  context.CancelFunc
}

func (l *joinTestList) ForEach(ctx context.Context, fn runtime.IndexReadablePredicate) error {
	if err := l.List.ForEach(ctx, fn); err != nil {
		return err
	}

	if l.cancel != nil {
		l.cancel()

		_, err := fn(ctx, runtime.String("second"), 1)

		return err
	}

	return l.failure
}

func (l *joinTestList) New(ctx context.Context) (runtime.List, error) {
	base, err := l.List.New(ctx)
	if err != nil {
		return nil, err
	}

	result := *l
	result.List = base

	return &result, nil
}
