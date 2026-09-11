package objects_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type configuredMap struct {
	newErr error
	*runtime.Object
	backend   string
	creations int
}

var (
	_ runtime.Map                  = (*configuredMap)(nil)
	_ runtime.Factory[runtime.Map] = (*configuredMap)(nil)
)

func (m *configuredMap) New(ctx context.Context) (runtime.Map, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	m.creations++
	if m.newErr != nil {
		return nil, m.newErr
	}

	return &configuredMap{Object: runtime.NewObject(), backend: m.backend}, nil
}
