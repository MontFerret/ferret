package vm

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/test/spec/mock"
)

type streamGroupSetupObservable struct {
	fail        error
	started     *atomic.Int32
	gate        chan struct{}
	gateOnce    *sync.Once
	subscribeFn func(context.Context) (runtime.Stream, error)
	expected    int32
}

func newStreamGroupSetupObservable(
	started *atomic.Int32,
	gate chan struct{},
	gateOnce *sync.Once,
	expected int32,
	fail error,
) *streamGroupSetupObservable {
	return &streamGroupSetupObservable{
		started:  started,
		gate:     gate,
		gateOnce: gateOnce,
		expected: expected,
		fail:     fail,
	}
}

func (o *streamGroupSetupObservable) Subscribe(ctx context.Context, _ runtime.Subscription) (runtime.Stream, error) {
	if o.subscribeFn != nil {
		return o.subscribeFn(ctx)
	}

	if o.started.Add(1) == o.expected {
		o.gateOnce.Do(func() { close(o.gate) })
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-o.gate:
	}

	if o.fail != nil {
		return nil, o.fail
	}

	return &mock.TestStream{}, nil
}

func (o *streamGroupSetupObservable) String() string {
	return "stream_group_setup_observable"
}

func (o *streamGroupSetupObservable) Hash() uint64 {
	return 0
}

func (o *streamGroupSetupObservable) Copy() runtime.Value {
	return o
}
