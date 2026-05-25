package mock

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type JSONQueryable struct {
	Queryable
}

func NewJSONQueryable() *JSONQueryable {
	return &JSONQueryable{}
}

func (m *JSONQueryable) Query(_ context.Context, q runtime.Query) (runtime.List, error) {
	m.queries = append(m.queries, q)

	if q.Kind.String() != "jp" {
		return runtime.NewArray(0), nil
	}

	orders := runtime.NewArrayWith(
		runtime.NewObjectWith(map[string]runtime.Value{
			"id":    runtime.NewInt(1),
			"total": runtime.NewInt(150),
			"items": runtime.NewArrayWith(
				runtime.NewObjectWith(map[string]runtime.Value{"name": runtime.NewString("Item A")}),
				runtime.NewObjectWith(map[string]runtime.Value{"name": runtime.NewString("Item B")}),
			),
		}),
		runtime.NewObjectWith(map[string]runtime.Value{
			"id":    runtime.NewInt(2),
			"total": runtime.NewInt(80),
			"items": runtime.NewArrayWith(
				runtime.NewObjectWith(map[string]runtime.Value{"name": runtime.NewString("Item C")}),
			),
		}),
	)

	return orders, nil
}

func (m *JSONQueryable) QueryOne(ctx context.Context, q runtime.Query) (runtime.Value, error) {
	return runtime.DefaultQueryOne(ctx, q, m.Query)
}

func (m *JSONQueryable) QueryCount(ctx context.Context, q runtime.Query) (runtime.Int, error) {
	return runtime.DefaultQueryCount(ctx, q, m.Query)
}

func (m *JSONQueryable) QueryExists(ctx context.Context, q runtime.Query) (runtime.Boolean, error) {
	return runtime.DefaultQueryExists(ctx, q, m.Query)
}
