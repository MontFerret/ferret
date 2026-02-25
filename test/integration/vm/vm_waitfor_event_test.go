package vm_test

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type testObservable struct {
	events    []runtime.Value
	readCount int32
}

func newTestObservable(events []runtime.Value) *testObservable {
	return &testObservable{
		events: events,
	}
}

func (o *testObservable) Subscribe(ctx context.Context, subscription runtime.Subscription) (runtime.Stream, error) {
	atomic.StoreInt32(&o.readCount, 0)

	ch := make(chan runtime.Message, len(o.events))
	for _, evt := range o.events {
		ch <- &testMessage{value: evt, obs: o}
	}
	close(ch)

	return &testStream{ch: ch}, nil
}

func (o *testObservable) ReadCount() int32 {
	return atomic.LoadInt32(&o.readCount)
}

func (o *testObservable) String() string {
	return "observable"
}

func (o *testObservable) Hash() uint64 {
	return 0
}

func (o *testObservable) Copy() runtime.Value {
	return o
}

type testStream struct {
	ch <-chan runtime.Message
}

func (s *testStream) Read(ctx context.Context) <-chan runtime.Message {
	return s.ch
}

func (s *testStream) Close() error {
	return nil
}

type testMessage struct {
	value runtime.Value
	obs   *testObservable
	err   error
}

func (m *testMessage) Value() runtime.Value {
	if m.obs != nil {
		atomic.AddInt32(&m.obs.readCount, 1)
	}

	return m.value
}

func (m *testMessage) Err() error {
	return m.err
}

func newEventType(eventType string) runtime.Value {
	return runtime.NewObjectWith(map[string]runtime.Value{
		"type": runtime.NewString(eventType),
	})
}

func assertReturnOneAndReads(obs *testObservable, expectedReads int32) func(actual any, expected ...any) string {
	return func(actual any, expected ...any) string {
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
			return fmt.Sprintf("expected return value 1, got %v", actual)
		}

		if reads := obs.ReadCount(); reads != expectedReads {
			return fmt.Sprintf("expected %d reads, got %d", expectedReads, reads)
		}

		return ""
	}
}

func TestWaitforEvent(t *testing.T) {
	matchFirst := newTestObservable([]runtime.Value{
		newEventType("match"),
		newEventType("other"),
	})
	matchSecond := newTestObservable([]runtime.Value{
		newEventType("other"),
		newEventType("match"),
	})

	RunUseCases(t, []UseCase{
		CaseRuntimeError(`LET obj = {}

WAITFOR EVENT "test" IN obj

RETURN NONE`, "Should compile but return an error during execution because the object does not implement the interface"),
		Options(
			CaseFn(`LET obs = @obs
WAITFOR EVENT "test" IN obs FILTER .type == "match"
RETURN 1`, assertReturnOneAndReads(matchFirst, 1)),
			vm.WithParams(map[string]runtime.Value{
				"obs": matchFirst,
			}),
		),
		Options(
			CaseFn(`LET obs = @obs
WAITFOR EVENT "test" IN obs FILTER .type == "match"
RETURN 1`, assertReturnOneAndReads(matchSecond, 2)),
			vm.WithParams(map[string]runtime.Value{
				"obs": matchSecond,
			}),
		),
	})
}
