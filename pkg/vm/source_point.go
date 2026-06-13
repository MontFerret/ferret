package vm

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

type (
	sourcePointAction uint8

	sourcePointState struct {
		pc      int
		pointID bytecode.DebugPointID
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
