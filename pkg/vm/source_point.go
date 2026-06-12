package vm

import "context"

type (
	sourcePointAction uint8

	sourcePointState struct {
		pc      int
		pointID int
		depth   int
	}

	sourcePointObserver interface {
		onSourcePoint(context.Context, sourcePointState) (sourcePointAction, error)
	}
)

const (
	sourcePointContinue sourcePointAction = iota
	sourcePointPause
	sourcePointTerminate
)
