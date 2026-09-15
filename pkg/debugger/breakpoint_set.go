package debugger

import (
	"context"
	"sync/atomic"

	apidebugger "github.com/MontFerret/api/debugger"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/internal/debugpoint"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

// breakpointSet owns immutable published records and serialized construction.
// The lifecycle command lock protects only captured hit IDs; writers never
// touch them, and listing reads the published snapshot without either lock.
type breakpointSet struct {
	lifecycle  *sessionLifecycle
	pointIndex *debugpoint.Index
	source     *source.Source
	write      chan struct{}
	data       atomic.Pointer[breakpointSnapshot]
	predicate  vm.DebugBreakpointPredicate
	hitIDs     []BreakpointID
}

// Initialize in place so the predicate and atomic snapshot retain one owner.
func (s *breakpointSet) initialize(lifecycle *sessionLifecycle, src *source.Source, points *debugpoint.Index) {
	s.lifecycle = lifecycle
	s.source = src
	s.pointIndex = points
	s.write = make(chan struct{}, 1)
	s.data.Store(&breakpointSnapshot{nextID: 1})
	s.predicate = s.hasBreakpoint
}

func (s *breakpointSet) replace(ctx context.Context, sourceName string, requests []BreakpointRequest) ([]Breakpoint, error) {
	if err := s.lockBreakpoints(ctx, true); err != nil {
		return nil, err
	}
	defer func() { <-s.write }()

	check := func() error { return s.lifecycle.checkBreakpointContext(ctx, true) }
	if sourceName == "" {
		sourceName = s.source.Name()
	}

	previous := s.data.Load()
	available := make(map[BreakpointRequest][]Breakpoint)
	records := make([]Breakpoint, 0, len(previous.breakpoints)+len(requests))
	for _, breakpoint := range previous.breakpoints {
		if err := check(); err != nil {
			return nil, err
		}

		if breakpoint.RequestedLocation.SourceName != sourceName {
			records = append(records, breakpoint)

			continue
		}

		request := BreakpointRequest{
			Position: breakpoint.RequestedLocation.Position,
			Options:  BreakpointOptions{BindingMode: breakpoint.BindingMode},
		}
		available[request] = append(available[request], breakpoint)
	}

	nextID := previous.nextID
	results := make([]Breakpoint, 0, len(requests))
	for _, request := range requests {
		if err := check(); err != nil {
			return nil, err
		}

		location := source.Location{SourceName: sourceName, Position: request.Position}
		breakpoint, err := s.resolveBreakpoint(location, request.Options, check)
		if err != nil {
			return nil, err
		}

		if matches := available[request]; len(matches) != 0 {
			breakpoint.ID = matches[0].ID
			available[request] = matches[1:]
		} else {
			if nextID <= 0 {
				return nil, runtime.Error(runtime.ErrInvalidOperation, "breakpoint identities exhausted")
			}

			breakpoint.ID = nextID
			nextID++
		}

		results = append(results, breakpoint)
		records = append(records, breakpoint)
	}

	snapshot, err := newBreakpointSnapshot(s.pointIndex, records, nextID, check)
	if err != nil {
		return nil, err
	}

	if err := s.publishBreakpoints(ctx, snapshot, true); err != nil {
		return nil, err
	}

	return results, nil
}

func (s *breakpointSet) add(location source.Location, opts BreakpointOptions) (Breakpoint, error) {
	ctx := context.Background()
	if err := s.lockBreakpoints(ctx, false); err != nil {
		return Breakpoint{}, err
	}
	defer func() { <-s.write }()

	check := func() error { return s.lifecycle.checkBreakpointContext(ctx, false) }
	breakpoint, err := s.resolveBreakpoint(location, opts, check)
	if err != nil {
		return Breakpoint{}, err
	}

	previous := s.data.Load()
	if previous.nextID <= 0 {
		return Breakpoint{}, runtime.Error(runtime.ErrInvalidOperation, "breakpoint identities exhausted")
	}

	breakpoint.ID = previous.nextID
	records := make([]Breakpoint, len(previous.breakpoints), len(previous.breakpoints)+1)
	copy(records, previous.breakpoints)
	records = append(records, breakpoint)
	snapshot, err := newBreakpointSnapshot(s.pointIndex, records, previous.nextID+1, check)
	if err != nil {
		return Breakpoint{}, err
	}

	if err := s.publishBreakpoints(ctx, snapshot, false); err != nil {
		return Breakpoint{}, err
	}

	return breakpoint, nil
}

func (s *breakpointSet) delete(id BreakpointID) error {
	ctx := context.Background()
	if err := s.lockBreakpoints(ctx, false); err != nil {
		return err
	}
	defer func() { <-s.write }()

	check := func() error { return s.lifecycle.checkBreakpointContext(ctx, false) }
	previous := s.data.Load()
	records := make([]Breakpoint, 0, len(previous.breakpoints))
	found := false
	for _, breakpoint := range previous.breakpoints {
		if err := check(); err != nil {
			return err
		}

		if breakpoint.ID == id {
			found = true

			continue
		}

		records = append(records, breakpoint)
	}

	if !found {
		return runtime.Errorf(runtime.ErrNotFound, "breakpoint %d", id)
	}

	snapshot, err := newBreakpointSnapshot(s.pointIndex, records, previous.nextID, check)
	if err != nil {
		return err
	}

	return s.publishBreakpoints(ctx, snapshot, false)
}

func (s *breakpointSet) list() []Breakpoint {
	snapshot := s.data.Load()
	if snapshot == nil {
		return []Breakpoint{}
	}

	records := snapshot.breakpoints
	out := make([]Breakpoint, len(records))
	copy(out, records)

	return out
}

func (s *breakpointSet) resolveBreakpoint(location source.Location, opts BreakpointOptions, check func() error) (Breakpoint, error) {
	if location.Line <= 0 {
		return Breakpoint{}, runtime.Error(runtime.ErrInvalidArgument, "breakpoint line must be positive")
	}

	if location.Column < 0 {
		return Breakpoint{}, runtime.Error(runtime.ErrInvalidArgument, "breakpoint column must not be negative")
	}

	if opts.BindingMode < BreakpointBindNextExecutableInSource || opts.BindingMode > BreakpointBindNextExecutableInFunction {
		return Breakpoint{}, runtime.Errorf(runtime.ErrInvalidArgument, "unknown breakpoint binding mode %d", opts.BindingMode)
	}

	if location.SourceName == "" {
		location.SourceName = s.source.Name()
	}

	breakpoint := Breakpoint{RequestedLocation: location, BindingMode: opts.BindingMode}
	if location.SourceName != s.source.Name() {
		return breakpoint, nil
	}

	point, err := s.breakpointPoint(location, opts.BindingMode, check)
	if err != nil {
		return Breakpoint{}, err
	}

	if point != nil {
		breakpoint.Bound = true
		breakpoint.PointID = apidebugger.PointID(point.ID)
		breakpoint.FunctionID = apidebugger.FunctionID(point.FunctionID)
		breakpoint.Location = s.source.RangeAt(point.Span)
	}

	return breakpoint, nil
}

func (s *breakpointSet) lockBreakpoints(ctx context.Context, replacement bool) error {
	if err := s.lifecycle.checkBreakpointContext(ctx, replacement); err != nil {
		return err
	}

	var done <-chan struct{}
	if replacement {
		done = s.lifecycle.terminalDone
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return s.lifecycle.checkBreakpointContext(ctx, replacement)
	case s.write <- struct{}{}:
	}

	if err := s.lifecycle.checkBreakpointContext(ctx, replacement); err != nil {
		<-s.write

		return err
	}

	return nil
}

// Construction holds the writer gate; only the atomic commit shares the
// lifecycle lock with terminal commitment and closure.
func (s *breakpointSet) publishBreakpoints(ctx context.Context, snapshot *breakpointSnapshot, replacement bool) error {
	if err := s.lifecycle.lockPublication(ctx, replacement); err != nil {
		return err
	}
	defer s.lifecycle.unlockPublication()

	s.data.Store(snapshot)

	return nil
}

// Called synchronously by the VM under the lifecycle command lock.
func (s *breakpointSet) hasBreakpoint(pc int) bool {
	ids := s.data.Load().byPC[pc]
	if len(ids) == 0 {
		return false
	}

	s.hitIDs = ids

	return true
}

func (s *breakpointSet) resetHit() {
	s.hitIDs = nil
}

func (s *breakpointSet) hitBreakpointIDs() []BreakpointID {
	return append([]BreakpointID(nil), s.hitIDs...)
}

func (s *breakpointSet) breakpointPoint(location source.Location, mode BreakpointBindingMode, check func() error) (*bytecode.DebugPoint, error) {
	exact, err := s.exactBreakpointPoint(location, check)
	if err != nil || exact != nil || mode == BreakpointBindExact {
		return exact, err
	}

	if mode == BreakpointBindNextExecutableInFunction {
		return s.nextBreakpointPointInFunction(location, check)
	}

	return s.nextBreakpointPoint(location, check)
}

func (s *breakpointSet) exactBreakpointPoint(location source.Location, check func() error) (*bytecode.DebugPoint, error) {
	var best *bytecode.DebugPoint
	points := s.pointIndex.Points()

	for i := range points {
		if err := check(); err != nil {
			return nil, err
		}

		point := &points[i]
		pos := s.source.PositionAt(point.Span)

		if pos.Line != location.Line || (location.Column > 0 && pos.Column != location.Column) {
			continue
		}

		if best == nil || s.debugPointLess(point, best) {
			best = point
		}
	}

	return best, nil
}

func (s *breakpointSet) nextBreakpointPoint(location source.Location, check func() error) (*bytecode.DebugPoint, error) {
	var best *bytecode.DebugPoint
	points := s.pointIndex.Points()

	for i := range points {
		if err := check(); err != nil {
			return nil, err
		}

		point := &points[i]
		pos := s.source.PositionAt(point.Span)

		if sourcePositionBefore(pos, location.Position) {
			continue
		}

		if best == nil || s.debugPointLess(point, best) {
			best = point
		}
	}

	return best, nil
}

func (s *breakpointSet) nextBreakpointPointInFunction(location source.Location, check func() error) (*bytecode.DebugPoint, error) {
	var previous, next *bytecode.DebugPoint
	previousAmbiguous := false
	nextAmbiguous := false
	points := s.pointIndex.Points()

	for i := range points {
		if err := check(); err != nil {
			return nil, err
		}

		point := &points[i]
		pos := s.source.PositionAt(point.Span)

		if sourcePositionBefore(pos, location.Position) {
			switch {
			case previous == nil || sourcePointPositionBefore(*s.source, previous, point):
				previous = point
				previousAmbiguous = false
			case sameSourcePosition(*s.source, previous, point) && previous.FunctionID != point.FunctionID:
				previousAmbiguous = true
			}

			continue
		}

		switch {
		case next == nil || sourcePointPositionBefore(*s.source, point, next):
			next = point
			nextAmbiguous = false
		case sameSourcePosition(*s.source, next, point) && next.FunctionID != point.FunctionID:
			nextAmbiguous = true
		}
	}

	if previous == nil || next == nil || previousAmbiguous || nextAmbiguous || previous.FunctionID != next.FunctionID {
		return nil, nil
	}

	return next, nil
}

func (s *breakpointSet) debugPointLess(left, right *bytecode.DebugPoint) bool {
	leftPos := s.source.PositionAt(left.Span)
	rightPos := s.source.PositionAt(right.Span)

	if leftPos.Line != rightPos.Line || leftPos.Column != rightPos.Column {
		return sourcePositionBefore(leftPos, rightPos)
	}

	if left.PC != right.PC {
		return left.PC < right.PC
	}

	return left.ID < right.ID
}
