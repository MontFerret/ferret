package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/pkg/errors"
)

type (
	FilterPredicate func(set ResultSet) (bool, error)

	FilterIterator struct {
		src       Iterator
		predicate FilterPredicate
		result    ResultSet
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

	return iterator.result != nil
}

func (iterator *FilterIterator) Next() (ResultSet, error) {
	if iterator.HasNext() == true {
		result := iterator.result

		iterator.filter()

		return result, nil
	}

	return nil, ErrExhausted
}

func (iterator *FilterIterator) filter() {
	var doNext bool

	for iterator.src.HasNext() {
		result, err := iterator.src.Next()

		if err != nil {
			doNext = false
			break
		}

		take, err := iterator.predicate(result)

		if err != nil {
			doNext = false
			break
		}

		if take == true {
			doNext = true
			iterator.result = result
			break
		}
	}

	if doNext == false {
		iterator.result = nil
	}
}
