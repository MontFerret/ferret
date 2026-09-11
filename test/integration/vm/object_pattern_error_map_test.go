package vm_test

import (
	"context"
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type errorMap struct {
	*runtime.Object
}

func (m *errorMap) ContainsKey(ctx context.Context, key runtime.Value) (runtime.Boolean, error) {
	return runtime.False, errors.New("boom")
}

func (m *errorMap) New(ctx context.Context) (runtime.Map, error) {
	base, err := m.Object.New(ctx)
	if err != nil {
		return nil, err
	}

	result := *m
	result.Object = base.(*runtime.Object)

	return &result, nil
}
