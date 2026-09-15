package debugger

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestReplacementResolutionIdentityAndIsolation(t *testing.T) {
	session := replacementSession(t)
	first, err := session.ReplaceBreakpoints(t.Context(), "", lineRequests(1, 3, 1, 99))
	if err != nil {
		t.Fatal(err)
	}
	if len(first) != 4 || !first[0].Bound || !first[1].Bound || first[2].ID == first[0].ID || first[3].Bound {
		t.Fatalf("unexpected initial results: %+v", first)
	}

	other, err := session.ReplaceBreakpoints(t.Context(), "other.fql", lineRequests(1))
	if err != nil || len(other) != 1 || other[0].Bound {
		t.Fatalf("unrelated source: %+v, %v", other, err)
	}

	second, err := session.ReplaceBreakpoints(t.Context(), "replace.fql", lineRequests(3, 1, 1, 1, 2))
	if err != nil {
		t.Fatal(err)
	}
	if second[0].ID != first[1].ID || second[1].ID != first[0].ID || second[2].ID != first[2].ID ||
		second[3].ID <= other[0].ID || second[4].Location.Line != 3 {
		t.Fatalf("identity or relocation changed: first=%+v second=%+v", first, second)
	}
	if got := session.Breakpoints(); len(got) != 6 {
		t.Fatalf("source replacement disturbed other source: %+v", got)
	}

	second[0].ID = -1
	snapshot := session.Breakpoints()
	for i := 1; i < len(snapshot); i++ {
		if snapshot[i-1].ID >= snapshot[i].ID {
			t.Fatalf("snapshot is not ordered: %+v", snapshot)
		}
	}
	snapshot[0].RequestedLocation.Line = 99
	if got := session.Breakpoints()[0]; got.RequestedLocation.Line != 1 {
		t.Fatal("caller modified immutable snapshot")
	}

	if _, err := session.ReplaceBreakpoints(t.Context(), "", nil); err != nil {
		t.Fatal(err)
	}
	if got := session.Breakpoints(); !reflect.DeepEqual(got, other) {
		t.Fatalf("clear disturbed another source: %+v", got)
	}

	readded, err := session.ReplaceBreakpoints(t.Context(), "", lineRequests(1))
	if err != nil || readded[0].ID <= second[4].ID {
		t.Fatalf("removed identity reused: %+v, %v", readded, err)
	}
}

func TestEmptyBreakpointSnapshotSurvivesZeroValueClose(t *testing.T) {
	session := &Session{}
	if got := session.Breakpoints(); got == nil || len(got) != 0 {
		t.Fatalf("expected detached empty snapshot: %+v", got)
	}
	if err := session.Close(); err != nil {
		t.Fatal(err)
	}
	if got := session.Breakpoints(); got == nil || len(got) != 0 {
		t.Fatalf("close changed empty snapshot: %+v", got)
	}
}

func TestReplacementPreservesMixedBindingPolicies(t *testing.T) {
	session := replacementSession(t)
	requests := []BreakpointRequest{
		{Position: source.Position{Line: 1, Column: 1}, Options: BreakpointOptions{BindingMode: BreakpointBindExact}},
		{Position: source.Position{Line: 2}, Options: BreakpointOptions{BindingMode: BreakpointBindExact}},
		{Position: source.Position{Line: 2}},
		{Position: source.Position{Line: 2}, Options: BreakpointOptions{BindingMode: BreakpointBindNextExecutableInFunction}},
	}
	got, err := session.ReplaceBreakpoints(t.Context(), "", requests)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 4 || !got[0].Bound || got[1].Bound || !got[2].Bound || !got[3].Bound ||
		got[2].Location.Line != 3 || got[3].Location.Line != 3 {
		t.Fatalf("binding policies changed: %+v", got)
	}
	for i := range requests {
		if got[i].RequestedLocation.Position != requests[i].Position || got[i].BindingMode != requests[i].Options.BindingMode {
			t.Fatalf("request %d lost its policy: %+v", i, got[i])
		}
	}
}

func TestReplacementFailureDoesNotPublishOrConsumeIDs(t *testing.T) {
	for _, invalid := range []BreakpointRequest{
		{Position: source.Position{Line: 0}},
		{Position: source.Position{Line: 1, Column: -1}},
		{Position: source.Position{Line: 1}, Options: BreakpointOptions{BindingMode: 99}},
	} {
		session := replacementSession(t)
		old := replaceLines(t, session, 1)
		requests := append(lineRequests(3), invalid)
		if _, err := session.ReplaceBreakpoints(t.Context(), "", requests); !errors.Is(err, runtime.ErrInvalidArgument) {
			t.Fatalf("expected validation failure: %v", err)
		}
		if got := session.Breakpoints(); !reflect.DeepEqual(got, old) {
			t.Fatalf("failed replacement partially published: %+v", got)
		}
		if got := replaceLines(t, session, 3); got[0].ID != old[0].ID+1 {
			t.Fatalf("failure consumed identity: %+v", got)
		}
	}
}

func TestReplacementCancellationDoesNotCancelExecutionOrCommittedState(t *testing.T) {
	session := replacementSession(t)
	ctx, cancel := context.WithCancel(t.Context())
	old, err := session.ReplaceBreakpoints(ctx, "", lineRequests(1))
	if err != nil {
		t.Fatal(err)
	}
	cancel()
	if _, err := session.ReplaceBreakpoints(ctx, "", lineRequests(3)); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled replacement: %v", err)
	}
	if _, err := session.ReplaceBreakpoints(nil, "", nil); !errors.Is(err, runtime.ErrInvalidArgument) {
		t.Fatalf("nil context: %v", err)
	}
	if got := session.Breakpoints(); !reflect.DeepEqual(got, old) {
		t.Fatalf("cancellation changed committed state: %+v", got)
	}
	// Incremental methods retain their pre-existing post-completion admission.
	session.lifecycle.finishTerminal("completed")
	added, err := session.SetBreakpoint(source.Location{Position: source.Position{Line: 3}})
	if err != nil {
		t.Fatal(err)
	}
	if err := session.DeleteBreakpoint(added.ID); err != nil {
		t.Fatal(err)
	}
}

func TestBreakpointSnapshotReadsAreAtomic(t *testing.T) {
	session := replacementSession(t)
	replaceLines(t, session, 1, 1)
	var readers sync.WaitGroup
	done := make(chan struct{})
	for range 4 {
		readers.Add(1)
		go func() {
			defer readers.Done()
			for {
				select {
				case <-done:
					return
				default:
				}
				got := session.Breakpoints()
				if len(got) != 2 || got[0].RequestedLocation.Line != got[1].RequestedLocation.Line || got[0].ID >= got[1].ID {
					t.Errorf("partial snapshot: %+v", got)

					return
				}
			}
		}()
	}
	for i := range 100 {
		replaceLines(t, session, 1+2*(i%2), 1+2*(i%2))
	}
	close(done)
	readers.Wait()
}

func TestConcurrentSourceReplacementsPreserveEachOther(t *testing.T) {
	session := replacementSession(t)
	var writers sync.WaitGroup
	for _, name := range []string{"a", "b", "c", "d"} {
		writers.Add(1)
		go func() {
			defer writers.Done()
			for range 20 {
				if _, err := session.ReplaceBreakpoints(t.Context(), name, lineRequests(1, 3)); err != nil {
					t.Errorf("source %s: %v", name, err)

					return
				}
			}
		}()
	}
	writers.Wait()
	got := session.Breakpoints()
	counts := make(map[string]int)
	for _, breakpoint := range got {
		counts[breakpoint.RequestedLocation.SourceName]++
		if breakpoint.Bound {
			t.Fatalf("unrelated source became executable: %+v", breakpoint)
		}
	}
	if !reflect.DeepEqual(counts, map[string]int{"a": 2, "b": 2, "c": 2, "d": 2}) {
		t.Fatalf("concurrent replacement lost source state: %+v", got)
	}
}

func replacementSession(t *testing.T) *Session {
	t.Helper()
	src := source.New("replace.fql", "RETURN 1\n\nRETURN 2")
	points := []bytecode.DebugPoint{
		{ID: 1, PC: 1, Span: source.Span{Start: 0, End: 8}, FunctionID: bytecode.NoFunction},
		{ID: 2, PC: 3, Span: source.Span{Start: 10, End: 18}, FunctionID: bytecode.NoFunction},
	}
	session := newBreakpointSession(t, src, points, &fakeExecution{status: vm.DebugExecutionNew})
	t.Cleanup(func() { _ = session.Close() })

	return session
}
