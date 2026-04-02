package mock

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type BlockingObservable struct{}

func NewBlockingObservable() *BlockingObservable {
	return &BlockingObservable{}
}

func (o *BlockingObservable) Subscribe(ctx context.Context, subscription runtime.Subscription) (runtime.Stream, error) {
	return NewBlockingStream(), nil
}

func (o *BlockingObservable) ReadCount() int32 {
	return 0
}

func (o *BlockingObservable) String() string {
	return "blocking_observable"
}

func (o *BlockingObservable) Hash() uint64 {
	return 0
}

func (o *BlockingObservable) Copy() runtime.Value {
	return o
}
