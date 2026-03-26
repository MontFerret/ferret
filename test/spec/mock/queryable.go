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

func (q *Queryable) Query(_ context.Context, query runtime.Query) (runtime.List, error) {
	q.queries = append(q.queries, query)
	if q.err != nil {
		return nil, q.err
	}
	if q.result != nil {
		return q.result, nil
	}
	return runtime.NewArray(0), nil
}

func (q *Queryable) MarshalJSON() ([]byte, error) {
	return json.Marshal("queryable")
}

func (q *Queryable) String() string {
	return "queryable"
}

func (q *Queryable) Unwrap() interface{} {
	return "queryable"
}

func (q *Queryable) Hash() uint64 {
	return 0
}

func (q *Queryable) Copy() runtime.Value {
	return q
}
