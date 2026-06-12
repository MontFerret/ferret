package vm

import "github.com/MontFerret/ferret/v2/pkg/bytecode"

type debugControl struct {
	owner       *debugExecution
	points      map[int]*bytecode.DebugPoint
	breakpoints map[int]struct{}
	startDepth  int
	skipPC      int
	skipDepth   int
	mode        DebugResumeMode
	reason      DebugStopReason
	skip        bool
	entry       bool
}

func (c *debugControl) shouldStop(pc, depth int) bool {
	point := c.points[pc]
	if point == nil {
		return false
	}
	c.owner.current = point
	if c.skip && c.skipPC == pc && c.skipDepth == depth {
		c.skip = false
		return false
	}
	if c.owner.pauseRequested.Swap(false) {
		c.reason = DebugStopPause
		return true
	}
	if c.entry {
		c.entry = false
		c.reason = DebugStopEntry
		return true
	}
	if _, ok := c.breakpoints[pc]; ok {
		c.reason = DebugStopBreakpoint
		return true
	}
	switch c.mode {
	case DebugResumeStep:
		c.reason = DebugStopStep
		return true
	case DebugResumeNext:
		if depth <= c.startDepth {
			c.reason = DebugStopStep
			return true
		}
	case DebugResumeOut:
		if depth < c.startDepth {
			c.reason = DebugStopStep
			return true
		}
	}
	return false
}
