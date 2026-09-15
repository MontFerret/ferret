package debugger

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func (s *Session) breakpointPoint(location source.Location, mode BreakpointBindingMode, check func() error) (*bytecode.DebugPoint, error) {
	exact, err := s.exactBreakpointPoint(location, check)
	if err != nil || exact != nil || mode == BreakpointBindExact {
		return exact, err
	}

	if mode == BreakpointBindNextExecutableInFunction {
		return s.nextBreakpointPointInFunction(location, check)
	}

	return s.nextBreakpointPoint(location, check)
}

func (s *Session) exactBreakpointPoint(location source.Location, check func() error) (*bytecode.DebugPoint, error) {
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

func (s *Session) nextBreakpointPoint(location source.Location, check func() error) (*bytecode.DebugPoint, error) {
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

func (s *Session) nextBreakpointPointInFunction(location source.Location, check func() error) (*bytecode.DebugPoint, error) {
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
			case previous == nil || sourcePointPositionBefore(s.source, previous, point):
				previous = point
				previousAmbiguous = false
			case sameSourcePosition(s.source, previous, point) && previous.FunctionID != point.FunctionID:
				previousAmbiguous = true
			}

			continue
		}

		switch {
		case next == nil || sourcePointPositionBefore(s.source, point, next):
			next = point
			nextAmbiguous = false
		case sameSourcePosition(s.source, next, point) && next.FunctionID != point.FunctionID:
			nextAmbiguous = true
		}
	}

	if previous == nil || next == nil || previousAmbiguous || nextAmbiguous || previous.FunctionID != next.FunctionID {
		return nil, nil
	}

	return next, nil
}

func (s *Session) debugPointLess(left, right *bytecode.DebugPoint) bool {
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
