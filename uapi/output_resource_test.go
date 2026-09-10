package uapi

import (
	"sync/atomic"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type outputResource struct {
	closes atomic.Int32
}

func (r *outputResource) String() string {
	return "output resource"
}

func (r *outputResource) Hash() uint64 {
	return runtime.NewString(r.String()).Hash()
}

func (r *outputResource) Copy() runtime.Value {
	return r
}

func (r *outputResource) Close() error {
	r.closes.Add(1)

	return nil
}
