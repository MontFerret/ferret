package debugger

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type fakeSessionServices struct {
	afterRunErr error
	afterCalls  int
	closed      bool
}

func (f *fakeSessionServices) BeforeRun(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (f *fakeSessionServices) AfterRun(_ context.Context, runErr error) error {
	f.afterCalls++
	f.afterRunErr = runErr

	return nil
}

func (f *fakeSessionServices) ExtendContext(ctx context.Context) context.Context {
	return ctx
}

func (f *fakeSessionServices) Materialize(*vm.Result) (*encoding.Output, error) {
	return &encoding.Output{}, nil
}

func (f *fakeSessionServices) Close() error {
	f.closed = true
	return nil
}
