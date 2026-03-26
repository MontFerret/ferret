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

	params, err := runtime.ToMap(ctx, q.Options)
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
