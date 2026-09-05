package vm

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type debugControl struct {
	owner       *debugExecution
	breakpoints map[int]struct{}
	startDepth  int
	skipPC      int
	skipDepth   int
	mode        DebugResumeMode
	reason      DebugStopReason
	skip        bool
	entry       bool
}

func (c *debugControl) onSourcePoint(ctx context.Context, state sourcePointState) (sourcePointAction, error) {
	if ctx != nil {
		if err := ctx.Err(); err != nil {
			return sourcePointTerminate, err
		}
	}

	point := c.owner.points.PointByID(state.pointID)
	if point == nil || point.PC != state.pc {
		return sourcePointTerminate, runtime.Errorf(runtime.ErrUnexpected, "source point id %d does not match pc %d", state.pointID, state.pc)
	}

	c.owner.current = point

	if c.shouldStop(state.pc, state.depth) {
		return sourcePointPause, nil
	}

	return sourcePointContinue, nil
}

func (c *debugControl) shouldStop(pc, depth int) bool {
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
	case DebugResumeStepIn:
		c.reason = DebugStopStep

		return true
	case DebugResumeStepOver:
		if depth <= c.startDepth {
			c.reason = DebugStopStep

			return true
		}
	case DebugResumeStepOut:
		if depth < c.startDepth {
			c.reason = DebugStopStep

			return true
		}
	default:
		return false
	}

	return false
}
