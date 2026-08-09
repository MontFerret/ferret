package data

import (
	"context"
	"hash/fnv"
	"io"
	"sync"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	streamGroupEvent struct {
		message runtime.Message
		index   int
		closed  bool
	}

	// StreamGroupIterator fans in active arms under one deadline and reports declaration indexes as keys.
	StreamGroupIterator struct {
		key         runtime.Value
		value       runtime.Value
		group       *StreamGroupValue
		done        chan struct{}
		events      chan streamGroupEvent
		cancel      context.CancelFunc
		timer       *time.Timer
		active      []bool
		armCancels  []context.CancelFunc
		workers     sync.WaitGroup
		activeCount int
		timeout     time.Duration
		startOnce   sync.Once
		mu          sync.Mutex
		closed      bool
	}
)

// newStreamGroupIterator defers consumption and its fixed deadline until the first wait.
func newStreamGroupIterator(group *StreamGroupValue, timeout runtime.Duration) *StreamGroupIterator {
	resolvedTimeout := time.Duration(timeout)
	if resolvedTimeout == 0 {
		resolvedTimeout = runtime.DefaultStreamTimeout
	}

	active := make([]bool, len(group.arms))
	for idx := range active {
		active[idx] = true
	}

	return &StreamGroupIterator{
		group:       group,
		timeout:     resolvedTimeout,
		done:        make(chan struct{}),
		events:      make(chan streamGroupEvent, len(group.arms)),
		armCancels:  make([]context.CancelFunc, len(group.arms)),
		active:      active,
		activeCount: len(active),
		value:       runtime.None,
		key:         runtime.None,
	}
}

func (it *StreamGroupIterator) Next(ctx context.Context) error {
	it.mu.Lock()
	closed := it.closed
	it.mu.Unlock()

	if closed {
		return io.EOF
	}

	it.start(ctx)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-it.timer.C:
			return runtime.ErrTimeout
		case event, open := <-it.events:
			if !open {
				return io.EOF
			}

			if !it.isActive(event.index) {
				continue
			}

			if event.closed {
				remaining, err := it.completeArm(event.index)
				if err != nil {
					return err
				}

				if remaining == 0 {
					return io.EOF
				}

				continue
			}

			if err := event.message.Err(); err != nil {
				return err
			}

			it.mu.Lock()
			it.value = event.message.Value()
			it.key = runtime.NewInt(event.index)
			it.mu.Unlock()

			return nil
		}
	}
}

func (it *StreamGroupIterator) Value() runtime.Value {
	it.mu.Lock()
	defer it.mu.Unlock()

	return it.value
}

func (it *StreamGroupIterator) Key() runtime.Value {
	it.mu.Lock()
	defer it.mu.Unlock()

	return it.key
}

// ArmDone prevents queued messages or errors from a satisfied arm from affecting the group.
func (it *StreamGroupIterator) ArmDone(index int) error {
	if index < 0 || index >= len(it.active) {
		return runtime.Error(runtime.ErrInvalidArgument, "stream group arm index is out of range")
	}

	_, err := it.completeArm(index)

	return err
}

func (it *StreamGroupIterator) Close() error {
	it.mu.Lock()

	if it.closed {
		it.mu.Unlock()
		return nil
	}

	it.closed = true
	cancel := it.cancel
	armCancels := append([]context.CancelFunc(nil), it.armCancels...)
	timer := it.timer
	it.mu.Unlock()

	if cancel != nil {
		cancel()
	}

	for _, armCancel := range armCancels {
		if armCancel != nil {
			armCancel()
		}
	}

	if timer != nil {
		timer.Stop()
	}

	err := it.group.Close()
	if cancel != nil {
		<-it.done
	}

	return err
}

func (it *StreamGroupIterator) String() string {
	return "[StreamGroupIterator]"
}

func (it *StreamGroupIterator) Hash() uint64 {
	hasher := fnv.New64a()
	_, _ = hasher.Write([]byte("vm.stream.group.iterator"))

	return hasher.Sum64()
}

func (it *StreamGroupIterator) Copy() runtime.Value {
	return it
}

func (it *StreamGroupIterator) start(ctx context.Context) {
	it.startOnce.Do(func() {
		groupCtx, cancel := context.WithCancel(ctx)
		it.mu.Lock()
		it.cancel = cancel
		it.timer = time.NewTimer(it.timeout)
		it.mu.Unlock()

		for idx, arm := range it.group.arms {
			armCtx, armCancel := context.WithCancel(groupCtx)
			it.mu.Lock()
			it.armCancels[idx] = armCancel
			it.mu.Unlock()
			it.workers.Add(1)

			go it.consumeArm(armCtx, idx, arm.stream)
		}

		go func() {
			it.workers.Wait()
			close(it.events)
			close(it.done)
		}()
	})
}

func (it *StreamGroupIterator) consumeArm(ctx context.Context, index int, stream runtime.Stream) {
	defer it.workers.Done()

	messages := stream.Read(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case message, open := <-messages:
			event := streamGroupEvent{index: index, message: message, closed: !open}
			select {
			case it.events <- event:
			case <-ctx.Done():
				return
			}
			if !open {
				return
			}
		}
	}
}

func (it *StreamGroupIterator) isActive(index int) bool {
	it.mu.Lock()
	defer it.mu.Unlock()

	return index >= 0 && index < len(it.active) && it.active[index]
}

func (it *StreamGroupIterator) completeArm(index int) (int, error) {
	it.mu.Lock()

	if !it.active[index] {
		remaining := it.activeCount
		it.mu.Unlock()

		return remaining, nil
	}

	it.active[index] = false
	it.activeCount--
	remaining := it.activeCount
	cancel := it.armCancels[index]
	it.mu.Unlock()

	if cancel != nil {
		cancel()
	}

	return remaining, it.group.closeArm(index)
}
