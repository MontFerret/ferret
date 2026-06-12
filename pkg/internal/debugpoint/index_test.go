package debugpoint

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

func TestIndexResolvesExactGlobalAndFunctionPoints(t *testing.T) {
	points := []bytecode.DebugPoint{
		{PC: 2, FunctionID: -1},
		{PC: 5, FunctionID: 0},
		{PC: 8, FunctionID: -1},
		{PC: 11, FunctionID: 1},
		{PC: 14, FunctionID: 0},
	}
	index := New(points)

	if got := index.PointForPC(8); got != &points[2] {
		t.Fatalf("unexpected exact point: %#v", got)
	}
	if got := index.PointForPC(9); got != nil {
		t.Fatalf("expected no exact point, got %#v", got)
	}
	if got := index.NearestBeforeOrAt(12); got != &points[3] {
		t.Fatalf("unexpected global point: %#v", got)
	}
	if got := index.NearestBeforeOrAtInFunction(0, 12); got != &points[1] {
		t.Fatalf("unexpected function point: %#v", got)
	}
	if got := index.NearestBeforeOrAtInFunction(-1, 12); got != &points[2] {
		t.Fatalf("unexpected top-level point: %#v", got)
	}
}

func TestIndexDoesNotCrossFunctionBoundaries(t *testing.T) {
	points := []bytecode.DebugPoint{
		{PC: 2, FunctionID: 0},
		{PC: 4, FunctionID: 1},
		{PC: 6, FunctionID: 0},
	}
	index := New(points)

	if got := index.NearestBeforeOrAtInFunction(0, 5); got != &points[0] {
		t.Fatalf("resolved point from another function: %#v", got)
	}
	if got := index.NearestBeforeOrAtInFunction(2, 10); got != nil {
		t.Fatalf("expected no point for missing function, got %#v", got)
	}
}
