package uapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/MontFerret/api"
	apidebugger "github.com/MontFerret/api/debugger"
	apidiagnostics "github.com/MontFerret/api/diagnostics"
	"github.com/MontFerret/ferret/v2"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/uapi"
)

func TestRealOutputEagerExecutionAndEarlyCleanup(t *testing.T) {
	for _, portable := range []bool{false, true} {
		for _, collect := range []bool{false, true} {
			t.Run(map[bool]string{false: "native", true: "uapi"}[portable]+map[bool]string{false: "/consume", true: "/collect"}[collect], func(t *testing.T) {
				var executions, sessions, plans atomic.Int32
				_, run := newOutputExecutor(t, portable,
					ferret.WithMaxActiveSessions(1),
					ferret.WithFunctionsRegistrar(func(ns runtime.Namespace) {
						ns.Function().A0().Add("PROBE", func(context.Context) (runtime.Value, error) { executions.Add(1); return runtime.NewInt(42), nil })
					}),
					ferret.WithSessionCloseHook(func() error { sessions.Add(1); return nil }),
					ferret.WithPlanCloseHook(func() error { plans.Add(1); return nil }),
				)
				o, err := run(t.Context(), "RETURN PROBE()")
				if err != nil || o == nil {
					t.Fatalf("admission=%v output=%v", err, o)
				}

				if executions.Load() != 1 || sessions.Load() != 1 || plans.Load() != 1 {
					t.Fatal("Run deferred execution or retained temporary resources")
				}

				metadata := o.Metadata()
				if metadata != (api.Metadata{ContentType: "application/json", Length: 2, LengthKnown: true}) {
					t.Fatal(metadata)
				}

				// An unread output must not retain the temporary session's permit or VM.
				ctx, cancel := context.WithTimeout(t.Context(), time.Second)
				defer cancel()
				other, err := run(ctx, "RETURN 1")
				if err != nil {
					t.Fatal(err)
				}

				if _, err := other.Collect(ctx); err != nil {
					t.Fatal(err)
				}

				if collect {
					content, err := o.Collect(t.Context())
					if err != nil || content == nil || string(content.Data) != "42" {
						t.Fatalf("collect=%v %v", content, err)
					}
				} else {
					var delivered []byte
					err := o.Consume(t.Context(), func(_ context.Context, chunk []byte) error { delivered = append(delivered, chunk...); return nil })
					if err != nil || string(delivered) != "42" {
						t.Fatalf("consume=%q %v", delivered, err)
					}
				}

				if executions.Load() != 1 || o.Metadata() != metadata {
					t.Fatal("consumption repeated execution or changed metadata")
				}

				if err := o.Close(); err != nil {
					t.Fatal(err)
				}

				if sessions.Load() != 2 || plans.Load() != 2 {
					t.Fatal("output repeated native cleanup")
				}
			})
		}
	}
}

func TestRealOutputEncodingPresenceAndDetachment(t *testing.T) {
	failure := errors.New("encoding failed")
	for _, portable := range []bool{false, true} {
		for _, tc := range []struct {
			err     error
			name    string
			data    []byte
			present bool
		}{
			{name: "nil empty", present: true}, {name: "slice empty", data: []byte{}, present: true},
			{name: "populated", data: []byte("42"), present: true},
			{name: "absent failure", err: failure}, {name: "prefix failure", data: []byte("4"), err: failure, present: true},
			{name: "empty failure", data: []byte{}, err: failure},
		} {
			t.Run(map[bool]string{false: "native/", true: "uapi/"}[portable]+tc.name, func(t *testing.T) {
				codec := &detachedTestCodec{data: tc.data, err: tc.err}
				native, _ := newOutputExecutor(t, portable, ferret.WithEncodingCodec(codec.ContentType(), codec))
				run := func(ctx context.Context, query string) (api.Output, error) {
					if portable {
						return uapi.Wrap(native, "test").Run(ctx, api.NewSource("output.fql", query), api.WithOutputContentType(codec.ContentType()))
					}

					return native.Run(ctx, source.New("output.fql", query), ferret.WithOutputContentType(codec.ContentType()))
				}
				o, err := run(t.Context(), "RETURN 42")
				if err != nil || o == nil {
					t.Fatalf("encoding failure escaped Run: %v", err)
				}

				content, err := o.Collect(t.Context())
				if (content != nil) != tc.present || !errors.Is(err, tc.err) {
					t.Fatalf("content=%v err=%v", content, err)
				}

				if content != nil {
					if !bytes.Equal(content.Data, tc.data) {
						t.Fatal("encoding changed bytes")
					}

					if content.Metadata.ContentType != codec.ContentType() || content.Metadata.LengthKnown != (tc.err == nil) {
						t.Fatal(content.Metadata)
					}

					if tc.err == nil && content.Metadata.Length != int64(len(tc.data)) {
						t.Fatal(content.Metadata)
					}

					codec.mu.Lock()
					if len(codec.data) > 0 {
						codec.data[0] = 'x'
					}

					codec.mu.Unlock()
					if len(content.Data) > 0 && content.Data[0] == 'x' {
						t.Fatal("content aliases reusable encoder storage")
					}

					before := bytes.Clone(content.Data)
					_ = native.Close()
					if !bytes.Equal(content.Data, before) {
						t.Fatal("parent cleanup invalidated content")
					}
				}

				if err := o.Close(); err != nil {
					t.Fatalf("encoding failure became cleanup: %v", err)
				}
			})
		}
	}
}

func TestRealOutputTerminalDiagnosticsAndJoinedCauses(t *testing.T) {
	for _, method := range []string{"Collect", "Consume", "Close"} {
		t.Run(method, func(t *testing.T) {
			hookCause, cleanupCause, consumerCause := errors.New("hook cause"), errors.New("cleanup cause"), errors.New("consumer cause")
			src := source.New("output.fql", "RETURN 42")
			hookErr := diagnostics.NewUnexpectedErrorWith(src, "after hook", hookCause)
			cleanupErr := diagnostics.NewUnexpectedErrorWith(src, "close hook", cleanupCause)
			_, run := newOutputExecutor(t, true, ferret.WithAfterRunHook(func(context.Context, error) error { return hookErr }), ferret.WithSessionCloseHook(func() error { return cleanupErr }))
			o, err := run(t.Context(), src.Content())
			if err != nil || o == nil {
				t.Fatalf("terminal errors escaped admission: %v", err)
			}

			switch method {
			case "Collect":
				var content *api.Content
				content, err = o.Collect(t.Context())
				if content == nil || string(content.Data) != "42" {
					t.Fatal("lost content with error")
				}
			case "Consume":
				err = o.Consume(t.Context(), func(context.Context, []byte) error { return consumerCause })
			case "Close":
				err = o.Close()
			}

			var projected apidiagnostics.Diagnostics
			if !errors.As(err, &projected) || !errors.Is(err, cleanupCause) || !errors.Is(err, cleanupErr) {
				t.Fatalf("missing cleanup diagnostic: %v", err)
			}

			if method != "Close" && (!errors.Is(err, hookCause) || !errors.Is(err, hookErr)) {
				t.Fatal("lost hook diagnostic")
			}

			if method == "Consume" && !errors.Is(err, consumerCause) {
				t.Fatal("lost consumer error")
			}

			if method == "Close" && errors.Is(err, hookCause) {
				t.Fatal("Close reported execution outcome")
			}

			for _, diagnostic := range projected {
				if diagnostic.Source.Name != src.Name() || diagnostic.Source.Content != src.Content() {
					t.Fatal("source lost")
				}
			}

			if err := o.Close(); !errors.Is(err, cleanupCause) {
				t.Fatal("repeated close lost cleanup")
			}
		})
	}

	_, run := newOutputExecutor(t, true)
	if o, err := run(t.Context(), "RETURN"); err == nil || o != nil {
		t.Fatal("compilation failure crossed admission")
	}

	o, err := run(t.Context(), "RETURN 1 / 0")
	if err != nil || o == nil {
		t.Fatal("VM failure escaped admission")
	}

	content, err := o.Collect(t.Context())
	var projected apidiagnostics.Diagnostics
	if content != nil || !errors.As(err, &projected) {
		t.Fatalf("VM terminal diagnostic=%v %v", content, err)
	}

	if err := o.Close(); err != nil {
		t.Fatal("VM failure polluted cleanup")
	}
}

func TestRealOutputCanceledAdmissionPreservesRecordedDiagnostics(t *testing.T) {
	for _, portable := range []bool{false, true} {
		for _, method := range []string{"Collect", "Consume"} {
			for _, duringRun := range []bool{false, true} {
				for _, shared := range []bool{false, true} {
					name := map[bool]string{false: "native", true: "uapi"}[portable] + "/" + method +
						map[bool]string{false: "/after Run", true: "/during Run"}[duringRun] +
						map[bool]string{false: "/distinct contexts", true: "/shared context"}[shared]
					t.Run(name, func(t *testing.T) {
						invocation, cancelInvocation := context.WithCancelCause(t.Context())
						consumption, cancelConsumption := context.WithCancelCause(t.Context())
						defer cancelInvocation(nil)
						defer cancelConsumption(nil)
						invocationCause, consumptionCause := errors.New("invocation canceled"), errors.New("consumption canceled")
						hookCause, cleanupCause := errors.New("after-run failure"), errors.New("temporary-session cleanup failure")
						src := source.New("output.fql", "RETURN 42")
						hookErr := diagnostics.NewUnexpectedErrorWith(src, "after hook", hookCause)
						cleanupErr := diagnostics.NewUnexpectedErrorWith(src, "close hook", cleanupCause)
						var cleanupCalls atomic.Int32
						_, run := newOutputExecutor(t, portable,
							ferret.WithAfterRunHook(func(context.Context, error) error {
								if duringRun {
									cancelInvocation(invocationCause)
								}

								return hookErr
							}),
							ferret.WithSessionCloseHook(func() error { cleanupCalls.Add(1); return cleanupErr }),
						)
						o, err := run(invocation, src.Content())
						if err != nil || o == nil {
							t.Fatalf("recorded failures escaped Run: output=%v err=%v", o, err)
						}

						defer o.Close()
						if cleanupCalls.Load() != 1 {
							t.Fatal("Run did not finish temporary-session cleanup")
						}

						cancelInvocation(invocationCause)
						if shared {
							consumption = invocation
						} else {
							cancelConsumption(consumptionCause)
						}

						if method == "Collect" {
							var content *api.Content
							content, err = o.Collect(consumption)
							if content != nil {
								t.Fatal("rejected collection returned content")
							}
						} else {
							err = o.Consume(consumption, func(context.Context, []byte) error { t.Fatal("rejected callback"); return nil })
						}

						for _, cause := range []error{context.Canceled, invocationCause, hookCause, hookErr, cleanupCause, cleanupErr} {
							if !errors.Is(err, cause) {
								t.Errorf("lost recorded cause %v: %v", cause, err)
							}
						}

						if !shared && !errors.Is(err, consumptionCause) {
							t.Errorf("lost separate consumption cause: %v", err)
						}

						if portable {
							var projected apidiagnostics.Diagnostics
							if !errors.As(err, &projected) || len(projected) != 2 {
								t.Fatalf("lost recorded diagnostic projection: %v", err)
							}

							messages := []string{hookErr.Message, cleanupErr.Message}
							for index, diagnostic := range projected {
								if diagnostic.Message != messages[index] {
									t.Errorf("recorded diagnostic message=%q want=%q", diagnostic.Message, messages[index])
								}

								if diagnostic.Source.Name != src.Name() || diagnostic.Source.Content != src.Content() {
									t.Error("recorded diagnostic source lost")
								}
							}
						}

						for range 2 {
							if err := o.Close(); !errors.Is(err, cleanupCause) || !errors.Is(err, cleanupErr) || errors.Is(err, hookCause) || errors.Is(err, context.Canceled) {
								t.Fatalf("Close did not preserve cleanup-only outcome: %v", err)
							}
						}

						if cleanupCalls.Load() != 1 {
							t.Fatal("rejected admission or Close repeated session cleanup")
						}
					})
				}
			}
		}
	}
}

func TestRealOutputContextsAndSessionReuse(t *testing.T) {
	for _, portable := range []bool{false, true} {
		t.Run(map[bool]string{false: "native", true: "uapi"}[portable], func(t *testing.T) {
			var executionCancel context.CancelFunc
			native, run := newOutputExecutor(t, portable,
				ferret.WithBeforeRunHook(func(ctx context.Context) (context.Context, error) {
					derived, cancel := context.WithCancel(ctx)
					executionCancel = cancel

					return derived, nil
				}),
				ferret.WithAfterRunHook(func(context.Context, error) error { executionCancel(); return nil }),
			)
			o, err := run(t.Context(), "RETURN 42")
			if err != nil {
				t.Fatal(err)
			}

			if content, err := o.Collect(t.Context()); err != nil || string(content.Data) != "42" {
				t.Fatalf("execution-only context canceled output: %v", err)
			}

			invocation, cancel := context.WithCancelCause(t.Context())
			o, err = run(invocation, "RETURN 42")
			if err != nil {
				t.Fatal(err)
			}

			cause := errors.New("invocation canceled after Run")
			cancel(cause)
			if _, err := o.Collect(t.Context()); !errors.Is(err, context.Canceled) || !errors.Is(err, cause) {
				t.Fatal(err)
			}

			if err := o.Close(); err != nil {
				t.Fatal(err)
			}

			var runSession func(context.Context) (api.Output, error)
			if portable {
				plan, err := uapi.Wrap(native, "test").Compile(t.Context(), api.NewAnonymousSource("RETURN 42"))
				if err != nil {
					t.Fatal(err)
				}

				defer plan.Close()
				session, err := plan.NewSession(t.Context())
				if err != nil {
					t.Fatal(err)
				}

				defer session.Close()
				runSession = session.Run
			} else {
				plan, err := native.Compile(t.Context(), source.NewAnonymous("RETURN 42"))
				if err != nil {
					t.Fatal(err)
				}

				defer plan.Close()
				session, err := plan.NewSession(t.Context())
				if err != nil {
					t.Fatal(err)
				}

				defer session.Close()
				runSession = session.Run
			}

			var previous *api.Content
			for range 2 {
				o, err := runSession(t.Context())
				if err != nil {
					t.Fatal(err)
				}

				content, err := o.Collect(t.Context())
				if err != nil {
					t.Fatal(err)
				}

				if previous != nil {
					previous.Data[0] = '9'
					if string(content.Data) != "42" {
						t.Fatal("later output aliases prior collection")
					}
				}

				previous = content
				if err := o.Close(); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestCollectedContentSerialization(t *testing.T) {
	_, run := newOutputExecutor(t, true)
	o, err := run(t.Context(), "RETURN 42")
	if err != nil {
		t.Fatal(err)
	}

	content, err := o.Collect(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	data, err := json.Marshal(content)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != `{"metadata":{"contentType":"application/json","length":2,"lengthKnown":true},"data":"NDI="}` {
		t.Fatalf("serialization=%s", data)
	}
}

func TestRealOutputConsumePreservesWrittenPrefixAndError(t *testing.T) {
	for _, portable := range []bool{false, true} {
		t.Run(map[bool]string{false: "native", true: "uapi"}[portable], func(t *testing.T) {
			failure := errors.New("partial encoding")
			codec := &detachedTestCodec{data: []byte("4"), err: failure}
			native, _ := newOutputExecutor(t, portable, ferret.WithEncodingCodec(codec.ContentType(), codec))
			var out api.Output
			var err error
			if portable {
				out, err = uapi.Wrap(native, "test").Run(t.Context(), api.NewAnonymousSource("RETURN 42"), api.WithOutputContentType(codec.ContentType()))
			} else {
				out, err = native.Run(t.Context(), source.NewAnonymous("RETURN 42"), ferret.WithOutputContentType(codec.ContentType()))
			}

			if err != nil {
				t.Fatal(err)
			}

			calls := 0
			err = out.Consume(t.Context(), func(_ context.Context, data []byte) error {
				calls++
				if string(data) != "4" {
					t.Fatal("consumption or adaptation changed the written prefix")
				}

				return nil
			})
			if calls != 1 || !errors.Is(err, failure) || out.Metadata().LengthKnown {
				t.Fatalf("partial consumption: calls=%d metadata=%v err=%v", calls, out.Metadata(), err)
			}
		})
	}
}

func TestDebugCompletionRetainsPartialContentAndErrors(t *testing.T) {
	encodingErr, hookErr := errors.New("partial encoding"), errors.New("after hook")
	codec := &detachedTestCodec{data: []byte("4"), err: encodingErr}
	native, _ := newOutputExecutor(t, true,
		ferret.WithEncodingCodec(codec.ContentType(), codec),
		ferret.WithAfterRunHook(func(_ context.Context, err error) error {
			if err != nil || !codec.written {
				t.Error("debugger after hook lost its completion ordering")
			}

			return hookErr
		}),
	)
	plan, err := uapi.Wrap(native, "test").CompileDebug(t.Context(), api.NewAnonymousSource("RETURN 42"))
	if err != nil {
		t.Fatal(err)
	}

	defer plan.Close()
	session, err := plan.NewDebugSession(t.Context(), api.WithOutputContentType(codec.ContentType()))
	if err != nil {
		t.Fatal(err)
	}

	defer session.Close()
	if _, err := session.Start(t.Context()); err != nil {
		t.Fatal(err)
	}

	event, err := session.Continue(t.Context())
	if event == nil || event.Reason != apidebugger.ReasonCompleted || event.Output == nil || !errors.Is(err, encodingErr) || !errors.Is(err, hookErr) {
		t.Fatalf("completion=%v err=%v", event, err)
	}

	if event.Output.Metadata.LengthKnown || string(event.Output.Data) != "4" {
		t.Fatal(event.Output)
	}

	codec.data[0] = '9'
	_ = session.Close()
	_ = plan.Close()
	_ = native.Close()
	if string(event.Output.Data) != "4" {
		t.Fatal("completion content did not survive encoder reuse and cleanup")
	}
}

func TestRealOutputConcurrentCloseWaitsThroughAdapter(t *testing.T) {
	for _, portable := range []bool{false, true} {
		t.Run(map[bool]string{false: "native", true: "uapi"}[portable], func(t *testing.T) {
			var cleanupCalls atomic.Int32
			_, run := newOutputExecutor(t, portable, ferret.WithSessionCloseHook(func() error { cleanupCalls.Add(1); return nil }))
			o, err := run(t.Context(), "RETURN 42")
			if err != nil {
				t.Fatal(err)
			}

			if cleanupCalls.Load() != 1 {
				t.Fatal("temporary session retained by output")
			}

			entered, canceled, release := make(chan struct{}), make(chan struct{}), make(chan struct{})
			consumed, closed := make(chan error, 1), make(chan error, 3)
			consumerErr := errors.New("consumer stopped")
			go func() {
				consumed <- o.Consume(t.Context(), func(ctx context.Context, data []byte) error {
					close(entered)
					<-ctx.Done()
					close(canceled)
					<-release
					if string(data) != "42" {
						t.Error("closure invalidated borrowed bytes")
					}

					return consumerErr
				})
			}()
			<-entered
			if _, err := o.Collect(t.Context()); !errors.Is(err, api.ErrOutputInUse) {
				t.Fatal(err)
			}

			for range 3 {
				go func() { closed <- o.Close() }()
			}

			<-canceled
			select {
			case err := <-closed:
				t.Fatalf("closer returned with active callback: %v", err)
			default:
			}

			close(release)
			if err := <-consumed; !errors.Is(err, api.ErrOutputClosed) || !errors.Is(err, consumerErr) {
				t.Fatalf("terminal error=%v", err)
			}

			for range 3 {
				if err := <-closed; err != nil {
					t.Fatal(err)
				}
			}

			if cleanupCalls.Load() != 1 {
				t.Fatal("output repeated session cleanup")
			}
		})
	}
}

func newOutputExecutor(t *testing.T, portable bool, opts ...ferret.Option) (*ferret.Engine, func(context.Context, string) (api.Output, error)) {
	t.Helper()
	native, err := ferret.New(opts...)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = native.Close() })
	if portable {
		adapter := uapi.Wrap(native, "test")

		return native, func(ctx context.Context, query string) (api.Output, error) {
			return adapter.Run(ctx, api.NewSource("output.fql", query))
		}
	}

	return native, func(ctx context.Context, query string) (api.Output, error) {
		return native.Run(ctx, source.New("output.fql", query))
	}
}
