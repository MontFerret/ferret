package debugger

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func sourcePositionBefore(pos, requestedPos source.Position) bool {
	return pos.Line < requestedPos.Line || (pos.Line == requestedPos.Line && requestedPos.Column > 0 && pos.Column < requestedPos.Column)
}

func sourcePointPositionBefore(src source.Source, left, right *bytecode.DebugPoint) bool {
	leftPos := src.PositionAt(left.Span)
	rightPos := src.PositionAt(right.Span)

	return sourcePositionBefore(leftPos, rightPos)
}

func sameSourcePosition(src source.Source, left, right *bytecode.DebugPoint) bool {
	leftPos := src.PositionAt(left.Span)
	rightPos := src.PositionAt(right.Span)

	return leftPos == rightPos
}
