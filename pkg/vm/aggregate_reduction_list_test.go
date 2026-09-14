package vm

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// The nil List exposes accidental copying, indexed access, or other capabilities.
type reductionList struct {
	runtime.List
	err    error
	cancel context.CancelFunc
	values []runtime.Value
	closed int
}

func (l *reductionList) ForEach(ctx context.Context, visit runtime.IndexReadablePredicate) error {
	for index, value := range l.values {
		more, err := visit(ctx, value, runtime.Int(index))
		if err != nil {
			return err
		}

		if !more {
			return nil
		}
	}

	if l.cancel != nil {
		l.cancel()
	}

	return l.err
}

func (l *reductionList) Close() error {
	l.closed++

	return nil
}
