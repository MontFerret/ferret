package vm_test

import (
	"sync/atomic"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type bodyCompletionResource struct {
	closeCount *atomic.Int32
	id         uint64
}

func newBodyCompletionResource(id uint64, closeCount *atomic.Int32) *bodyCompletionResource {
	return &bodyCompletionResource{id: id, closeCount: closeCount}
}

func (r *bodyCompletionResource) String() string {
	return "bodyCompletionResource"
}

func (r *bodyCompletionResource) Hash() uint64 {
	return r.id
}

func (r *bodyCompletionResource) Copy() runtime.Value {
	return r
}

func (r *bodyCompletionResource) ResourceID() uint64 {
	return r.id
}

func (r *bodyCompletionResource) Close() error {
	r.closeCount.Add(1)

	return nil
}
