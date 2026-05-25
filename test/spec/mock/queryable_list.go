package mock

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/goccy/go-json"
)

type NilListQueryable struct{}

func NewNilListQueryable() *NilListQueryable {
	return &NilListQueryable{}
}

func (n *NilListQueryable) Query(_ context.Context, _ runtime.Query) (runtime.List, error) {
	return nil, nil
}

func (n *NilListQueryable) QueryOne(ctx context.Context, q runtime.Query) (runtime.Value, error) {
	return runtime.DefaultQueryOne(ctx, q, n.Query)
}

func (n *NilListQueryable) QueryCount(ctx context.Context, q runtime.Query) (runtime.Int, error) {
	return runtime.DefaultQueryCount(ctx, q, n.Query)
}

func (n *NilListQueryable) QueryExists(ctx context.Context, q runtime.Query) (runtime.Boolean, error) {
	return runtime.DefaultQueryExists(ctx, q, n.Query)
}

func (n *NilListQueryable) MarshalJSON() ([]byte, error) {
	return json.Marshal("nil-queryable")
}

func (n *NilListQueryable) String() string {
	return "nil-queryable"
}

func (n *NilListQueryable) Hash() uint64 {
	return 0
}

func (n *NilListQueryable) Copy() runtime.Value {
	return n
}
