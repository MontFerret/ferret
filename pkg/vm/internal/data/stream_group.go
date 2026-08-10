package data

import (
	"errors"
	"hash/fnv"
	"sync"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	streamGroupArm struct {
		stream     runtime.Stream
		closeErr   error
		panicValue any
		closeOnce  sync.Once
	}

	// StreamGroupValue owns the subscriptions established for one grouped wait.
	StreamGroupValue struct {
		arms []*streamGroupArm
	}
)

// NewStreamGroupValue creates one VM-owned value that closes every subscription.
func NewStreamGroupValue(streams []runtime.Stream) runtime.Value {
	arms := make([]*streamGroupArm, len(streams))

	for idx, stream := range streams {
		arms[idx] = &streamGroupArm{stream: stream}
	}

	return &StreamGroupValue{arms: arms}
}

func (v *StreamGroupValue) Iterate(timeout runtime.Duration) IteratorState {
	return newStreamGroupIterator(v, timeout)
}

func (v *StreamGroupValue) Close() error {
	var errs []error

	for idx := range v.arms {
		if err := v.closeArm(idx); err != nil {
			errs = append(errs, err)
		}

	}

	return errors.Join(errs...)
}

func (v *StreamGroupValue) String() string {
	return "[StreamGroup]"
}

func (v *StreamGroupValue) Hash() uint64 {
	hasher := fnv.New64a()
	_, _ = hasher.Write([]byte("vm.stream.group"))

	return hasher.Sum64()
}

func (v *StreamGroupValue) Copy() runtime.Value {
	return v
}

func (v *StreamGroupValue) closeArm(idx int) error {
	if idx < 0 || idx >= len(v.arms) {
		return runtime.Error(runtime.ErrInvalidArgument, "stream group arm index is out of range")
	}

	arm := v.arms[idx]
	arm.closeOnce.Do(func() {
		if arm.stream != nil {
			arm.closeErr = arm.stream.Close()
		}
	})

	return arm.closeErr
}
