package mock

import (
	"context"
	"sync/atomic"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type Observable struct {
	events    []runtime.Value
	readCount int32
}

func NewObservable(events []runtime.Value) *Observable {
	return &Observable{
		events: events,
	}
}

func (o *Observable) Subscribe(_ context.Context, _ runtime.Subscription) (runtime.Stream, error) {
	atomic.StoreInt32(&o.readCount, 0)

	ch := make(chan runtime.Message, len(o.events))
	for _, evt := range o.events {
		ch <- &Message{value: evt, obs: o}
	}
	close(ch)

	return &TestStream{ch: ch}, nil
}

func (o *Observable) ReadCount() int32 {
	return atomic.LoadInt32(&o.readCount)
}

func (o *Observable) String() string {
	return "observable"
}

func (o *Observable) Hash() uint64 {
	return 0
}

func (o *Observable) Copy() runtime.Value {
	return o
}

type TestStream struct {
	ch <-chan runtime.Message
}

func (s *TestStream) Read(ctx context.Context) <-chan runtime.Message {
	return s.ch
}

func (s *TestStream) Close() error {
	return nil
}

func NewTestEventType(eventType string) runtime.Value {
	return runtime.NewObjectWith(map[string]runtime.Value{
		"type": runtime.NewString(eventType),
	})
}
