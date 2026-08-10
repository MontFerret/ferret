package vm_test

import (
	"context"
	"sync/atomic"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type panicPolicyObservable struct {
	subscribeErr   error
	subscribePanic any
	readPanic      any
	closePanic     any
	waitForCancel  bool
	closeCount     atomic.Int32
}

func (o *panicPolicyObservable) Subscribe(ctx context.Context, _ runtime.Subscription) (runtime.Stream, error) {
	if o.waitForCancel {
		<-ctx.Done()
		return nil, ctx.Err()
	}

	if o.subscribePanic != nil {
		panic(o.subscribePanic)
	}

	if o.subscribeErr != nil {
		return nil, o.subscribeErr
	}

	return &panicPolicyStream{
		messages:   make(chan runtime.Message),
		readPanic:  o.readPanic,
		closePanic: o.closePanic,
		closeCount: &o.closeCount,
	}, nil
}

func (o *panicPolicyObservable) String() string {
	return "panic_policy_observable"
}

func (o *panicPolicyObservable) Hash() uint64 {
	return 0
}

func (o *panicPolicyObservable) Copy() runtime.Value {
	return o
}
