package debugger

import (
	"context"
	"reflect"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestReplacementPublicationFailurePreservesPreparedSet(t *testing.T) {
	for _, failure := range []string{"cancel", "complete", "close"} {
		t.Run(failure, func(t *testing.T) {
			session := replacementSession(t)
			old := replaceLines(t, session, 1)
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			if err := session.breakpoints.lockBreakpoints(ctx, true); err != nil {
				t.Fatal(err)
			}
			defer func() { <-session.breakpoints.write }()

			// Prepare a complete candidate, then fail at the actual commit boundary.
			resolved, err := session.breakpoints.resolveBreakpoint(source.Location{Position: source.Position{Line: 3}}, BreakpointOptions{}, ctx.Err)
			if err != nil {
				t.Fatal(err)
			}
			resolved.ID = old[0].ID + 1
			candidate, err := newBreakpointSnapshot(session.breakpoints.pointIndex, []Breakpoint{resolved}, resolved.ID+1, ctx.Err)
			if err != nil {
				t.Fatal(err)
			}
			switch failure {
			case "cancel":
				cancel()
			case "complete":
				session.lifecycle.finishTerminal("completed")
			case "close":
				if err := session.Close(); err != nil {
					t.Fatal(err)
				}
			}
			if err := session.breakpoints.publishBreakpoints(ctx, candidate, true); err == nil {
				t.Fatal("failed commit succeeded")
			}
			if got := session.Breakpoints(); !reflect.DeepEqual(got, old) {
				t.Fatalf("failed commit changed snapshot: %+v", got)
			}
			if session.breakpoints.data.Load().nextID != old[0].ID+1 {
				t.Fatal("failed commit advanced identity")
			}
		})
	}
}
