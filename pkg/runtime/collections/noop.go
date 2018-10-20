package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type noopIterator struct{}

var NoopIterator = &noopIterator{}

func (iterator *noopIterator) HasNext() bool {
	return false
}

func (iterator *noopIterator) Next() (core.Value, core.Value, error) {
	return values.None, values.None, ErrExhausted
}
