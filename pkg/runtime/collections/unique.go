package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type (
	UniqueIterator struct {
		src    Iterator
		hashes map[uint64]bool
		result ResultSet
		err    error
	}
)

func NewUniqueIterator(src Iterator) (*UniqueIterator, error) {
	if src == nil {
		return nil, core.Error(core.ErrMissedArgument, "source")
	}

	return &UniqueIterator{
		src:    src,
		hashes: make(map[uint64]bool),
	}, nil
}

func (iterator *UniqueIterator) HasNext() bool {
	if !iterator.src.HasNext() {
		return false
	}

	iterator.doNext()

	return iterator.result != nil
}

func (iterator *UniqueIterator) Next() (ResultSet, error) {
	return iterator.result, iterator.err
}

func (iterator *UniqueIterator) doNext() {
	// reset state
	iterator.err = nil
	iterator.result = nil

	// iterate over source until we find a non-unique item
	for iterator.src.HasNext() {
		set, err := iterator.src.Next()

		if err != nil {
			iterator.err = err

			return
		}

		if len(set) == 0 {
			continue
		}

		h := set[0].Hash()

		_, exists := iterator.hashes[h]

		if exists {
			continue
		}

		iterator.hashes[h] = true
		iterator.result = set

		return
	}
}
