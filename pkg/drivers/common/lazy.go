package common

import (
	"context"
	"sync"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	LazyFactory func(ctx context.Context) (core.Value, error)

	LazyValue struct {
		sync.Mutex
		factory LazyFactory
		ready   bool
		value   core.Value
		err     error
	}
)

func NewLazyValue(factory LazyFactory) *LazyValue {
	lz := new(LazyValue)
	lz.ready = false
	lz.factory = factory
	lz.value = values.None

	return lz
}

// Ready indicates whether the value is ready.
// @returns (Boolean) - Boolean value indicating whether the value is ready.
func (lv *LazyValue) Ready() bool {
	lv.Lock()
	defer lv.Unlock()

	return lv.ready
}

// Read returns an underlying value.
// Not thread safe. Should not mutated.
// @returns (Value) - Underlying value if successfully loaded, otherwise error
func (lv *LazyValue) Read(ctx context.Context) (core.Value, error) {
	lv.Lock()
	defer lv.Unlock()

	if !lv.ready {
		lv.load(ctx)
	}

	return lv.value, lv.err
}

// Write safely mutates an underlying value.
// Loads a value if it's not ready.
// Thread safe.
func (lv *LazyValue) Write(ctx context.Context, writer func(v core.Value, err error)) {
	lv.Lock()
	defer lv.Unlock()

	if !lv.ready {
		lv.load(ctx)
	}

	writer(lv.value, lv.err)
}

// Reset resets the storage.
// Next call of Read will trigger the factory function again.
func (lv *LazyValue) Reset() {
	lv.Lock()
	defer lv.Unlock()

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
