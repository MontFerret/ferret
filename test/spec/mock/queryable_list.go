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
