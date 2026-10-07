package engine

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestOutputPresenceTransferAndOutcomes(t *testing.T) {
	operationErr, cleanupErr := errors.New("execution"), errors.New("cleanup")
	for _, tc := range []struct {
		operationErr error
		cleanupErr   error
		content      *encoding.Content
		name         string
	}{
		{name: "absent"},
		{name: "absent failure", operationErr: operationErr},
		{name: "present nil bytes", content: &encoding.Content{}},
		{name: "present empty bytes", content: &encoding.Content{Data: []byte{}}},
		{name: "populated", content: &encoding.Content{Metadata: encoding.Metadata{ContentType: "custom", Length: 2, LengthKnown: true}, Data: []byte("42")}},
		{name: "content and errors", content: &encoding.Content{Data: []byte("4")}, operationErr: operationErr, cleanupErr: cleanupErr},
	} {
		t.Run(tc.name, func(t *testing.T) {
			o := newOutput(t.Context(), runOutcome{content: tc.content, operationErr: tc.operationErr, cleanupErr: tc.cleanupErr})
			metadata := o.Metadata()
			content, err := o.Collect(t.Context())
			if content != tc.content {
				t.Fatal("collection did not transfer the owned content")
			}

			for _, cause := range []error{tc.operationErr, tc.cleanupErr} {
				if cause != nil && !errors.Is(err, cause) {
					t.Fatalf("lost cause %v: %v", cause, err)
				}
			}

			if tc.operationErr == nil && tc.cleanupErr == nil && err != nil {
				t.Fatal(err)
			}

			if o.content != nil || o.cancel != nil || o.operationErr != nil {
				t.Fatal("finalized output retained consumption resources")
			}

			if content != nil {
				content.Metadata.ContentType = "caller mutation"
			}

			if o.Metadata() != metadata {
				t.Fatal("descriptor changed after transfer")
			}

			for range 3 {
				if err := o.Close(); !errors.Is(err, tc.cleanupErr) || (tc.cleanupErr == nil && err != nil) {
					t.Fatalf("close=%v", err)
				}
			}

			if _, err := o.Collect(t.Context()); !errors.Is(err, encoding.ErrOutputClosed) {
				t.Fatalf("repeated collect=%v", err)
			}

			if err := o.Consume(t.Context(), func(context.Context, []byte) error { t.Fatal("late callback"); return nil }); !errors.Is(err, encoding.ErrOutputClosed) {
				t.Fatal(err)
			}
		})
	}
}

func TestOutputConsumeBorrowingPresenceAndConsumerFailure(t *testing.T) {
	failure := errors.New("consumer failed")
	for _, content := range []*encoding.Content{nil, {}, {Data: []byte{}}, {Data: []byte("42")}} {
		o := newOutput(t.Context(), runOutcome{content: content})
		calls := 0
		err := o.Consume(t.Context(), func(ctx context.Context, data []byte) error {
			calls++
			if len(data) != len(content.Data) {
				t.Fatal("payload changed")
			}

			if len(data) > 0 && &data[0] != &content.Data[0] {
				t.Fatal("borrowed delivery copied the payload")
			}

			_ = o.Metadata()
			if _, err := o.Collect(ctx); !errors.Is(err, encoding.ErrOutputInUse) {
				t.Fatalf("reentrant collect=%v", err)
			}

			return failure
		})
		if content == nil {
			if calls != 0 || err != nil {
				t.Fatalf("absent calls=%d err=%v", calls, err)
			}
		} else if calls != 1 || !errors.Is(err, failure) {
			t.Fatalf("present calls=%d err=%v", calls, err)
		}

		if err := o.Close(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestOutputValidationOrderingAndRetry(t *testing.T) {
	o := newOutput(t.Context(), runOutcome{content: &encoding.Content{Data: []byte("42")}})
	consumer := func(context.Context, []byte) error { return nil }
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	for _, call := range []func() error{
		func() error { return o.Consume(nil, nil) },
		func() error { return o.Consume(ctx, nil) },
		func() error { _, err := o.Collect(nil); return err },
	} {
		if err := call(); !errors.Is(err, runtime.ErrInvalidArgument) {
			t.Fatalf("invalid call=%v", err)
		}
	}

	if err := o.Consume(ctx, consumer); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}

	if _, err := o.Collect(t.Context()); err != nil {
		t.Fatalf("rejected calls claimed output: %v", err)
	}

	if err := o.Consume(nil, consumer); !errors.Is(err, runtime.ErrInvalidArgument) {
		t.Fatal(err)
	}

	if err := o.Consume(ctx, consumer); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}

	if err := o.Consume(t.Context(), consumer); !errors.Is(err, encoding.ErrOutputClosed) {
		t.Fatal(err)
	}
}

func TestOutputCloseCancelsAndWaitsForBorrowedCallback(t *testing.T) {
	cleanupErr := errors.New("early cleanup")
	o := newOutput(t.Context(), runOutcome{content: &encoding.Content{Data: []byte("borrowed")}, cleanupErr: cleanupErr})
	entered, canceled, release := make(chan struct{}), make(chan struct{}), make(chan struct{})
	consumed, closed := make(chan error, 1), make(chan error, 4)
	go func() {
		consumed <- o.Consume(t.Context(), func(ctx context.Context, chunk []byte) error {
			close(entered)
			<-ctx.Done()
			close(canceled)
			<-release
			if string(chunk) != "borrowed" {
				t.Error("borrowed bytes invalidated before callback settlement")
			}

			return nil
		})
	}()
	<-entered
	if _, err := o.Collect(t.Context()); !errors.Is(err, encoding.ErrOutputInUse) {
		t.Fatal(err)
	}

	if err := o.Consume(nil, nil); !errors.Is(err, runtime.ErrInvalidArgument) {
		t.Fatal(err)
	}

	rejected, reject := context.WithCancel(t.Context())
	reject()
	if err := o.Consume(rejected, func(context.Context, []byte) error { t.Fatal("rejected callback"); return nil }); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}

	for range 4 {
		go func() { closed <- o.Close() }()
	}

	<-canceled
	select {
	case err := <-closed:
		t.Fatalf("close returned with callback active: %v", err)
	default:
	}

	if _, err := o.Collect(t.Context()); !errors.Is(err, encoding.ErrOutputClosed) {
		t.Fatal(err)
	}

	close(release)
	if err := <-consumed; !errors.Is(err, encoding.ErrOutputClosed) || !errors.Is(err, cleanupErr) {
		t.Fatalf("consumption=%v", err)
	}

	for range 4 {
		if err := <-closed; !errors.Is(err, cleanupErr) || errors.Is(err, encoding.ErrOutputClosed) {
			t.Fatalf("close=%v", err)
		}
	}
}

func TestOutputCancellationContexts(t *testing.T) {
	type contextKey struct{}
	for _, which := range []string{"invocation", "consumption"} {
		t.Run(which, func(t *testing.T) {
			cause := errors.New(which)
			invocation, cancelInvocation := context.WithCancelCause(t.Context())
			defer cancelInvocation(nil)
			shortDeadline := time.Now().Add(time.Hour)
			invocation, cancelDeadline := context.WithDeadline(invocation, shortDeadline)
			defer cancelDeadline()
			consumption, cancelConsumption := context.WithCancelCause(context.WithValue(t.Context(), contextKey{}, "consumer"))
			defer cancelConsumption(nil)
			o := newOutput(invocation, runOutcome{content: &encoding.Content{Data: []byte("42")}})
			err := o.Consume(consumption, func(ctx context.Context, _ []byte) error {
				if deadline, ok := ctx.Deadline(); !ok || deadline != shortDeadline {
					t.Fatal("effective deadline did not bound both contexts")
				}

				if ctx.Value(contextKey{}) != "consumer" {
					t.Fatal("consumption context values lost")
				}

				if which == "invocation" {
					cancelInvocation(cause)
				} else {
					cancelConsumption(cause)
				}

				<-ctx.Done()

				return nil
			})
			if !errors.Is(err, context.Canceled) || !errors.Is(err, cause) || errors.Is(err, encoding.ErrOutputClosed) {
				t.Fatalf("cancellation=%v", err)
			}

			if err := o.Close(); err != nil {
				t.Fatalf("cancellation polluted cleanup: %v", err)
			}
		})
	}
}

func TestOutputPanicFinalizesAndPropagates(t *testing.T) {
	cleanupErr := errors.New("cleanup")
	o := newOutput(t.Context(), runOutcome{content: &encoding.Content{}, cleanupErr: cleanupErr})
	func() {
		defer func() {
			if recover() != "consumer panic" {
				t.Fatal("panic was swallowed or replaced")
			}
		}()
		_ = o.Consume(t.Context(), func(context.Context, []byte) error { panic("consumer panic") })
	}()
	if err := o.Close(); !errors.Is(err, cleanupErr) {
		t.Fatal(err)
	}

	if o.content != nil {
		t.Fatal("panic retained borrowed storage")
	}
}

func TestOutputPreservesBothCancellationCauses(t *testing.T) {
	for _, invocationFirst := range []bool{false, true} {
		t.Run(map[bool]string{false: "consumption first", true: "invocation first"}[invocationFirst], func(t *testing.T) {
			invocation, cancelInvocation := context.WithCancelCause(t.Context())
			consumption, cancelConsumption := context.WithCancelCause(t.Context())
			defer cancelInvocation(nil)
			defer cancelConsumption(nil)
			invocationCause, consumptionCause := errors.New("invocation"), errors.New("consumption")
			o := newOutput(invocation, runOutcome{content: &encoding.Content{}})
			err := o.Consume(consumption, func(ctx context.Context, _ []byte) error {
				if invocationFirst {
					cancelInvocation(invocationCause)
					<-ctx.Done()
					cancelConsumption(consumptionCause)
				} else {
					cancelConsumption(consumptionCause)
					<-ctx.Done()
					cancelInvocation(invocationCause)
				}

				return nil
			})
			for _, cause := range []error{context.Canceled, invocationCause, consumptionCause} {
				if !errors.Is(err, cause) {
					t.Fatalf("lost %v: %v", cause, err)
				}
			}
		})
	}
}

func TestOutputCloseWaitsForCommittedFinalization(t *testing.T) {
	o := newOutput(t.Context(), runOutcome{content: &encoding.Content{}})
	c, err := o.admit(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	entered, release := make(chan struct{}), make(chan struct{})
	stop := c.stopInvocation
	c.stopInvocation = func() bool {
		close(entered)
		<-release
		if stop != nil {
			return stop()
		}

		return true
	}
	finished, closed := make(chan error, 1), make(chan error, 1)
	go func() { finished <- o.finish(c, nil) }()
	<-entered
	go func() { closed <- o.Close() }()
	select {
	case err := <-closed:
		t.Fatalf("close skipped finalization: %v", err)
	default:
	}

	close(release)
	if err := <-finished; err != nil {
		t.Fatalf("late close changed committed success: %v", err)
	}

	if err := <-closed; err != nil {
		t.Fatal(err)
	}
}

func TestOutputCompletionAndCloseRace(t *testing.T) {
	for range 100 {
		o := newOutput(t.Context(), runOutcome{content: &encoding.Content{Data: []byte("42")}})
		var callbacks atomic.Int32
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			if err := o.Consume(t.Context(), func(context.Context, []byte) error { callbacks.Add(1); return nil }); err != nil && !errors.Is(err, encoding.ErrOutputClosed) {
				t.Error(err)
			}
		}()
		go func() {
			defer wg.Done()
			if err := o.Close(); err != nil {
				t.Error(err)
			}
		}()
		wg.Wait()
		if callbacks.Load() > 1 {
			t.Fatal("duplicate callback")
		}

		if err := o.Close(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestOutputLengthMismatchPreservesAvailableContent(t *testing.T) {
	o := newOutput(t.Context(), runOutcome{content: &encoding.Content{Metadata: encoding.Metadata{Length: 10, LengthKnown: true}, Data: []byte("42")}})
	content, err := o.Collect(t.Context())
	if !errors.Is(err, runtime.ErrInvalidOperation) || content == nil || string(content.Data) != "42" || content.Metadata.Length != 10 {
		t.Fatalf("mismatch lost content or descriptor: %v %v", content, err)
	}
}

func TestOutputCanceledInvocationPreservesEagerFailuresWithoutClaiming(t *testing.T) {
	invocation, cancel := context.WithCancelCause(t.Context())
	cause, operationErr, cleanupErr := errors.New("invocation"), errors.New("execution"), errors.New("cleanup")
	o := newOutput(invocation, runOutcome{operationErr: operationErr, cleanupErr: cleanupErr})
	cancel(cause)
	_, err := o.Collect(t.Context())
	for _, expected := range []error{context.Canceled, cause, operationErr, cleanupErr} {
		if !errors.Is(err, expected) {
			t.Fatalf("lost %v: %v", expected, err)
		}
	}

	if o.state != outputReady {
		t.Fatal("rejected invocation claimed output")
	}

	if err := o.Close(); !errors.Is(err, cleanupErr) {
		t.Fatal(err)
	}
}
