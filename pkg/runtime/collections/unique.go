package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type (
	UniqueIterator struct {
		src     Iterator
		hashes  map[uint64]bool
		hashKey string
		dataSet DataSet
		err     error
	}
)

func NewUniqueIterator(src Iterator, hashKey string) (*UniqueIterator, error) {
	if src == nil {
		return nil, core.Error(core.ErrMissedArgument, "source")
	}

	return &UniqueIterator{
		src:     src,
		hashes:  make(map[uint64]bool),
		hashKey: hashKey,
	}, nil
}

func (iterator *UniqueIterator) HasNext() bool {
	if !iterator.src.HasNext() {
		return false
	}

	iterator.doNext()

	if iterator.err != nil {
		return false
	}

	if iterator.dataSet != nil {
		return true
	}

	return false
}

func (iterator *UniqueIterator) Next() (DataSet, error) {
	return iterator.dataSet, iterator.err
}

func (iterator *UniqueIterator) doNext() {
	// reset state
	iterator.err = nil
	iterator.dataSet = nil

	// iterate over source until we find a non-unique item
	for iterator.src.HasNext() {
		ds, err := iterator.src.Next()

		if err != nil {
			iterator.err = err

			return
		}

		h := ds.Get(iterator.hashKey).Hash()

		_, exists := iterator.hashes[h]

		if exists {
			continue
		}

		iterator.hashes[h] = true
		iterator.dataSet = ds

		return
	}
}
