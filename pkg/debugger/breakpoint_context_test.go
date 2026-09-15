package debugger

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestIncrementalBreakpointCancellationWhileWaiting(t *testing.T) {
	for _, remove := range []bool{false, true} {
		t.Run(fmt.Sprintf("delete=%t", remove), func(t *testing.T) {
			session := replacementSession(t)
			old := replaceLines(t, session, 1)
			before := session.breakpoints.data.Load()
			session.breakpoints.write <- struct{}{}
			defer func() { <-session.breakpoints.write }()
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			waiting := &breakpointWaitContext{Context: ctx, waiting: make(chan struct{})}
			done := make(chan error, 1)
			go func() {
				if remove {
					done <- session.DeleteBreakpoint(waiting, old[0].ID)

					return
				}
				_, err := session.SetBreakpoint(waiting, source.Location{Position: source.Position{Line: 3}})
				done <- err
			}()
			waitForSignal(t, waiting.waiting, "incremental writer admission")
			cancel()
			if err := receiveLiveError(t, done); !errors.Is(err, context.Canceled) {
				t.Fatalf("waiting writer did not cancel: %v", err)
			}

			if session.breakpoints.data.Load() != before {
				t.Fatal("canceled writer published or consumed an ID")
			}
		})
	}
}

func TestIncrementalBreakpointCancellationCheckpoints(t *testing.T) {
	for _, remove := range []bool{false, true} {
		t.Run(fmt.Sprintf("delete=%t", remove), func(t *testing.T) {
			// Sweep cancellation checks through admission, construction and publication
			// without depending on the exact number of checks in those phases.
			sawFailure, sawSuccess := false, false
			for checkpoint := 1; checkpoint <= 32; checkpoint++ {
				session := replacementSession(t)
				old := replaceLines(t, session, 1, 3)
				before := session.breakpoints.data.Load()
				ctx, cancel := context.WithCancel(t.Context())
				observed := &cancellationCheckpoint{Context: ctx, cancel: cancel, remaining: checkpoint}
				var err error
				if remove {
					err = session.DeleteBreakpoint(observed, old[0].ID)
				} else {
					_, err = session.SetBreakpointAt(observed, source.Location{Position: source.Position{Line: 3}}, BreakpointOptions{})
				}

				cancel()
				if err != nil {
					sawFailure = true
					if !errors.Is(err, context.Canceled) || session.breakpoints.data.Load() != before {
						t.Fatalf("checkpoint %d: error=%v, canceled mutation changed snapshot or IDs", checkpoint, err)
					}

					next, err := session.SetBreakpoint(t.Context(), source.Location{Position: source.Position{Line: 1}})
					if err != nil || next.ID != before.nextID {
						t.Fatalf("checkpoint %d consumed an ID: %+v, %v", checkpoint, next, err)
					}
				} else {
					sawSuccess = true
					values, listErr := session.Breakpoints(t.Context())
					want := len(old) + 1
					if remove {
						want = len(old) - 1
					}

					if listErr != nil || len(values) != want {
						t.Fatalf("post-commit cancellation changed success: %+v, %v", values, listErr)
					}
				}
			}

			if !sawFailure || !sawSuccess {
				t.Fatal("did not cover both pre-commit and post-commit cancellation")
			}
		})
	}
}

func TestBreakpointListingContextAfterClose(t *testing.T) {
	session := replacementSession(t)
	old := replaceLines(t, session, 1)
	if err := session.Close(); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if values, err := session.Breakpoints(ctx); values != nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("closed listing ignored cancellation: %+v, %v", values, err)
	}
	values, err := session.Breakpoints(t.Context())
	if err != nil || !reflect.DeepEqual(values, old) {
		t.Fatalf("closed listing: %+v, %v", values, err)
	}

	values[0].ID = -1
	if again, err := session.Breakpoints(t.Context()); err != nil || !reflect.DeepEqual(again, old) {
		t.Fatalf("closed listing was not detached: %+v, %v", again, err)
	}

	var absent *Session
	if values, err := absent.Breakpoints(t.Context()); values != nil || err != nil {
		t.Fatalf("nil session listing: %+v, %v", values, err)
	}
}
