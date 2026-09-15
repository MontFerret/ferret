package debugger

import (
	"context"

	apidebugger "github.com/MontFerret/api/debugger"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

// ReplaceBreakpoints atomically replaces one source's complete breakpoint set.
// It is valid before execution, while paused, and while running. Empty requests
// clear the source; an empty sourceName selects the launched source. Results
// follow request order. Unbound results are successful verification outcomes;
// invalid positions or binding modes fail the whole operation.
//
// Unchanged requests retain IDs, including duplicate requests matched in ID
// order. Removed IDs are never reused. Other sources are unaffected. A stop
// already decided using a previous snapshot retains its hit IDs and remains
// inspectable; additions only affect subsequent source-point checks.
//
// Cancellation or native termination before publication leaves the set intact.
// Once publication commits, the operation returns success even if cancellation
// follows. Concurrent writers are ordered by admission to the writer gate.
func (s *Session) ReplaceBreakpoints(ctx context.Context, sourceName string, requests []BreakpointRequest) ([]Breakpoint, error) {
	if err := s.lockBreakpoints(ctx, true); err != nil {
		return nil, err
	}
	defer func() { <-s.breakpointWrite }()

	check := func() error { return s.checkBreakpointContext(ctx, true) }
	if sourceName == "" {
		sourceName = s.source.Name()
	}

	previous := s.breakpointData.Load()
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

	snapshot, err := newBreakpointSnapshot(&s.pointIndex, records, nextID, check)
	if err != nil {
		return nil, err
	}

	if err := s.publishBreakpoints(ctx, snapshot, true); err != nil {
		return nil, err
	}

	return results, nil
}

// SetBreakpoint adds a location breakpoint using next-executable-in-source binding.
// It is safe while execution is running.
func (s *Session) SetBreakpoint(location source.Location) (Breakpoint, error) {
	return s.SetBreakpointAt(location, BreakpointOptions{BindingMode: BreakpointBindNextExecutableInSource})
}

// SetBreakpointAt adds a breakpoint at an explicit source location, including
// while execution is running. Each addition receives a distinct identity.
func (s *Session) SetBreakpointAt(location source.Location, opts BreakpointOptions) (Breakpoint, error) {
	ctx := context.Background()
	if err := s.lockBreakpoints(ctx, false); err != nil {
		return Breakpoint{}, err
	}
	defer func() { <-s.breakpointWrite }()

	check := func() error { return s.checkBreakpointContext(ctx, false) }
	breakpoint, err := s.resolveBreakpoint(location, opts, check)
	if err != nil {
		return Breakpoint{}, err
	}

	previous := s.breakpointData.Load()
	if previous.nextID <= 0 {
		return Breakpoint{}, runtime.Error(runtime.ErrInvalidOperation, "breakpoint identities exhausted")
	}

	breakpoint.ID = previous.nextID
	records := make([]Breakpoint, len(previous.breakpoints), len(previous.breakpoints)+1)
	copy(records, previous.breakpoints)
	records = append(records, breakpoint)
	snapshot, err := newBreakpointSnapshot(&s.pointIndex, records, previous.nextID+1, check)
	if err != nil {
		return Breakpoint{}, err
	}

	if err := s.publishBreakpoints(ctx, snapshot, false); err != nil {
		return Breakpoint{}, err
	}

	return breakpoint, nil
}

// DeleteBreakpoint removes a breakpoint by ID, including while execution is
// running. A previously decided hit remains a valid stop.
func (s *Session) DeleteBreakpoint(id BreakpointID) error {
	ctx := context.Background()
	if err := s.lockBreakpoints(ctx, false); err != nil {
		return err
	}
	defer func() { <-s.breakpointWrite }()

	check := func() error { return s.checkBreakpointContext(ctx, false) }
	previous := s.breakpointData.Load()
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

	snapshot, err := newBreakpointSnapshot(&s.pointIndex, records, previous.nextID, check)
	if err != nil {
		return err
	}

	return s.publishBreakpoints(ctx, snapshot, false)
}

// Breakpoints returns a detached, ID-ordered snapshot, including while running
// or after Close. Listing never waits for execution or an in-progress replacement.
func (s *Session) Breakpoints() []Breakpoint {
	if s == nil {
		return nil
	}

	snapshot := s.breakpointData.Load()
	if snapshot == nil {
		return []Breakpoint{}
	}

	records := snapshot.breakpoints
	out := make([]Breakpoint, len(records))
	copy(out, records)

	return out
}

func (s *Session) resolveBreakpoint(location source.Location, opts BreakpointOptions, check func() error) (Breakpoint, error) {
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

func (s *Session) lockBreakpoints(ctx context.Context, replacement bool) error {
	if err := s.checkBreakpointContext(ctx, replacement); err != nil {
		return err
	}

	var done <-chan struct{}
	if replacement {
		done = s.breakpointDone
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return s.checkBreakpointContext(ctx, replacement)
	case s.breakpointWrite <- struct{}{}:
	}

	if err := s.checkBreakpointContext(ctx, replacement); err != nil {
		<-s.breakpointWrite

		return err
	}

	return nil
}

func (s *Session) checkBreakpointContext(ctx context.Context, replacement bool) error {
	if err := s.checkContext(ctx); err != nil {
		return err
	}

	if err := s.ensureOpen(); err != nil {
		return err
	}

	if replacement {
		select {
		case <-s.breakpointDone:
			s.lifecycleMu.Lock()
			state := s.breakpointState
			s.lifecycleMu.Unlock()

			return &StateError{Operation: "replace breakpoints", State: state}
		default:
		}
	}

	return nil
}

// The writer gate owns construction. Only publication and terminal commitment
// share lifecycleMu; neither resolution nor execution holds it. No path holding
// lifecycleMu may acquire the writer gate or commandMu.
func (s *Session) publishBreakpoints(ctx context.Context, snapshot *breakpointSnapshot, replacement bool) error {
	s.lifecycleMu.Lock()
	defer s.lifecycleMu.Unlock()

	if err := s.checkContext(ctx); err != nil {
		return err
	}

	if err := s.ensureOpen(); err != nil {
		return err
	}

	if replacement && s.breakpointState != "" {
		return &StateError{Operation: "replace breakpoints", State: s.breakpointState}
	}

	s.breakpointData.Store(snapshot)

	return nil
}

// Called synchronously by the VM under commandMu. Writers never touch hit IDs.
func (s *Session) hasBreakpoint(pc int) bool {
	ids := s.breakpointData.Load().byPC[pc]
	if len(ids) == 0 {
		return false
	}

	s.hitBreakpointIDs = ids

	return true
}

func (s *Session) finishBreakpoints(state string) {
	s.lifecycleMu.Lock()
	defer s.lifecycleMu.Unlock()
	s.finishBreakpointsLocked(state)
}

func (s *Session) finishBreakpointsLocked(state string) {
	if s.breakpointState == "" {
		s.breakpointState = state
		if s.breakpointDone != nil {
			close(s.breakpointDone)
		}
	}
}

// Status can acquire the VM run lock, so this is only called after an execution
// command has returned, never from a breakpoint operation.
func (s *Session) recordExecutionTerminal() {
	switch s.execution.Status() {
	case vm.DebugExecutionCompleted:
		s.finishBreakpoints("completed")
	case vm.DebugExecutionTerminated, vm.DebugExecutionClosed:
		s.finishBreakpoints("terminated")
	}
}
