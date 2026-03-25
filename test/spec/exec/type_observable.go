package exec

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/test/spec/assert"
)

type TestObservable struct {
	events    []runtime.Value
	readCount int32
}

func NewTestObservable(events []runtime.Value) *TestObservable {
	return &TestObservable{
		events: events,
	}
}

func (o *TestObservable) Subscribe(ctx context.Context, subscription runtime.Subscription) (runtime.Stream, error) {
	atomic.StoreInt32(&o.readCount, 0)

	ch := make(chan runtime.Message, len(o.events))
	for _, evt := range o.events {
		ch <- &TestMessage{value: evt, obs: o}
	}
	close(ch)

	return &TestStream{ch: ch}, nil
}

func (o *TestObservable) ReadCount() int32 {
	return atomic.LoadInt32(&o.readCount)
}

func (o *TestObservable) String() string {
	return "observable"
}

func (o *TestObservable) Hash() uint64 {
	return 0
}

func (o *TestObservable) Copy() runtime.Value {
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

type TestMessage struct {
	value runtime.Value
	obs   *TestObservable
	err   error
}

func (m *TestMessage) Value() runtime.Value {
	if m.obs != nil {
		atomic.AddInt32(&m.obs.readCount, 1)
	}

	return m.value
}

func (m *TestMessage) Err() error {
	return m.err
}

func NewTestEventType(eventType string) runtime.Value {
	return runtime.NewObjectWith(map[string]runtime.Value{
		"type": runtime.NewString(eventType),
	})
}

func ObservableReturnOneAndReads(obs *TestObservable, expectedReads int32) assert.Unary {
	return func(actual any) error {
		var ok bool
		switch v := actual.(type) {
		case float64:
			ok = v == 1
		case int:
			ok = v == 1
		case int64:
			ok = v == 1
		}

		if !ok {
			return fmt.Errorf("expected return value 1, got %v", actual)
		}

		if reads := obs.ReadCount(); reads != expectedReads {
			return fmt.Errorf("expected %d reads, got %d", expectedReads, reads)
		}

		return nil
	}
}
