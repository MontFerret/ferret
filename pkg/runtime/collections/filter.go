package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/pkg/errors"
)

type (
	FilterPredicate func(set DataSet) (bool, error)

	FilterIterator struct {
		src       Iterator
		predicate FilterPredicate
		dataSet   DataSet
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

	return iterator.dataSet != nil
}

func (iterator *FilterIterator) Next() (DataSet, error) {
	if iterator.HasNext() == true {
		ds := iterator.dataSet

		iterator.filter()

		return ds, nil
	}

	return nil, ErrExhausted
}

func (iterator *FilterIterator) filter() {
	var doNext bool

	for iterator.src.HasNext() {
		set, err := iterator.src.Next()

		if err != nil {
			doNext = false
			break
		}

		take, err := iterator.predicate(set)

		if err != nil {
			doNext = false
			break
		}

		if take == true {
			doNext = true
			iterator.dataSet = set
			break
		}
	}

	if doNext == false {
		iterator.dataSet = nil
	}
}
