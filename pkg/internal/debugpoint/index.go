package debugpoint

import (
	"sort"
	"sync"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

// Index resolves ordered debug points globally or within a function.
// It must not be copied after its first lookup.
type Index struct {
	byFunction map[int][]*bytecode.DebugPoint
	points     []bytecode.DebugPoint
	functions  sync.Once
}

// New creates an index over points that are strictly ordered by PC.
func New(points []bytecode.DebugPoint) Index {
	return Index{points: points}
}

// PointForPC returns the point at pc, when one exists.
func (i *Index) PointForPC(pc int) *bytecode.DebugPoint {
	pos := sort.Search(len(i.points), func(pos int) bool {
		return i.points[pos].PC >= pc
	})

	if pos >= len(i.points) || i.points[pos].PC != pc {
		return nil
	}

	return &i.points[pos]
}

// NearestBeforeOrAt returns the nearest global point at or before pc.
func (i *Index) NearestBeforeOrAt(pc int) *bytecode.DebugPoint {
	pos := sort.Search(len(i.points), func(pos int) bool {
		return i.points[pos].PC > pc
	})

	if pos == 0 {
		return nil
	}

	return &i.points[pos-1]
}

// NearestBeforeOrAtInFunction returns the nearest point in functionID at or
// before pc.
func (i *Index) NearestBeforeOrAtInFunction(functionID, pc int) *bytecode.DebugPoint {
	i.functions.Do(i.indexFunctions)

	points := i.byFunction[functionID]
	pos := sort.Search(len(points), func(pos int) bool {
		return points[pos].PC > pc
	})

	if pos == 0 {
		return nil
	}

	return points[pos-1]
}

func (i *Index) indexFunctions() {
	i.byFunction = make(map[int][]*bytecode.DebugPoint)

	for pos := range i.points {
		point := &i.points[pos]
		i.byFunction[point.FunctionID] = append(i.byFunction[point.FunctionID], point)
	}
}
