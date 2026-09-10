package objects_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	readableObject struct {
		runtime.Value
	}

	fallibleRemovalMap struct {
		*hostMap
		failure error
		calls   int
	}

	fallibleWriteMap struct {
		*hostMap
		failure error
		calls   int
	}
)

func (m *readableObject) Get(context.Context, runtime.Value) (runtime.Value, error) {
	return runtime.None, nil
}

func (m *fallibleRemovalMap) RemoveKey(ctx context.Context, key runtime.Value) error {
	m.calls++
	if m.calls == 2 {
		return m.failure
	}

	return m.hostMap.RemoveKey(ctx, key)
}

func (m *fallibleWriteMap) Set(ctx context.Context, key, value runtime.Value) error {
	m.calls++
	if m.calls == 2 {
		return m.failure
	}

	return m.hostMap.Set(ctx, key, value)
}
