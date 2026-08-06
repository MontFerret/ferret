package vm

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type cancellationProbeValue struct {
	nextErr        error
	nextStarted    chan struct{}
	cancel         context.CancelFunc
	lastContext    context.Context
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
	cooperative    bool
	blockNext      bool
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

func (v *cancellationProbeValue) Get(ctx context.Context, _ runtime.Value) (runtime.Value, error) {
	v.getCalls++
	if err := v.contextErr(ctx); err != nil {
		return runtime.None, err
	}

	return runtime.NewInt(1), nil
}

func (v *cancellationProbeValue) At(ctx context.Context, _ runtime.Int) (runtime.Value, error) {
	v.getCalls++
	if err := v.contextErr(ctx); err != nil {
		return runtime.None, err
	}

	return runtime.NewInt(1), nil
}

func (v *cancellationProbeValue) Set(ctx context.Context, _, _ runtime.Value) error {
	v.setCalls++
	return v.contextErr(ctx)
}

func (v *cancellationProbeValue) SetAt(ctx context.Context, _ runtime.Int, _ runtime.Value) error {
	v.setCalls++
	return v.contextErr(ctx)
}

func (v *cancellationProbeValue) RemoveKey(ctx context.Context, _ runtime.Value) error {
	v.removeCalls++
	return v.contextErr(ctx)
}

func (v *cancellationProbeValue) RemoveAt(ctx context.Context, _ runtime.Int) (runtime.Value, error) {
	v.removeCalls++
	if err := v.contextErr(ctx); err != nil {
		return runtime.None, err
	}

	return runtime.None, nil
}

func (v *cancellationProbeValue) Append(ctx context.Context, _ runtime.Value) error {
	v.appendCalls++
	return v.contextErr(ctx)
}

func (v *cancellationProbeValue) Equal(ctx context.Context, _ runtime.Value) (bool, error) {
	v.equalCalls++
	if err := v.contextErr(ctx); err != nil {
		return false, err
	}
	if v.cancel != nil {
		v.cancel()
	}

	return true, nil
}

func (v *cancellationProbeValue) Length(ctx context.Context) (runtime.Int, error) {
	v.lengthCalls++
	if err := v.contextErr(ctx); err != nil {
		return 0, err
	}

	return 1, nil
}

func (v *cancellationProbeValue) Iterate(ctx context.Context) (runtime.Iterator, error) {
	v.iterateCalls++
	if err := v.contextErr(ctx); err != nil {
		return nil, err
	}

	return v, nil
}

func (v *cancellationProbeValue) Next(ctx context.Context) (runtime.Value, runtime.Value, error) {
	v.nextCalls++
	if v.nextStarted != nil {
		close(v.nextStarted)
	}
	if v.blockNext {
		<-ctx.Done()
		return runtime.None, runtime.None, ctx.Err()
	}
	if err := v.contextErr(ctx); err != nil {
		return runtime.None, runtime.None, err
	}
	if v.cancel != nil {
		v.cancel()
	}
	if v.nextErr != nil {
		return runtime.None, runtime.None, v.nextErr
	}

	return runtime.None, runtime.None, io.EOF
}

func (v *cancellationProbeValue) Dispatch(ctx context.Context, _ runtime.DispatchEvent) error {
	v.dispatchCalls++
	return v.contextErr(ctx)
}

func (v *cancellationProbeValue) Subscribe(ctx context.Context, _ runtime.Subscription) (runtime.Stream, error) {
	v.subscribeCalls++
	if err := v.contextErr(ctx); err != nil {
		return nil, err
	}

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

func (v *cancellationProbeValue) Query(ctx context.Context, _ runtime.Query) (runtime.List, error) {
	v.queryCalls++
	if err := v.contextErr(ctx); err != nil {
		return nil, err
	}

	return runtime.NewArray(0), nil
}

func (v *cancellationProbeValue) QueryOne(ctx context.Context, _ runtime.Query) (runtime.Value, error) {
	v.queryCalls++
	if err := v.contextErr(ctx); err != nil {
		return runtime.None, err
	}

	return runtime.None, nil
}

func (v *cancellationProbeValue) QueryCount(ctx context.Context, _ runtime.Query) (runtime.Int, error) {
	v.queryCalls++
	if err := v.contextErr(ctx); err != nil {
		return 0, err
	}

	return runtime.ZeroInt, nil
}

func (v *cancellationProbeValue) QueryExists(ctx context.Context, _ runtime.Query) (runtime.Boolean, error) {
	v.queryCalls++
	if err := v.contextErr(ctx); err != nil {
		return false, err
	}

	return runtime.False, nil
}

func (v *cancellationProbeValue) contextErr(ctx context.Context) error {
	v.lastContext = ctx
	if !v.cooperative || ctx == nil {
		return nil
	}

	return ctx.Err()
}
