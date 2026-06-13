package debugpoint

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestIndexResolvesExactGlobalAndFunctionPoints(t *testing.T) {
	points := []bytecode.DebugPoint{
		{ID: 13, PC: 2, FunctionID: -1},
		{ID: 3, PC: 5, FunctionID: 0},
		{ID: 21, PC: 8, FunctionID: -1},
		{ID: 8, PC: 11, FunctionID: 1},
		{ID: 5, PC: 14, FunctionID: 0},
	}
	index := mustNewIndex(t, points)

	if got := index.PointByID(8); got == nil || got.ID != 8 {
		t.Fatalf("unexpected point by id: %#v", got)
	}
	if got := index.PointByID(4); got != nil {
		t.Fatalf("expected no point for missing id, got %#v", got)
	}
	if got := index.PointForPC(8); got == nil || got.ID != 21 {
		t.Fatalf("unexpected exact point: %#v", got)
	}
	if got := index.PointForPC(9); got != nil {
		t.Fatalf("expected no exact point, got %#v", got)
	}
	if got := index.NearestBeforeOrAt(12); got == nil || got.ID != 8 {
		t.Fatalf("unexpected global point: %#v", got)
	}
	if got := index.NearestBeforeOrAtInFunction(0, 12); got == nil || got.ID != 3 {
		t.Fatalf("unexpected function point: %#v", got)
	}
	if got := index.NearestBeforeOrAtInFunction(-1, 12); got == nil || got.ID != 21 {
		t.Fatalf("unexpected top-level point: %#v", got)
	}
	if got := index.PointsInFunction(0); len(got) != 2 || got[0].ID != 3 || got[1].ID != 5 {
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
	index := mustNewIndex(t, points)

	if got := index.NearestBeforeOrAtInFunction(0, 5); got == nil || got.ID != 4 {
		t.Fatalf("resolved point from another function: %#v", got)
	}
	if got := index.NearestBeforeOrAtInFunction(2, 10); got != nil {
		t.Fatalf("expected no point for missing function, got %#v", got)
	}
}

func TestIndexSortsAndDefensivelyCopiesPoints(t *testing.T) {
	points := []bytecode.DebugPoint{
		{ID: 2, PC: 8, FunctionID: -1},
		{ID: 1, PC: 3, FunctionID: -1},
	}
	index := mustNewIndex(t, points)
	points[0].PC = 1

	if got := index.Points(); len(got) != 2 || got[0].ID != 1 || got[1].ID != 2 {
		t.Fatalf("unexpected sorted points: %#v", got)
	}
}

func TestIndexRejectsMalformedPoints(t *testing.T) {
	tests := []struct {
		name   string
		points []bytecode.DebugPoint
	}{
		{
			name: "duplicate_id",
			points: []bytecode.DebugPoint{
				{ID: 1, PC: 1, FunctionID: -1},
				{ID: 1, PC: 2, FunctionID: -1},
			},
		},
		{
			name: "duplicate_pc",
			points: []bytecode.DebugPoint{
				{ID: 1, PC: 1, FunctionID: -1},
				{ID: 2, PC: 1, FunctionID: -1},
			},
		},
		{
			name:   "negative_id",
			points: []bytecode.DebugPoint{{ID: -1, PC: 1, FunctionID: -1}},
		},
		{
			name:   "negative_pc",
			points: []bytecode.DebugPoint{{ID: 1, PC: -1, FunctionID: -1}},
		},
		{
			name:   "invalid_function",
			points: []bytecode.DebugPoint{{ID: 1, PC: 1, FunctionID: -2}},
		},
		{
			name:   "invalid_kind",
			points: []bytecode.DebugPoint{{ID: 1, PC: 1, FunctionID: -1, Kind: bytecode.DebugPointSynthetic + 1}},
		},
		{
			name:   "invalid_span",
			points: []bytecode.DebugPoint{{ID: 1, PC: 1, FunctionID: -1, Span: source.Span{Start: -1, End: -1}}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := New(tc.points); err == nil {
				t.Fatal("expected malformed debug points to be rejected")
			}
		})
	}
}

func mustNewIndex(t *testing.T, points []bytecode.DebugPoint) Index {
	t.Helper()

	index, err := New(points)
	if err != nil {
		t.Fatal(err)
	}

	return index
}
