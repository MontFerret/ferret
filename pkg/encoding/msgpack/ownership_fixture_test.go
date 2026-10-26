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

	hostKeyMap struct {
		*runtime.Object
		key runtime.Value
	}

	hostStringKey struct {
		runtime.Value
		calls int
	}
)

func (ctx ownershipContext) Value(key any) any {
	*ctx.lookups++

	return ctx.Context.Value(key)
}

func (value *ownershipMap) ForEach(_ context.Context, predicate runtime.KeyReadablePredicate) error {
	return value.Object.ForEach(context.Background(), predicate)
}

func (value *hostKeyMap) ForEach(ctx context.Context, predicate runtime.KeyReadablePredicate) error {
	return value.Object.ForEach(ctx, func(ctx context.Context, item, _ runtime.Value) (runtime.Boolean, error) {
		return predicate(ctx, item, value.key)
	})
}

func (key *hostStringKey) String() string {
	key.calls++

	return "host key"
}
