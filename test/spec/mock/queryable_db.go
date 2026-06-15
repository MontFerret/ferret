package mock

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type DBQueryable struct {
	Queryable
}

func NewDBQueryable() *DBQueryable {
	return &DBQueryable{}
}

func (m *DBQueryable) Query(ctx context.Context, q runtime.Query) (runtime.List, error) {
	m.queries = append(m.queries, q)

	if q.Kind.String() == "nil" {
		return runtime.NewArray(0), nil
	}

	if q.Kind.String() != "sql" {
		return runtime.NewArray(0), nil
	}

	params, err := runtime.ToMap(ctx, q.Params)
	if err != nil {
		return nil, err
	}

	category, _ := params.Get(ctx, runtime.NewString("c"))
	if category == runtime.NewString("laptops") {
		return runtime.NewArrayWith(
			runtime.NewObjectWith(map[string]runtime.Value{
				"name":  runtime.NewString("Laptop Pro"),
				"price": runtime.NewInt(200),
			}),
		), nil
	}

	return runtime.NewArrayWith(
		runtime.NewObjectWith(map[string]runtime.Value{
			"name":  runtime.NewString("Laptop Pro"),
			"price": runtime.NewInt(200),
		}),
		runtime.NewObjectWith(map[string]runtime.Value{
			"name":  runtime.NewString("Mouse"),
			"price": runtime.NewInt(50),
		}),
	), nil
}

func (m *DBQueryable) QueryOne(ctx context.Context, q runtime.Query) (runtime.Value, error) {
	return runtime.DefaultQueryOne(ctx, q, m.Query)
}

func (m *DBQueryable) QueryCount(ctx context.Context, q runtime.Query) (runtime.Int, error) {
	return runtime.DefaultQueryCount(ctx, q, m.Query)
}

func (m *DBQueryable) QueryExists(ctx context.Context, q runtime.Query) (runtime.Boolean, error) {
	return runtime.DefaultQueryExists(ctx, q, m.Query)
}
