package debugpoint

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

func TestIndexResolvesExactGlobalAndFunctionPoints(t *testing.T) {
	points := []bytecode.DebugPoint{
		{ID: 13, PC: 2, FunctionID: -1},
		{ID: 3, PC: 5, FunctionID: 0},
		{ID: 21, PC: 8, FunctionID: -1},
		{ID: 8, PC: 11, FunctionID: 1},
		{ID: 5, PC: 14, FunctionID: 0},
	}
	index := New(points)

	if got := index.PointByID(8); got != &points[3] {
		t.Fatalf("unexpected point by id: %#v", got)
	}
	if got := index.PointByID(4); got != nil {
		t.Fatalf("expected no point for missing id, got %#v", got)
	}
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
	if got := index.PointsInFunction(0); len(got) != 2 || got[0] != &points[1] || got[1] != &points[4] {
		t.Fatalf("unexpected function points: %#v", got)
	}
	if got := index.Points(); len(got) != len(points) || got[0].ID != points[0].ID {
		t.Fatalf("unexpected ordered points: %#v", got)
	}
}

func TestIndexDoesNotCrossFunctionBoundaries(t *testing.T) {
	points := []bytecode.DebugPoint{
		{ID: 4, PC: 2, FunctionID: 0},
		{ID: 9, PC: 4, FunctionID: 1},
		{ID: 2, PC: 6, FunctionID: 0},
	}
	index := New(points)

	if got := index.NearestBeforeOrAtInFunction(0, 5); got != &points[0] {
		t.Fatalf("resolved point from another function: %#v", got)
	}
	if got := index.NearestBeforeOrAtInFunction(2, 10); got != nil {
		t.Fatalf("expected no point for missing function, got %#v", got)
	}
}
