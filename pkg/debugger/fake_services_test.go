package debugger

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type fakeSessionServices struct {
	beforeRun   func(context.Context) (context.Context, error)
	afterRun    func(context.Context, error) error
	afterRunErr error
	closeErr    error
	afterCalls  int
	closeCalls  int
	closed      bool
}

func (f *fakeSessionServices) BeforeRun(ctx context.Context) (context.Context, error) {
	if f.beforeRun != nil {
		return f.beforeRun(ctx)
	}

	return ctx, nil
}

func (f *fakeSessionServices) AfterRun(ctx context.Context, runErr error) error {
	f.afterCalls++
	f.afterRunErr = runErr
	if f.afterRun != nil {
		return f.afterRun(ctx, runErr)
	}

	return nil
}

func (f *fakeSessionServices) ExtendContext(ctx context.Context) context.Context {
	return ctx
}

func (f *fakeSessionServices) Materialize(*vm.Result) (*encoding.Output, error) {
	return &encoding.Output{}, nil
}

func (f *fakeSessionServices) Close() error {
	f.closeCalls++
	f.closed = true

	return f.closeErr
}
