package debugger

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestSessionSerializesConcurrentCommands(t *testing.T) {
	session, execution := newBlockingSession(t)
	if _, err := session.Start(context.Background()); err != nil {
		t.Fatal(err)
	}

	continueDone := make(chan error, 1)
	go func() {
		_, err := session.Continue(context.Background())
		continueDone <- err
	}()
	waitForSignal(t, execution.resumeStarted, "first resume")

	stepDone := make(chan error, 1)
	go func() {
		_, err := session.Step(context.Background())
		stepDone <- err
	}()
	breakpointDone := make(chan error, 1)
	go func() {
		_, err := session.SetBreakpoint("", 1)
		breakpointDone <- err
	}()

	assertBlocked(t, stepDone, "step")
	assertBlocked(t, breakpointDone, "set breakpoint")
	if calls, maxActive, _ := execution.stats(); calls != 1 || maxActive != 1 {
		t.Fatalf("commands entered execution concurrently: calls=%d max=%d", calls, maxActive)
	}

	execution.release()
	waitForError(t, continueDone, "continue")
	waitForError(t, stepDone, "step")
	waitForError(t, breakpointDone, "set breakpoint")

	if calls, maxActive, _ := execution.stats(); calls != 2 || maxActive != 1 {
		t.Fatalf("unexpected serialized execution stats: calls=%d max=%d", calls, maxActive)
	}
}

func TestSessionPauseDoesNotWaitForRunningCommand(t *testing.T) {
	session, execution := newBlockingSession(t)
	if _, err := session.Start(context.Background()); err != nil {
		t.Fatal(err)
	}

	continueDone := make(chan error, 1)
	go func() {
		_, err := session.Continue(context.Background())
		continueDone <- err
	}()
	waitForSignal(t, execution.resumeStarted, "resume")

	pauseDone := make(chan error, 1)
	go func() {
		pauseDone <- session.Pause()
	}()

	waitForError(t, pauseDone, "pause")
	waitForSignal(t, execution.pauseCalled, "pause request")
	waitForError(t, continueDone, "continue")
}

func TestSessionCloseSerializesWithCommandsAndPreservesBreakpointSnapshot(t *testing.T) {
	session, execution := newBlockingSession(t)
	breakpoint, err := session.SetBreakpoint("", 1)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := session.Start(context.Background()); err != nil {
		t.Fatal(err)
	}

	continueDone := make(chan error, 1)
	go func() {
		_, err := session.Continue(context.Background())
		continueDone <- err
	}()
	waitForSignal(t, execution.resumeStarted, "resume")

	closeDone := make(chan error, 1)
	go func() {
		closeDone <- session.Close()
	}()
	assertBlocked(t, closeDone, "close")
	if _, _, calls := execution.stats(); calls != 0 {
		t.Fatalf("close entered execution while command was running: %d", calls)
	}

	if err := session.Pause(); err != nil {
		t.Fatal(err)
	}
	waitForError(t, continueDone, "continue")
	waitForError(t, closeDone, "close")

	if _, _, calls := execution.stats(); calls != 1 {
		t.Fatalf("unexpected execution close count: %d", calls)
	}
	if got := session.Breakpoints(); len(got) != 1 || got[0].ID != breakpoint.ID {
		t.Fatalf("unexpected post-close breakpoint snapshot: %#v", got)
	}
	if _, err := session.Continue(context.Background()); err == nil || !errors.Is(err, &StateError{}) {
		t.Fatalf("expected closed-state error, got %v", err)
	}
	if err := session.Close(); err != nil {
		t.Fatal(err)
	}
}

func newBlockingSession(t *testing.T) (*Session, *blockingExecution) {
	t.Helper()

	src := source.New("concurrent.fql", "RETURN 1")
	point := bytecode.DebugPoint{ID: 5, PC: 0, Span: source.Span{Start: 0, End: 8}, FunctionID: -1}
	execution := newBlockingExecution(point)
	session, err := NewSession(Config{
		Execution:   execution,
		Values:      &fakeValueAccess{inner: vm.NewDebugValueAccess()},
		Services:    &fakeSessionServices{},
		Source:      src,
		DebugPoints: []bytecode.DebugPoint{point},
	})
	if err != nil {
		t.Fatal(err)
	}

	return session, execution
}

func waitForSignal(t *testing.T, ch <-chan struct{}, operation string) {
	t.Helper()

	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Fatalf("timed out waiting for %s", operation)
	}
}

func waitForError(t *testing.T, ch <-chan error, operation string) {
	t.Helper()

	select {
	case err := <-ch:
		if err != nil {
			t.Fatalf("%s failed: %v", operation, err)
		}
	case <-time.After(time.Second):
		t.Fatalf("timed out waiting for %s", operation)
	}
}

func assertBlocked(t *testing.T, ch <-chan error, operation string) {
	t.Helper()

	select {
	case err := <-ch:
		t.Fatalf("%s completed before active command: %v", operation, err)
	case <-time.After(25 * time.Millisecond):
	}
}
