package debugger

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func sourcePositionBefore(line, column, requestedLine, requestedColumn int) bool {
	return line < requestedLine || (line == requestedLine && requestedColumn > 0 && column < requestedColumn)
}

func sourcePointPositionBefore(src *source.Source, left, right *bytecode.DebugPoint) bool {
	leftLine, leftColumn := src.LocationAt(left.Span)
	rightLine, rightColumn := src.LocationAt(right.Span)

	return leftLine < rightLine || (leftLine == rightLine && leftColumn < rightColumn)
}

func sameSourcePosition(src *source.Source, left, right *bytecode.DebugPoint) bool {
	leftLine, leftColumn := src.LocationAt(left.Span)
	rightLine, rightColumn := src.LocationAt(right.Span)

	return leftLine == rightLine && leftColumn == rightColumn
}
