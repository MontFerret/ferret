package vm

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type cancellationProbeValue struct {
	nextErr        error
	cancel         context.CancelFunc
	dispatchCalls  int
	closeCalls     int
	equalCalls     int
	getCalls       int
	iterateCalls   int
	lengthCalls    int
	nextCalls      int
	appendCalls    int
	queryCalls     int
	removeCalls    int
	setCalls       int
	stringCalls    int
	subscribeCalls int
}

func (v *cancellationProbeValue) String() string {
	v.stringCalls++
	return "probe"
}

func (*cancellationProbeValue) Hash() uint64 {
	return 1
}

func (v *cancellationProbeValue) Copy() runtime.Value {
	return v
}

func (*cancellationProbeValue) Type() runtime.Type {
	return runtime.TypeObject
}

func (v *cancellationProbeValue) Get(_ context.Context, _ runtime.Value) (runtime.Value, error) {
	v.getCalls++
	return runtime.NewInt(1), nil
}

func (v *cancellationProbeValue) At(_ context.Context, _ runtime.Int) (runtime.Value, error) {
	v.getCalls++
	return runtime.NewInt(1), nil
}

func (v *cancellationProbeValue) Set(_ context.Context, _, _ runtime.Value) error {
	v.setCalls++
	return nil
}

func (v *cancellationProbeValue) SetAt(_ context.Context, _ runtime.Int, _ runtime.Value) error {
	v.setCalls++
	return nil
}

func (v *cancellationProbeValue) RemoveKey(_ context.Context, _ runtime.Value) error {
	v.removeCalls++
	return nil
}

func (v *cancellationProbeValue) RemoveAt(_ context.Context, _ runtime.Int) (runtime.Value, error) {
	v.removeCalls++
	return runtime.None, nil
}

func (v *cancellationProbeValue) Append(_ context.Context, _ runtime.Value) error {
	v.appendCalls++
	return nil
}

func (v *cancellationProbeValue) Equal(_ context.Context, _ runtime.Value) (bool, error) {
	v.equalCalls++
	if v.cancel != nil {
		v.cancel()
	}

	return true, nil
}

func (v *cancellationProbeValue) Length(_ context.Context) (runtime.Int, error) {
	v.lengthCalls++
	return 1, nil
}

func (v *cancellationProbeValue) Iterate(_ context.Context) (runtime.Iterator, error) {
	v.iterateCalls++
	return v, nil
}

func (v *cancellationProbeValue) Next(_ context.Context) (runtime.Value, runtime.Value, error) {
	v.nextCalls++
	if v.cancel != nil {
		v.cancel()
	}
	if v.nextErr != nil {
		return runtime.None, runtime.None, v.nextErr
	}

	return runtime.None, runtime.None, io.EOF
}

func (v *cancellationProbeValue) Dispatch(_ context.Context, _ runtime.DispatchEvent) error {
	v.dispatchCalls++
	return nil
}

func (v *cancellationProbeValue) Subscribe(_ context.Context, _ runtime.Subscription) (runtime.Stream, error) {
	v.subscribeCalls++
	return v, nil
}

func (*cancellationProbeValue) Read(_ context.Context) <-chan runtime.Message {
	messages := make(chan runtime.Message)
	close(messages)

	return messages
}

func (v *cancellationProbeValue) Close() error {
	v.closeCalls++
	return nil
}

func (v *cancellationProbeValue) Query(_ context.Context, _ runtime.Query) (runtime.List, error) {
	v.queryCalls++
	return runtime.NewArray(0), nil
}

func (v *cancellationProbeValue) QueryOne(_ context.Context, _ runtime.Query) (runtime.Value, error) {
	v.queryCalls++
	return runtime.None, nil
}

func (v *cancellationProbeValue) QueryCount(_ context.Context, _ runtime.Query) (runtime.Int, error) {
	v.queryCalls++
	return runtime.ZeroInt, nil
}

func (v *cancellationProbeValue) QueryExists(_ context.Context, _ runtime.Query) (runtime.Boolean, error) {
	v.queryCalls++
	return runtime.False, nil
}
