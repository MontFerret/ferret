package runtime_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type opaqueHostValue struct {
	typ  runtime.Type
	hash uint64
}

func (v *opaqueHostValue) String() string {
	panic("comparison dispatch inspected String")
}

func (v *opaqueHostValue) Hash() uint64 {
	return v.hash
}

func (v *opaqueHostValue) Copy() runtime.Value {
	panic("comparison dispatch inspected Copy")
}

func (v *opaqueHostValue) Type() runtime.Type {
	if v.typ == nil {
		return comparisonContractType
	}

	return v.typ
}

func (v *opaqueHostValue) Unwrap() any {
	panic("comparison dispatch inspected Unwrap")
}

func (v *opaqueHostValue) MarshalJSON() ([]byte, error) {
	panic("comparison dispatch inspected MarshalJSON")
}

func (v *opaqueHostValue) Iterate(context.Context) (runtime.Iterator, error) {
	panic("comparison dispatch inspected Iterate")
}
