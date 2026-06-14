package debugger

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestSessionSetBreakpointAtSupportsExplicitBindingPolicies(t *testing.T) {
	src := source.New("binding.fql", "first second\n\ninside one\n\ninside two\n\nlast")
	points := []bytecode.DebugPoint{
		{ID: 1, PC: 1, Span: source.Span{Start: 0, End: 5}, FunctionID: -1},
		{ID: 2, PC: 2, Span: source.Span{Start: 6, End: 12}, FunctionID: -1},
		{ID: 3, PC: 3, Span: source.Span{Start: 14, End: 24}, FunctionID: 0},
		{ID: 4, PC: 4, Span: source.Span{Start: 26, End: 36}, FunctionID: 0},
		{ID: 5, PC: 5, Span: source.Span{Start: 38, End: 42}, FunctionID: -1},
	}
	session := newBreakpointSession(t, src, points, &fakeExecution{status: vm.DebugExecutionNew})
	defer session.Close()

	tests := []struct {
		name       string
		location   SourceLocation
		options    BreakpointOptions
		wantPoint  bytecode.DebugPointID
		wantBound  bool
		wantLine   int
		wantColumn int
	}{
		{
			name:       "default_next_in_file",
			location:   SourceLocation{Line: 2},
			wantPoint:  3,
			wantBound:  true,
			wantLine:   3,
			wantColumn: 1,
		},
		{
			name:       "exact_line_selects_first_point",
			location:   SourceLocation{Line: 1},
			options:    BreakpointOptions{BindingMode: BreakpointBindExact},
			wantPoint:  1,
			wantBound:  true,
			wantLine:   1,
			wantColumn: 1,
		},
		{
			name:       "exact_column",
			location:   SourceLocation{Line: 1, Column: 7},
			options:    BreakpointOptions{BindingMode: BreakpointBindExact},
			wantPoint:  2,
			wantBound:  true,
			wantLine:   1,
			wantColumn: 7,
		},
		{
			name:     "exact_miss_is_unbound",
			location: SourceLocation{Line: 2},
			options:  BreakpointOptions{BindingMode: BreakpointBindExact},
		},
		{
			name:       "next_in_function",
			location:   SourceLocation{Line: 4},
			options:    BreakpointOptions{BindingMode: BreakpointBindNextExecutableInFunction},
			wantPoint:  4,
			wantBound:  true,
			wantLine:   5,
			wantColumn: 1,
		},
		{
			name:     "next_in_function_does_not_enter_udf",
			location: SourceLocation{Line: 2},
			options:  BreakpointOptions{BindingMode: BreakpointBindNextExecutableInFunction},
		},
		{
			name:     "next_in_function_does_not_leave_udf",
			location: SourceLocation{Line: 6},
			options:  BreakpointOptions{BindingMode: BreakpointBindNextExecutableInFunction},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			breakpoint, err := session.SetBreakpointAt(tc.location, tc.options)
			if err != nil {
				t.Fatal(err)
			}
			if breakpoint.Bound != tc.wantBound || breakpoint.PointID != tc.wantPoint ||
				breakpoint.Line != tc.wantLine || breakpoint.Column != tc.wantColumn {
				t.Fatalf("unexpected breakpoint: %#v", breakpoint)
			}
			if breakpoint.BindingMode != tc.options.BindingMode {
				t.Fatalf("unexpected binding mode: %#v", breakpoint)
			}
			if breakpoint.File != src.Name() {
				t.Fatalf("expected default source file, got %#v", breakpoint)
			}
		})
	}

	legacy, err := session.SetBreakpoint("", 2)
	if err != nil {
		t.Fatal(err)
	}
	if !legacy.Bound || legacy.PointID != 3 || legacy.BindingMode != BreakpointBindNextExecutableInFile {
		t.Fatalf("legacy helper did not preserve next-in-file binding: %#v", legacy)
	}
}

func TestSessionSetBreakpointAtValidatesLocationAndMode(t *testing.T) {
	src := source.New("binding.fql", "RETURN 1")
	point := bytecode.DebugPoint{ID: 1, PC: 1, Span: source.Span{Start: 0, End: 8}, FunctionID: -1}
	session := newBreakpointSession(t, src, []bytecode.DebugPoint{point}, &fakeExecution{status: vm.DebugExecutionNew})
	defer session.Close()

	for _, tc := range []struct {
		name     string
		location SourceLocation
		options  BreakpointOptions
	}{
		{name: "line", location: SourceLocation{Line: 0}},
		{name: "column", location: SourceLocation{Line: 1, Column: -1}},
		{name: "mode", location: SourceLocation{Line: 1}, options: BreakpointOptions{BindingMode: BreakpointBindingMode(99)}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := session.SetBreakpointAt(tc.location, tc.options); !errors.Is(err, runtime.ErrInvalidArgument) {
				t.Fatalf("expected invalid argument, got %v", err)
			}
		})
	}
}

func TestSessionBreakpointEventIncludesAllMatchingBreakpointIDs(t *testing.T) {
	src := source.New("binding.fql", "RETURN 1")
	point := bytecode.DebugPoint{ID: 7, PC: 3, Span: source.Span{Start: 0, End: 8}, FunctionID: -1}
	execution := &fakeExecution{
		startEvent:  &vm.DebugExecutionEvent{Reason: vm.DebugStopEntry, Point: &point},
		resumeEvent: &vm.DebugExecutionEvent{Reason: vm.DebugStopBreakpoint, Point: &point},
		status:      vm.DebugExecutionNew,
	}
	session := newBreakpointSession(t, src, []bytecode.DebugPoint{point}, execution)
	defer session.Close()

	first, err := session.SetBreakpointAt(SourceLocation{Line: 1}, BreakpointOptions{BindingMode: BreakpointBindExact})
	if err != nil {
		t.Fatal(err)
	}
	second, err := session.SetBreakpointAt(SourceLocation{Line: 1}, BreakpointOptions{BindingMode: BreakpointBindExact})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := session.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	event, err := session.Continue(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if event.Reason != ReasonBreakpoint || len(event.HitBreakpointIDs) != 2 ||
		event.HitBreakpointIDs[0] != first.ID || event.HitBreakpointIDs[1] != second.ID {
		t.Fatalf("unexpected breakpoint hit identity: %#v", event)
	}
}

func newBreakpointSession(t *testing.T, src *source.Source, points []bytecode.DebugPoint, execution vm.DebugExecution) *Session {
	t.Helper()

	session, err := NewSession(Config{
		Execution:   execution,
		Values:      &fakeValueAccess{inner: vm.NewDebugValueAccess()},
		Services:    &fakeSessionServices{},
		Source:      src,
		DebugPoints: points,
	})
	if err != nil {
		t.Fatal(err)
	}

	return session
}
