package msgpack_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	ownershipContext struct {
		context.Context
		lookups *int
	}

	ownershipMap struct {
		*runtime.Object
	}
)

func (ctx ownershipContext) Value(key any) any {
	*ctx.lookups++

	return ctx.Context.Value(key)
}

func (value *ownershipMap) ForEach(_ context.Context, predicate runtime.KeyReadablePredicate) error {
	return value.Object.ForEach(context.Background(), predicate)
}
