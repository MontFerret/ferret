package engine

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/debugger"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestDebugTerminationRetainsResourceCleanupFailures(t *testing.T) {
	for _, mode := range []string{"cancel", "deadline", "close", "close_runtime_error", "runtime_error"} {
		t.Run(mode, func(t *testing.T) {
			waitCtx, stopWaiting := context.WithTimeout(t.Context(), 10*time.Second)
			t.Cleanup(stopWaiting)
			firstErr, secondErr := errors.New("first cleanup"), errors.New("second cleanup")
			first, second, borrowed := newTrackingJSONCloser("first", "1"), newTrackingJSONCloser("second", "2"), newTrackingJSONCloser("borrowed", "3")
			first.closeErr, second.closeErr = firstErr, secondErr
			entered, release := make(chan struct{}), make(chan struct{})
			releaseHost := sync.OnceFunc(func() { close(release) })
			t.Cleanup(releaseHost)
			runtimeErr := errors.New("host failure")
			ctx, cancel := context.WithCancel(waitCtx)
			t.Cleanup(cancel)
			engine := mustNewEngine(t, WithFunctionsRegistrar(func(ns runtime.Namespace) {
				ns.Function().A0().Add("FIRST_RESOURCE", func(context.Context) (runtime.Value, error) { return first, nil })
				ns.Function().A0().Add("SECOND_RESOURCE", func(context.Context) (runtime.Value, error) { return second, nil })
				ns.Function().A0().Add("BLOCK_UNTIL_FAILURE", func(c context.Context) (runtime.Value, error) {
					close(entered)
					if mode == "deadline" || mode == "runtime_error" {
						<-release
						if mode == "deadline" {
							return nil, context.DeadlineExceeded
						}

						return nil, runtimeErr
					}

					<-c.Done()
					if mode == "close_runtime_error" {
						return nil, runtimeErr
					}

					return nil, c.Err()
				})
			}))
			t.Cleanup(func() { _ = engine.Close() })
			plan, err := engine.CompileDebug(t.Context(), source.NewAnonymous("LET first = FIRST_RESOURCE()\nLET alias = first\nLET second = SECOND_RESOURCE()\nLET borrowed = @borrowed\nBLOCK_UNTIL_FAILURE()\nRETURN [first, alias, second, borrowed]"))
			if err != nil {
				t.Fatal(err)
			}

			t.Cleanup(func() { _ = plan.Close() })
			session, err := plan.NewDebugSession(t.Context(), WithSessionRuntimeParam("borrowed", borrowed))
			if err != nil {
				t.Fatal(err)
			}

			t.Cleanup(func() { cancel(); releaseHost(); _ = session.Close() })
			if _, err := session.Start(ctx); err != nil {
				t.Fatal(err)
			}

			events := make(chan *debugger.Event, 1)
			commandErrors := make(chan error, 1)
			go func() {
				event, err := session.Continue(ctx)
				events <- event
				commandErrors <- err
			}()
			select {
			case <-entered:
			case <-waitCtx.Done():
				t.Fatal(waitCtx.Err())
			}

			closed := make(chan error, 4)
			closing := mode == "close" || mode == "close_runtime_error"
			if closing {
				for range cap(closed) {
					go func() { closed <- session.Close() }()
				}
			} else if mode == "cancel" {
				cancel()
			} else {
				releaseHost()
			}

			var event *debugger.Event
			select {
			case event = <-events:
			case <-waitCtx.Done():
				t.Fatal(waitCtx.Err())
			}

			if err := <-commandErrors; err != nil {
				t.Fatal(err)
			}

			cause := error(context.Canceled)
			if mode == "deadline" {
				cause = context.DeadlineExceeded
			} else if mode == "runtime_error" {
				cause = runtimeErr
				if event == nil || event.Reason != debugger.ReasonRuntimeError || first.closed != 0 || second.closed != 0 {
					t.Fatalf("runtime failure did not retain inspectable state: %+v", event)
				}

				event, err = session.Continue(t.Context())
				if err != nil {
					t.Fatal(err)
				}
			}

			if event == nil || event.Reason != debugger.ReasonTerminated || !errors.Is(event.Error, cause) || !errors.Is(event.Error, firstErr) || !errors.Is(event.Error, secondErr) {
				t.Fatalf("termination lost causes: %+v", event)
			}

			if mode == "close_runtime_error" && !errors.Is(event.Error, runtimeErr) {
				t.Fatalf("close lost concurrent runtime failure: %v", event.Error)
			}

			if first.closed != 1 || second.closed != 1 || borrowed.closed != 0 {
				t.Fatalf("closes before teardown: first=%d second=%d borrowed=%d", first.closed, second.closed, borrowed.closed)
			}

			if !closing {
				for range cap(closed) {
					go func() { closed <- session.Close() }()
				}
			}

			for range cap(closed) {
				var closeErr error
				select {
				case closeErr = <-closed:
				case <-waitCtx.Done():
					t.Fatal(waitCtx.Err())
				}

				if !errors.Is(closeErr, firstErr) || !errors.Is(closeErr, secondErr) || errors.Is(closeErr, cause) {
					t.Fatalf("close must retain only cleanup failures: %v", closeErr)
				}
			}

			if first.closed != 1 || second.closed != 1 || borrowed.closed != 0 {
				t.Fatal("repeated close changed resource ownership")
			}
		})
	}
}

func BenchmarkDebugSessionCancellationCleanup(b *testing.B) {
	var cancel context.CancelFunc
	resource := newTrackingJSONCloser("owned", "1")
	engine, err := New(WithFunctionsRegistrar(func(ns runtime.Namespace) {
		ns.Function().A0().Add("BENCH_RESOURCE", func(context.Context) (runtime.Value, error) { return resource, nil })
		ns.Function().A0().Add("BENCH_CANCEL", func(context.Context) (runtime.Value, error) {
			cancel()

			return nil, context.Canceled
		})
	}))
	if err != nil {
		b.Fatal(err)
	}

	defer engine.Close()
	plan, err := engine.CompileDebug(b.Context(), source.NewAnonymous("LET resource = BENCH_RESOURCE()\nBENCH_CANCEL()\nRETURN resource"))
	if err != nil {
		b.Fatal(err)
	}

	defer plan.Close()
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		var ctx context.Context
		ctx, cancel = context.WithCancel(b.Context())
		session, err := plan.NewDebugSession(ctx)
		if err != nil {
			b.Fatal(err)
		}

		if _, err := session.Start(ctx); err != nil {
			b.Fatal(err)
		}

		event, err := session.Continue(ctx)
		if err != nil || event == nil || !errors.Is(event.Error, context.Canceled) {
			b.Fatalf("cancellation failed: event=%+v err=%v", event, err)
		}

		if err := session.Close(); err != nil {
			b.Fatal(err)
		}

		cancel()
		if resource.closed != 1 {
			b.Fatalf("resource closes=%d", resource.closed)
		}

		resource.closed = 0
	}
}
