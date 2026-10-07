package engine

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestSessionAfterRunFailurePreservesOutputAndResourceOwnership(t *testing.T) {
	hookErr := errors.New("after run failed")
	closeErr := errors.New("result cleanup failed")
	encodeErr := errors.New("result encoding failed")

	for _, test := range []struct {
		encodeErr   error
		closeErr    error
		outputErr   error
		name        string
		contentType string
	}{
		{name: "successful output", contentType: "application/json"},
		{name: "output and cleanup failure", contentType: "application/json", closeErr: closeErr},
		{name: "missing codec", contentType: "application/missing", closeErr: closeErr, outputErr: encoding.ErrCodecNotFound},
		{name: "encoding failure", contentType: "application/json", encodeErr: encodeErr, closeErr: closeErr, outputErr: encodeErr},
	} {
		t.Run(test.name, func(t *testing.T) {
			resource := newTrackingJSONCloser("owned result", "42")
			resource.closeErr = test.closeErr
			resource.encodeErr = test.encodeErr
			engine := mustNewEngine(t,
				WithFunctionsRegistrar(func(ns runtime.Namespace) {
					ns.Function().A0().Add("OWNED_RESULT", func(context.Context) (runtime.Value, error) { return resource, nil })
				}),
				WithAfterRunHook(func(_ context.Context, err error) error {
					if err != nil {
						t.Errorf("after run received %v for successful execution", err)
					}

					resource.events = append(resource.events, "after")

					return hookErr
				}),
			)
			t.Cleanup(func() { _ = engine.Close() })
			plan := mustCompilePlan(t, engine, "RETURN OWNED_RESULT()")
			t.Cleanup(func() { _ = plan.Close() })
			session := mustNewSession(t, plan, WithOutputContentType(test.contentType))
			t.Cleanup(func() { _ = session.Close() })

			handle, err := session.Run(t.Context())
			if err != nil || handle == nil {
				t.Fatalf("terminal failure escaped Run: %v", err)
			}

			if resource.closed != 1 {
				t.Fatal("Run retained its materialized VM result")
			}

			output, err := handle.Collect(t.Context())
			for _, expected := range []error{hookErr, test.closeErr, test.outputErr} {
				if expected != nil && !errors.Is(err, expected) {
					t.Errorf("lost error %v: %v", expected, err)
				}
			}

			if test.outputErr == nil {
				if output == nil || output.Metadata.ContentType != test.contentType || string(output.Data) != "42" {
					t.Fatalf("lost successful output: %+v", output)
				}
			} else if output != nil {
				t.Fatalf("failed encoding returned output: %+v", output)
			}

			wantEvents := []string{"after", "encode", "close"}
			if test.outputErr == encoding.ErrCodecNotFound {
				wantEvents = []string{"after", "close"}
			}

			if !slices.Equal(resource.events, wantEvents) || resource.closed != 1 {
				t.Fatalf("resource lifecycle before teardown: events=%v closes=%d", resource.events, resource.closed)
			}

			for range 2 {
				if err := handle.Close(); !errors.Is(err, test.closeErr) || (test.closeErr == nil && err != nil) || errors.Is(err, hookErr) {
					t.Fatalf("output cleanup outcome: %v", err)
				}
			}

			for _, closeResource := range []func() error{session.Close, plan.Close, engine.Close} {
				if err := closeResource(); err != nil {
					t.Fatal(err)
				}
			}

			if !slices.Equal(resource.events, wantEvents) || resource.closed != 1 {
				t.Fatalf("teardown repeated result cleanup or hooks: events=%v closes=%d", resource.events, resource.closed)
			}

			if output != nil && string(output.Data) != "42" {
				t.Fatalf("teardown invalidated returned output: %+v", output)
			}
		})
	}
}

func TestRunPanicClosesMaterializedResultAndTemporaryResources(t *testing.T) {
	for _, stage := range []string{"after hook", "encoder"} {
		t.Run(stage, func(t *testing.T) {
			resource := newTrackingJSONCloser("owned result", "42")
			var value runtime.Value = resource
			if stage == "encoder" {
				value = &panickingJSONCloser{trackingJSONCloser: resource}
			}

			var sessionCloses, planCloses int
			panicEnabled := true
			engine := mustNewEngine(t,
				WithMaxActiveSessions(1),
				WithFunctionsRegistrar(func(ns runtime.Namespace) {
					ns.Function().A0().Add("OWNED_RESULT", func(context.Context) (runtime.Value, error) { return value, nil })
				}),
				WithAfterRunHook(func(context.Context, error) error {
					if stage == "after hook" && panicEnabled {
						panic("after hook panic")
					}

					return nil
				}),
				WithSessionCloseHook(func() error { sessionCloses++; return nil }),
				WithPlanCloseHook(func() error { planCloses++; return nil }),
			)
			t.Cleanup(func() { _ = engine.Close() })
			func() {
				defer func() {
					if got := recover(); got != stage+" panic" {
						t.Fatalf("unexpected panic: %v", got)
					}
				}()
				_, _ = engine.Run(t.Context(), source.NewAnonymous("RETURN OWNED_RESULT()"))
			}()

			if resource.closed != 1 || sessionCloses != 1 || planCloses != 1 {
				t.Fatalf("panic cleanup: result=%d session=%d plan=%d", resource.closed, sessionCloses, planCloses)
			}

			panicEnabled = false
			ctx, cancel := context.WithTimeout(t.Context(), time.Second)
			defer cancel()
			handle, err := engine.Run(ctx, source.NewAnonymous("RETURN 1"))
			if err != nil {
				t.Fatalf("panic retained capacity: %v", err)
			}

			if _, err := handle.Collect(ctx); err != nil {
				t.Fatal(err)
			}
		})
	}
}
