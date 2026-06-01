package mock

import (
	"context"
	"sync/atomic"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type TriggerObservable struct {
	ch                    chan runtime.Message
	dispatchErr           error
	readErr               error
	activeCount           atomic.Int32
	closeCount            atomic.Int32
	dispatchCount         atomic.Int32
	failDispatchRemaining atomic.Int32
	subscribeCount        atomic.Int32
}

func NewTriggerObservable() *TriggerObservable {
	return &TriggerObservable{
		ch: make(chan runtime.Message, 16),
	}
}

func (o *TriggerObservable) Subscribe(_ context.Context, _ runtime.Subscription) (runtime.Stream, error) {
	o.subscribeCount.Add(1)
	o.activeCount.Add(1)

	return &TestStream{
		ch: o.ch,
		onClose: func() {
			o.activeCount.Add(-1)
			o.closeCount.Add(1)
		},
	}, nil
}

func (o *TriggerObservable) Dispatch(_ context.Context, event runtime.DispatchEvent) error {
	o.dispatchCount.Add(1)

	if o.failDispatchRemaining.Load() > 0 {
		o.failDispatchRemaining.Add(-1)
		return o.dispatchErr
	}

	if o.activeCount.Load() <= 0 {
		return nil
	}

	if o.readErr != nil {
		o.ch <- runtime.NewErrorMessage(o.readErr)
		return nil
	}

	o.ch <- runtime.NewValueMessage(runtime.NewObjectWith(map[string]runtime.Value{
		"type": event.Name,
	}))

	return nil
}

func (o *TriggerObservable) FailNextDispatches(n int32, err error) {
	o.failDispatchRemaining.Store(n)
	o.dispatchErr = err
}

func (o *TriggerObservable) FailReadsWith(err error) {
	o.readErr = err
}

func (o *TriggerObservable) SubscribeCount() int32 {
	return o.subscribeCount.Load()
}

func (o *TriggerObservable) DispatchCount() int32 {
	return o.dispatchCount.Load()
}

func (o *TriggerObservable) CloseCount() int32 {
	return o.closeCount.Load()
}

func (o *TriggerObservable) String() string {
	return "trigger_observable"
}

func (o *TriggerObservable) Hash() uint64 {
	return 0
}

func (o *TriggerObservable) Copy() runtime.Value {
	return o
}
