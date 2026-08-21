package testing

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type countingMap struct {
	*runtime.Object
	lookups []string
}

func (m *countingMap) ContainsKey(ctx context.Context, key runtime.Value) (runtime.Boolean, error) {
	m.lookups = append(m.lookups, key.String())

	return m.Object.ContainsKey(ctx, key)
}
