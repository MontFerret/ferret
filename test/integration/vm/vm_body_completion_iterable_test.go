package vm_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type bodyCompletionIterable struct {
	iterator *bodyCompletionIterator
}

func newBodyCompletionIterable() *bodyCompletionIterable {
	return &bodyCompletionIterable{iterator: &bodyCompletionIterator{}}
}

func (i *bodyCompletionIterable) String() string {
	return "bodyCompletionIterable"
}

func (i *bodyCompletionIterable) Hash() uint64 {
	return 1
}

func (i *bodyCompletionIterable) Copy() runtime.Value {
	return i
}

func (i *bodyCompletionIterable) Iterate(context.Context) (runtime.Iterator, error) {
	return i.iterator, nil
}

func (i *bodyCompletionIterable) closed() int32 {
	return i.iterator.closeCount.Load()
}
