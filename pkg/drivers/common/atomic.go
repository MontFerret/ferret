package common

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"sync"
)

type (
	// AtomicValueWriter represents an atomic value writer
	AtomicValueWriter func(current core.Value) (core.Value, error)

	// AtomicValue represents an atomic value
	AtomicValue struct {
		mu    sync.Mutex
		value core.Value
	}
)

func NewAtomicValue(value core.Value) *AtomicValue {
	av := new(AtomicValue)
	av.value = value

	return av
}

// Read returns an underlying value.
// @returns (Value) - Underlying value
func (av *AtomicValue) Read() core.Value {
	av.mu.Lock()
	defer av.mu.Unlock()

	return av.value
}

// Write sets a new underlying value.
// If writer fails, the operations gets terminated and an underlying value remains.
// @param (AtomicValueWriter) - Writer function that receives a current value and returns new one.
// @returns (Error) - Error if write operation failed
func (av *AtomicValue) Write(writer AtomicValueWriter) error {
	av.mu.Lock()
	defer av.mu.Unlock()

	next, err := writer(av.value)

	if err != nil {
		return err
	}

	av.value = next

	return nil
}
