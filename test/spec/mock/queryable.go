package mock

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/goccy/go-json"
)

type Queryable struct {
	result  runtime.List
	err     error
	queries []runtime.Query
}

func NewQueryable(result runtime.List) *Queryable {
	return &Queryable{result: result}
}

func (q *Queryable) MockQueries() []runtime.Query {
	return q.queries
}

func (t *Queryable) Query(_ context.Context, q runtime.Query) (runtime.List, error) {
	t.queries = append(t.queries, q)
	if t.err != nil {
		return nil, t.err
	}
	if t.result != nil {
		return t.result, nil
	}
	return runtime.NewArray(0), nil
}

func (t *Queryable) MarshalJSON() ([]byte, error) {
	return json.Marshal("queryable")
}

func (t *Queryable) String() string {
	return "queryable"
}

func (t *Queryable) Unwrap() interface{} {
	return "queryable"
}

func (t *Queryable) Hash() uint64 {
	return 0
}

func (t *Queryable) Copy() runtime.Value {
	return t
}
