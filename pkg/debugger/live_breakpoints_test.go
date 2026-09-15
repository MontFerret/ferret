package debugger

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestSessionReplacesBreakpointsWhileRunning(t *testing.T) {
	for _, tc := range []struct {
		name    string
		initial []int
		updates [][]int
		want    int
	}{
		{name: "add", updates: [][]int{{3}}, want: 3},
		{name: "remove", initial: []int{3, 4}, updates: [][]int{{4}}, want: 4},
		{name: "replace", initial: []int{3}, updates: [][]int{{4}}, want: 4},
		{name: "clear", initial: []int{3, 4}, updates: [][]int{{}}},
		{name: "rapid", initial: []int{3}, updates: [][]int{{4}, {3, 4}, {}, {3}, {4}}, want: 4},
	} {
		t.Run(tc.name, func(t *testing.T) {
			session, gate := newLiveSession(t, nil)
			replaceLines(t, session, tc.initial...)
			done := continueLive(t, session, t.Context())
			gate.visit(t)
			var installed []Breakpoint
			for _, lines := range tc.updates {
				installed = replaceLines(t, session, lines...)
			}
			gate.proceed(t)

			if tc.want != 0 {
				event := liveEvent(t, done)
				if event.Reason != ReasonBreakpoint || event.Location.Line != tc.want ||
					len(event.HitBreakpointIDs) != 1 || event.HitBreakpointIDs[0] != installed[0].ID {
					t.Fatalf("unexpected live hit: %+v; installed %+v", event, installed)
				}

				done = continueLive(t, session, t.Context())
			}

			// Clear during the next run; both remaining visits must finish without a hit.
			gate.visit(t)
			replaceLines(t, session)
			gate.proceed(t)
			gate.visit(t)
			gate.proceed(t)
			if event := liveEvent(t, done); event.Reason != ReasonCompleted {
				t.Fatalf("cleared breakpoint stopped a subsequent visit: %+v", event)
			}
		})
	}
}

func TestSessionRetainsCommittedHitAfterReplacement(t *testing.T) {
	hold := &deferredStopExecution{
		reason: vm.DebugStopBreakpoint, ready: make(chan *vm.DebugExecutionEvent), release: make(chan struct{}),
	}
	session, gate := newLiveSession(t, hold)
	old := replaceLines(t, session, 3)[0]
	done := continueLive(t, session, t.Context())
	gate.visit(t)
	gate.proceed(t)
	select {
	case <-hold.ready:
	case <-time.After(3 * time.Second):
		t.Fatal("VM did not commit breakpoint hit")
	}

	replaceLines(t, session)
	added := replaceLines(t, session, 3)[0]
	if added.ID == old.ID {
		t.Fatal("removed breakpoint identity was reused")
	}
	close(hold.release)
	event := liveEvent(t, done)
	if event.Reason != ReasonBreakpoint || !reflect.DeepEqual(event.HitBreakpointIDs, []BreakpointID{old.ID}) {
		t.Fatalf("committed hit changed: %+v", event)
	}

	if frames, err := session.Frames(context.Background()); err != nil || len(frames) == 0 {
		t.Fatalf("committed stop is not inspectable: %+v, %v", frames, err)
	}

	if value, err := session.Evaluate(t.Context(), "i"); err != nil || value.Display != "1" {
		t.Fatalf("committed stop lost locals: %+v, %v", value, err)
	}

	// Caller mutation of an old event cannot alter snapshot-owned identities.
	event.HitBreakpointIDs[0] = -1
	done = continueLive(t, session, t.Context())
	gate.visit(t)
	replaceLines(t, session)
	gate.proceed(t)
	gate.visit(t)
	gate.proceed(t)
	if event := liveEvent(t, done); event.Reason != ReasonCompleted {
		t.Fatalf("removed hit remained active: %+v", event)
	}
}

func TestSessionLiveReplacementAndPause(t *testing.T) {
	session, gate := newLiveSession(t, nil)
	done := continueLive(t, session, t.Context())
	gate.visit(t)

	// Simulate another writer holding admission. Pause must not wait for it.
	session.breakpoints.write <- struct{}{}
	ctx, cancel := context.WithCancel(t.Context())
	waiting := &breakpointWaitContext{Context: ctx, waiting: make(chan struct{})}
	replaced := make(chan error, 1)
	go func() {
		_, err := session.ReplaceBreakpoints(waiting, "", lineRequests(3))
		replaced <- err
	}()
	waitForSignal(t, waiting.waiting, "replacement admission")
	if err := session.Pause(context.Background()); err != nil {
		t.Fatal(err)
	}
	gate.proceed(t)
	if event := liveEvent(t, done); event.Reason != ReasonPause {
		t.Fatalf("pause did not interrupt: %+v", event)
	}
	cancel()
	if err := receiveLiveError(t, replaced); !errors.Is(err, context.Canceled) {
		t.Fatalf("waiting replacement did not cancel: %v", err)
	}
	<-session.breakpoints.write

	installed := replaceLines(t, session, 4)
	for range 2 {
		event, err := session.StepIn(t.Context())
		if err != nil {
			t.Fatal(err)
		}
		if event.Reason == ReasonBreakpoint {
			if event.HitBreakpointIDs[0] != installed[0].ID {
				t.Fatalf("wrong hit after pause: %+v", event)
			}

			return
		}
		replaceLines(t, session, 4)
	}
	t.Fatal("stepping did not observe replacement")
}

func TestSessionCloseWakesWaitingReplacement(t *testing.T) {
	session, gate := newLiveSession(t, nil)
	old := replaceLines(t, session, 3)
	done := continueLive(t, session, t.Context())
	gate.visit(t)
	session.breakpoints.write <- struct{}{}
	defer func() { <-session.breakpoints.write }()
	waiting := &breakpointWaitContext{Context: t.Context(), waiting: make(chan struct{})}
	replaced := make(chan error, 1)
	go func() {
		_, err := session.ReplaceBreakpoints(waiting, "", lineRequests(4))
		replaced <- err
	}()
	waitForSignal(t, waiting.waiting, "replacement admission")
	closed := make(chan error, 1)
	go func() { closed <- session.Close() }()
	if err := receiveLiveError(t, closed); err != nil {
		t.Fatal(err)
	}
	if err := receiveLiveError(t, replaced); !errors.Is(err, &StateError{}) {
		t.Fatalf("replacement survived close: %v", err)
	}
	if event := liveEvent(t, done); event.Reason != ReasonTerminated {
		t.Fatalf("close failed to terminate execution: %+v", event)
	}
	if got, err := session.Breakpoints(context.Background()); err != nil || !reflect.DeepEqual(got, old) {
		t.Fatalf("close changed committed snapshot: %+v", got)
	}
}

func TestSessionReplacementOrdersWithCompletion(t *testing.T) {
	hold := &deferredStopExecution{
		reason: vm.DebugStopCompleted, ready: make(chan *vm.DebugExecutionEvent), release: make(chan struct{}),
	}
	session, gate := newLiveSession(t, hold)
	done := continueLive(t, session, t.Context())
	for range 3 {
		gate.visit(t)
		gate.proceed(t)
	}
	select {
	case <-hold.ready:
	case <-time.After(3 * time.Second):
		t.Fatal("VM did not finish")
	}

	// VM return precedes native terminal commitment; this ordering may publish.
	old := replaceLines(t, session, 3)
	session.breakpoints.write <- struct{}{}
	defer func() { <-session.breakpoints.write }()
	waiting := &breakpointWaitContext{Context: t.Context(), waiting: make(chan struct{})}
	replaced := make(chan error, 1)
	go func() {
		_, err := session.ReplaceBreakpoints(waiting, "", lineRequests(4))
		replaced <- err
	}()
	waitForSignal(t, waiting.waiting, "replacement admission")
	close(hold.release)
	if event := liveEvent(t, done); event.Reason != ReasonCompleted {
		t.Fatalf("expected completion: %+v", event)
	}
	var state *StateError
	if err := receiveLiveError(t, replaced); !errors.As(err, &state) || state.State != "completed" {
		t.Fatalf("terminal commitment did not reject replacement: %v", err)
	}
	if got, err := session.Breakpoints(context.Background()); err != nil || !reflect.DeepEqual(got, old) {
		t.Fatalf("failed replacement changed snapshot: %+v", got)
	}
}

func TestSessionReplacementAfterExecutionCancellation(t *testing.T) {
	session, gate := newLiveSession(t, nil)
	ctx, cancel := context.WithCancel(t.Context())
	done := continueLive(t, session, ctx)
	gate.visit(t)
	cancel()
	if event := liveEvent(t, done); event.Reason != ReasonTerminated || !errors.Is(event.Error, context.Canceled) {
		t.Fatalf("expected canceled execution: %+v", event)
	}
	var state *StateError
	if _, err := session.ReplaceBreakpoints(t.Context(), "", lineRequests(3)); !errors.As(err, &state) || state.State != "terminated" {
		t.Fatalf("replacement after termination: %v", err)
	}
}

func TestAddedBreakpointDoesNotStopAtPassedLocation(t *testing.T) {
	session, gate := newLiveSession(t, nil)
	replaceLines(t, session, 4)
	done := continueLive(t, session, t.Context())
	gate.visit(t)
	gate.proceed(t)
	if event := liveEvent(t, done); event.Location.Line != 4 {
		t.Fatalf("expected stop after candidate location: %+v", event)
	}

	added := replaceLines(t, session, 3)[0]
	done = continueLive(t, session, t.Context())
	gate.visit(t)
	gate.proceed(t)
	event := liveEvent(t, done)
	if event.Reason != ReasonBreakpoint || event.Location.Line != 3 || event.HitBreakpointIDs[0] != added.ID {
		t.Fatalf("expected newly added breakpoint on a later visit: %+v", event)
	}
	if value, err := session.Evaluate(t.Context(), "i"); err != nil || value.Display != "2" {
		t.Fatalf("breakpoint was retroactive: %+v, %v", value, err)
	}
}

func TestIncrementalBreakpointMutationPublishesWhileRunning(t *testing.T) {
	session, gate := newLiveSession(t, nil)
	done := continueLive(t, session, t.Context())
	gate.visit(t)
	breakpoint, err := session.SetBreakpoint(context.Background(), source.Location{Position: source.Position{Line: 3}})
	if err != nil {
		t.Fatal(err)
	}
	gate.proceed(t)
	if event := liveEvent(t, done); event.Reason != ReasonBreakpoint || event.HitBreakpointIDs[0] != breakpoint.ID {
		t.Fatalf("incremental addition did not publish: %+v", event)
	}

	done = continueLive(t, session, t.Context())
	gate.visit(t)
	if err := session.DeleteBreakpoint(context.Background(), breakpoint.ID); err != nil {
		t.Fatal(err)
	}
	gate.proceed(t)
	gate.visit(t)
	gate.proceed(t)
	if event := liveEvent(t, done); event.Reason != ReasonCompleted {
		t.Fatalf("incremental deletion did not publish: %+v", event)
	}
}

func TestReplacementDuringRepeatedStepCommands(t *testing.T) {
	c, err := compiler.New(compiler.WithDebugInfo())
	if err != nil {
		t.Fatal(err)
	}
	src := source.New("steps.fql", "LET a = 1\nLET b = 2\nLET c = 3\nLET d = 4\nRETURN a + b + c + d")
	program, err := c.Compile(t.Context(), src)
	if err != nil {
		t.Fatal(err)
	}
	instance, err := vm.New(program)
	if err != nil {
		t.Fatal(err)
	}
	execution, err := vm.NewDebugExecution(instance, nil)
	if err != nil {
		t.Fatal(err)
	}
	hold := &deferredStopExecution{
		DebugExecution: execution, reason: vm.DebugStopBreakpoint,
		ready: make(chan *vm.DebugExecutionEvent), release: make(chan struct{}),
	}
	session := newBreakpointSession(t, src, program.Metadata.DebugPoints, hold)
	t.Cleanup(func() { _ = session.Close() })
	if _, err := session.Start(t.Context()); err != nil {
		t.Fatal(err)
	}
	active := replaceLines(t, session, 2)
	for line := 2; line <= 5; line++ {
		done := make(chan *Event, 1)
		go func() {
			event, err := session.StepIn(t.Context())
			if err != nil {
				t.Errorf("step: %v", err)
			}
			done <- event
		}()
		select {
		case <-hold.ready:
		case <-time.After(3 * time.Second):
			t.Fatal("step did not decide a stop")
		}

		committedID := active[0].ID
		if line < 5 {
			active = replaceLines(t, session, line+1)
		} else {
			replaceLines(t, session)
		}
		hold.release <- struct{}{}
		if event := liveEvent(t, done); event.Reason != ReasonBreakpoint || event.Location.Line != line ||
			!reflect.DeepEqual(event.HitBreakpointIDs, []BreakpointID{committedID}) {
			t.Fatalf("replacement changed step's committed stop: %+v", event)
		}
	}
	if event, err := session.Continue(t.Context()); err != nil || event.Reason != ReasonCompleted {
		t.Fatalf("completion after repeated step/replacement: %+v, %v", event, err)
	}
}

func lineRequests(lines ...int) []BreakpointRequest {
	requests := make([]BreakpointRequest, len(lines))
	for i, line := range lines {
		requests[i] = BreakpointRequest{Position: source.Position{Line: line}}
	}

	return requests
}

func replaceLines(t *testing.T, session *Session, lines ...int) []Breakpoint {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Second)
	defer cancel()
	results, err := session.ReplaceBreakpoints(ctx, "", lineRequests(lines...))
	if err != nil {
		t.Fatal(err)
	}
	for _, result := range results {
		if !result.Bound {
			t.Fatalf("expected bound breakpoint: %+v", result)
		}
	}

	return results
}

func continueLive(t *testing.T, session *Session, ctx context.Context) <-chan *Event {
	t.Helper()
	done := make(chan *Event, 1)
	go func() {
		event, err := session.Continue(ctx)
		if err != nil {
			t.Errorf("continue: %v", err)
		}
		done <- event
	}()

	return done
}

func liveEvent(t *testing.T, done <-chan *Event) *Event {
	t.Helper()
	select {
	case event := <-done:
		if event == nil {
			t.Fatal("execution returned no event")
		}

		return event
	case <-time.After(3 * time.Second):
		t.Fatal("execution did not report an event")

		return nil
	}
}

func receiveLiveError(t *testing.T, done <-chan error) error {
	t.Helper()
	select {
	case err := <-done:
		return err
	case <-time.After(3 * time.Second):
		t.Fatal("operation did not finish")

		return nil
	}
}
