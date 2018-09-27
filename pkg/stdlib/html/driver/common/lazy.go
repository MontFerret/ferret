package common

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"sync"
)

type (
	LazyFactory func() (core.Value, error)

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

func (lv *LazyValue) Value() (core.Value, error) {
	lv.Lock()
	defer lv.Unlock()

	if !lv.ready {
		val, err := lv.factory()

		if err == nil {
			lv.value = val
			lv.err = nil
		} else {
			lv.value = values.None
			lv.err = err
		}
	}

	return lv.value, lv.err
}

func (lv *LazyValue) Reset() {
	lv.Lock()
	defer lv.Unlock()

	lv.ready = false
	lv.value = values.None
	lv.err = nil
}
