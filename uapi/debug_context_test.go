package uapi

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MontFerret/api"
	apidebugger "github.com/MontFerret/api/debugger"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestDebugAdapterPreservesContextErrors(t *testing.T) {
	portable := newTestRuntime(t)
	plan, err := portable.CompileDebug(t.Context(), api.NewAnonymousSource("RETURN 1"))
	if err != nil {
		t.Fatal(err)
	}
	defer plan.Close()
	session, err := plan.NewDebugSession(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	canceled, cancel := context.WithCancel(t.Context())
	cancel()
	expired, stop := context.WithDeadline(t.Context(), time.Now().Add(-time.Second))
	defer stop()
	for name, call := range map[string]func(context.Context) error{
		"pause": session.Pause,
		"add": func(ctx context.Context) error {
			_, err := session.SetBreakpoint(ctx, api.Location{Position: api.Position{Line: 1}})

			return err
		},
		"add_at": func(ctx context.Context) error {
			_, err := session.SetBreakpointAt(ctx, api.Location{Position: api.Position{Line: 1}}, apidebugger.BreakpointOptions{})

			return err
		},
		"delete": func(ctx context.Context) error { return session.DeleteBreakpoint(ctx, 1) },
		"list": func(ctx context.Context) error {
			_, err := session.Breakpoints(ctx)

			return err
		},
		"frames": func(ctx context.Context) error {
			_, err := session.Frames(ctx)

			return err
		},
		"locals": func(ctx context.Context) error {
			_, err := session.Locals(ctx)

			return err
		},
		"frame_locals": func(ctx context.Context) error {
			_, err := session.FrameLocals(ctx, 0)

			return err
		},
		"variables": func(ctx context.Context) error {
			_, err := session.Variables(ctx, 1)

			return err
		},
	} {
		t.Run(name, func(t *testing.T) {
			for _, tc := range []struct {
				ctx  context.Context
				want error
			}{
				{nil, runtime.ErrInvalidArgument},
				{canceled, context.Canceled},
				{expired, context.DeadlineExceeded},
			} {
				if err := call(tc.ctx); !errors.Is(err, tc.want) {
					t.Fatalf("got %v, want %v", err, tc.want)
				}
			}
		})
	}
	// Rejected request contexts must not become the retained execution context.
	if _, err := session.Start(t.Context()); err != nil {
		t.Fatal(err)
	}

	if event, err := session.Continue(t.Context()); err != nil || event.Reason != apidebugger.ReasonCompleted {
		t.Fatalf("request cancellation affected execution: %+v, %v", event, err)
	}
}
