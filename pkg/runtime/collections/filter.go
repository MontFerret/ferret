package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/pkg/errors"
)

type (
	FilterPredicate func(val core.Value, key core.Value) (bool, error)
	FilterIterator  struct {
		src       Iterator
		predicate FilterPredicate
		value     core.Value
		key       core.Value
		ready     bool
	}
)

func NewFilterIterator(src Iterator, predicate FilterPredicate) (*FilterIterator, error) {
	if core.IsNil(src) {
		return nil, errors.Wrap(core.ErrMissedArgument, "source")
	}

	if core.IsNil(predicate) {
		return nil, errors.Wrap(core.ErrMissedArgument, "predicate")
	}

	return &FilterIterator{src: src, predicate: predicate}, nil
}

func (iterator *FilterIterator) HasNext() bool {
	if !iterator.ready {
		iterator.filter()
		iterator.ready = true
	}

	return iterator.value != nil && iterator.value.Type() != core.NoneType
}

func (iterator *FilterIterator) Next() (core.Value, core.Value, error) {
	if iterator.HasNext() == true {
		val := iterator.value
		key := iterator.key

		iterator.filter()

		return val, key, nil
	}

	return values.None, values.None, ErrExhausted
}

func (iterator *FilterIterator) filter() {
	var doNext bool

	for iterator.src.HasNext() {
		val, key, err := iterator.src.Next()

		if err != nil {
			doNext = false
			break
		}

		take, err := iterator.predicate(val, key)

		if err != nil {
			doNext = false
			break
		}

		if take == true {
			doNext = true
			iterator.value = val
			iterator.key = key
			break
		}
	}

	if doNext == false {
		iterator.value = nil
		iterator.key = nil
	}
}
