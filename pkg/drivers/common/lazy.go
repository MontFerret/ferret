package common

import (
	"context"
	"sync"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	// LazyValueFactory represents a value initializer
	LazyValueFactory func(ctx context.Context) (core.Value, error)

	// LazyValue represents a value with late initialization
	LazyValue struct {
		mu      sync.Mutex
		factory LazyValueFactory
		ready   bool
		value   core.Value
		err     error
	}
)

func NewLazyValue(factory LazyValueFactory) *LazyValue {
	lz := new(LazyValue)
	lz.ready = false
	lz.factory = factory
	lz.value = values.None

	return lz
}

// Ready indicates whether the value is ready.
// @returns (Boolean) - Boolean value indicating whether the value is ready.
func (lv *LazyValue) Ready() bool {
	lv.mu.Lock()
	defer lv.mu.Unlock()

	return lv.ready
}

// Read returns an underlying value.
// Not thread safe. Should not mutated.
// @returns (Value) - Underlying value if successfully loaded, otherwise error
func (lv *LazyValue) Read(ctx context.Context) (core.Value, error) {
	lv.mu.Lock()
	defer lv.mu.Unlock()

	if !lv.ready {
		lv.load(ctx)
	}

	return lv.value, lv.err
}

// Mutate safely mutates an underlying value.
// Loads a value if it's not ready.
// Thread safe.
func (lv *LazyValue) Mutate(ctx context.Context, mutator func(v core.Value, err error)) {
	lv.mu.Lock()
	defer lv.mu.Unlock()

	if !lv.ready {
		lv.load(ctx)
	}

	mutator(lv.value, lv.err)
}

// MutateIfReady safely mutates an underlying value only if it's ready.
func (lv *LazyValue) MutateIfReady(mutator func(v core.Value, err error)) {
	lv.mu.Lock()
	defer lv.mu.Unlock()

	if lv.ready {
		mutator(lv.value, lv.err)
	}
}

// Reload resets the storage and loads data.
func (lv *LazyValue) Reload(ctx context.Context) {
	lv.mu.Lock()
	defer lv.mu.Unlock()

	lv.resetInternal()
	lv.load(ctx)
}

// Reset resets the storage.
// Next call of Read will trigger the factory function again.
func (lv *LazyValue) Reset() {
	lv.mu.Lock()
	defer lv.mu.Unlock()

	lv.resetInternal()
}

func (lv *LazyValue) resetInternal() {
	lv.ready = false
	lv.value = values.None
	lv.err = nil
}

func (lv *LazyValue) load(ctx context.Context) {
	val, err := lv.factory(ctx)

	if err == nil {
		lv.value = val
		lv.err = nil
	} else {
		lv.value = values.None
		lv.err = err
	}

	lv.ready = true
}
