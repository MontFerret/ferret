package uapi

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/MontFerret/api"
	apidebugger "github.com/MontFerret/api/debugger"
	"github.com/MontFerret/ferret/v2/pkg/engine"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestPortableBreakpointReplacementReachesRunningNativeSession(t *testing.T) {
	for _, tc := range []struct {
		name    string
		initial []int
		updates [][]int
		want    int
	}{
		{name: "add_and_relocate", updates: [][]int{{2}}, want: 3},
		{name: "remove", initial: []int{3, 4}, updates: [][]int{{4}}, want: 4},
		{name: "replace", initial: []int{3}, updates: [][]int{{4}}, want: 4},
		{name: "clear", initial: []int{3}, updates: [][]int{{}}},
		{name: "rapid", updates: [][]int{{3}, {4}, {}, {2}}, want: 3},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
			defer cancel()
			entered, release := make(chan struct{}), make(chan struct{})
			native, err := engine.New(engine.WithFunctionsRegistrar(func(ns runtime.Namespace) {
				ns.Function().A0().Add("WAIT_NATIVE", func(ctx context.Context) (runtime.Value, error) {
					close(entered)
					select {
					case <-release:
						return runtime.NewInt(1), nil
					case <-ctx.Done():
						return nil, ctx.Err()
					}
				})
			}))
			if err != nil {
				t.Fatal(err)
			}
			defer native.Close()
			portable := Wrap(native)
			plan, err := portable.CompileDebug(ctx, api.NewSource("live.fql", "LET value = WAIT_NATIVE()\n\nLET next = value + 1\nRETURN next"))
			if err != nil {
				t.Fatal(err)
			}
			defer plan.Close()
			session, err := plan.NewDebugSession(ctx)
			if err != nil {
				t.Fatal(err)
			}
			defer session.Close()

			replace := func(lines []int) []apidebugger.Breakpoint {
				t.Helper()
				requests := make([]apidebugger.BreakpointRequest, len(lines))
				for i, line := range lines {
					requests[i].Position.Line = line
				}
				results, err := session.ReplaceBreakpoints(ctx, "live.fql", requests)
				if err != nil {
					t.Fatal(err)
				}

				return results
			}
			replace(tc.initial)
			if _, err := session.Start(ctx); err != nil {
				t.Fatal(err)
			}
			done := make(chan *apidebugger.Event, 1)
			go func() {
				event, err := session.Continue(ctx)
				if err != nil {
					t.Errorf("continue: %v", err)
				}
				done <- event
			}()
			select {
			case <-entered:
			case <-ctx.Done():
				t.Fatal(ctx.Err())
			}
			var results []apidebugger.Breakpoint
			for _, update := range tc.updates {
				results = replace(update)
			}
			before := session.Breakpoints()
			canceled, stop := context.WithCancel(ctx)
			stop()
			if _, err := session.ReplaceBreakpoints(canceled, "live.fql", nil); !errors.Is(err, context.Canceled) {
				t.Fatalf("adapter lost context cancellation: %v", err)
			}
			if !reflect.DeepEqual(before, session.Breakpoints()) {
				t.Fatal("cancellation changed active state")
			}
			close(release)
			select {
			case event := <-done:
				if event == nil {
					t.Fatal("missing event")
				}
				if tc.want == 0 {
					if event.Reason != apidebugger.ReasonCompleted {
						t.Fatalf("clear paused execution: %+v", event)
					}
				} else if event.Reason != apidebugger.ReasonBreakpoint || event.Location.Line != tc.want ||
					!reflect.DeepEqual(event.HitBreakpointIDs, []apidebugger.BreakpointID{results[0].ID}) {
					t.Fatalf("unexpected live result: %+v; breakpoints %+v", event, results)
				}
			case <-ctx.Done():
				t.Fatal(ctx.Err())
			}
		})
	}
}
