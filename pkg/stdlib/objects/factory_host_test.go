package objects_test

import (
	"context"
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type configuredMap struct {
	newErr   error
	closeErr error
	*runtime.Object
	created   *configuredMap
	onNew     func(*configuredMap)
	onSet     func(context.Context, runtime.Value, runtime.Value) error
	onLookup  func(context.Context, runtime.Value) (runtime.Value, bool, error)
	onForEach func(context.Context, runtime.KeyReadablePredicate) error
	backend   string
	creations int
	closes    int
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
	m.created = &configuredMap{Object: runtime.NewObject(), backend: m.backend}
	if m.onNew != nil {
		m.onNew(m.created)
	}

	if m.newErr != nil {
		return nil, errors.Join(m.newErr, m.created.Close())
	}

	return m.created, nil
}

func (m *configuredMap) Set(ctx context.Context, key, value runtime.Value) error {
	if m.onSet != nil {
		return m.onSet(ctx, key, value)
	}

	return m.Object.Set(ctx, key, value)
}

func (m *configuredMap) Lookup(ctx context.Context, key runtime.Value) (runtime.Value, bool, error) {
	if m.onLookup != nil {
		return m.onLookup(ctx, key)
	}

	return m.Object.Lookup(ctx, key)
}

func (m *configuredMap) ForEach(ctx context.Context, predicate runtime.KeyReadablePredicate) error {
	if m.onForEach != nil {
		return m.onForEach(ctx, predicate)
	}

	return m.Object.ForEach(ctx, predicate)
}

func (m *configuredMap) Close() error {
	m.closes++

	return m.closeErr
}
