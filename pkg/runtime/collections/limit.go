package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/pkg/errors"
)

type LimitIterator struct {
	src       Iterator
	count     int
	offset    int
	currCount int
}

func NewLimitIterator(src Iterator, count, offset int) (*LimitIterator, error) {
	if core.IsNil(src) {
		return nil, errors.Wrap(core.ErrMissedArgument, "source")
	}

	return &LimitIterator{src, count, offset, 0}, nil
}

func (i *LimitIterator) HasNext() bool {
	i.verifyOffset()

	if i.src.HasNext() == false {
		return false
	}

	return i.counter() < i.count
}

func (i *LimitIterator) Next() (core.Value, core.Value, error) {
	if i.counter() <= i.count {
		i.currCount++

		return i.src.Next()
	}

	return nil, nil, ErrExhausted
}

func (i *LimitIterator) counter() int {
	return i.currCount - i.offset
}

func (i *LimitIterator) verifyOffset() {
	if i.offset == 0 {
		return
	}

	if (i.offset < i.currCount) || i.src.HasNext() == false {
		return
	}

	for (i.offset > i.currCount) && i.src.HasNext() {
		i.currCount++
		i.src.Next()
	}
}
