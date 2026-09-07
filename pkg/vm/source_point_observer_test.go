package vm

import "context"

type recordingSourcePointObserver struct {
	states []sourcePointState
	action sourcePointAction
}

func (o *recordingSourcePointObserver) onSourcePoint(_ context.Context, state sourcePointState) (sourcePointAction, error) {
	o.states = append(o.states, state)

	return o.action, nil
}
