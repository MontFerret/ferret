package vm

import "context"

type recordingSourcePointObserver struct {
	states []sourcePointState
}

func (o *recordingSourcePointObserver) onSourcePoint(_ context.Context, state sourcePointState) (sourcePointAction, error) {
	o.states = append(o.states, state)

	return sourcePointContinue, nil
}
