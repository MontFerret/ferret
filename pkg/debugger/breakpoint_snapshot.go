package debugger

import (
	"cmp"
	"slices"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/internal/debugpoint"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// breakpointSnapshot owns its records and indexes. They are immutable after
// publication, including the ID slices retained by an already-decided hit.
type breakpointSnapshot struct {
	byPC        map[int][]BreakpointID
	breakpoints []Breakpoint
	nextID      BreakpointID
}

func newBreakpointSnapshot(points *debugpoint.Index, records []Breakpoint, nextID BreakpointID, check func() error) (*breakpointSnapshot, error) {
	slices.SortFunc(records, func(a, b Breakpoint) int { return cmp.Compare(a.ID, b.ID) })
	snapshot := &breakpointSnapshot{breakpoints: records, nextID: nextID}

	for _, breakpoint := range records {
		if err := check(); err != nil {
			return nil, err
		}

		if !breakpoint.Bound {
			continue
		}

		point := points.PointByID(bytecode.DebugPointID(breakpoint.PointID))
		if point == nil {
			return nil, runtime.Error(runtime.ErrUnexpected, "bound breakpoint has no debug point")
		}

		if snapshot.byPC == nil {
			snapshot.byPC = make(map[int][]BreakpointID)
		}

		snapshot.byPC[point.PC] = append(snapshot.byPC[point.PC], breakpoint.ID)
	}

	return snapshot, nil
}
