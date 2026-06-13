package debugpoint

import (
	"sort"
	"sync"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

// Index resolves ordered debug points globally or within a function.
// It must not be copied after its first lookup.
type Index struct {
	byID       map[bytecode.DebugPointID]*bytecode.DebugPoint
	byFunction map[int][]*bytecode.DebugPoint
	points     []bytecode.DebugPoint
	functions  sync.Once
}

// New creates an index over points that are strictly ordered by PC.
func New(points []bytecode.DebugPoint) Index {
	byID := make(map[bytecode.DebugPointID]*bytecode.DebugPoint, len(points))

	for pos := range points {
		point := &points[pos]
		byID[point.ID] = point
	}

	return Index{points: points, byID: byID}
}

// Points returns all debug points in PC order.
func (i *Index) Points() []bytecode.DebugPoint {
	return i.points
}

// PointByID returns the point with id, when one exists.
func (i *Index) PointByID(id bytecode.DebugPointID) *bytecode.DebugPoint {
	pos := int(id)
	if pos >= 0 && pos < len(i.points) && i.points[pos].ID == id {
		return &i.points[pos]
	}

	return i.byID[id]
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

// PointsInFunction returns the debug points in functionID ordered by PC.
func (i *Index) PointsInFunction(functionID int) []*bytecode.DebugPoint {
	i.functions.Do(i.indexFunctions)

	return i.byFunction[functionID]
}

func (i *Index) indexFunctions() {
	i.byFunction = make(map[int][]*bytecode.DebugPoint)

	for pos := range i.points {
		point := &i.points[pos]
		i.byFunction[point.FunctionID] = append(i.byFunction[point.FunctionID], point)
	}
}
