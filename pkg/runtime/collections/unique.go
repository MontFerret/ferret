package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type (
	UniqueIterator struct {
		src    Iterator
		hashes map[uint64]bool
		value  core.Value
		key    core.Value
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

	if iterator.err != nil {
		return false
	}

	if !core.IsNil(iterator.value) {
		return true
	}

	return false
}

func (iterator *UniqueIterator) Next() (core.Value, core.Value, error) {
	return iterator.value, iterator.key, iterator.err
}

func (iterator *UniqueIterator) doNext() {
	// reset state
	iterator.err = nil
	iterator.value = nil
	iterator.key = nil

	// iterate over source until we find a non-unique item
	for iterator.src.HasNext() {
		val, key, err := iterator.src.Next()

		if err != nil {
			iterator.err = err

			return
		}

		h := val.Hash()

		_, exists := iterator.hashes[h]

		if exists {
			continue
		}

		iterator.hashes[h] = true
		iterator.key = key
		iterator.value = val

		return
	}
}
