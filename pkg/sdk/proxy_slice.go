package sdk

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// ProxySlice is a proxy for a slice of type T that implements runtime.KeyReadable, runtime.KeyWritable, runtime.Sortable, and runtime.Iterable interfaces if the underlying slice supports them.
type ProxySlice[T any] struct {
	*Proxy[[]T]
}

// NewProxySlice creates a new ProxySlice for the given target slice.
func NewProxySlice[T any](target []T) *ProxySlice[T] {
	return &ProxySlice[T]{
		Proxy: NewProxy[[]T](target),
	}
}

func (p *ProxySlice[T]) At(ctx context.Context, key runtime.Int) (runtime.Value, error) {
	keyReadable, ok := p.target.(runtime.KeyReadable)

	if ok {
		return keyReadable.Get(ctx, key)
	}

	return runtime.None, ProxyError(p.target, runtime.TypeKeyReadable)
}

func (p *ProxySlice[T]) Set(ctx context.Context, key, value runtime.Int) error {
	keyWritable, ok := p.target.(runtime.KeyWritable)

	if ok {
		return keyWritable.Set(ctx, key, value)
	}

	return ProxyError(p.target, runtime.TypeKeyWritable)
}

func (p *ProxySlice[T]) SortAsc(ctx context.Context) error {
	sortable, ok := p.target.(runtime.Sortable)

	if ok {
		return sortable.SortAsc(ctx)
	}

	return ProxyError(p.target, runtime.TypeSortable)
}

func (p *ProxySlice[T]) SortDesc(ctx context.Context) error {
	sortable, ok := p.target.(runtime.Sortable)

	if ok {
		return sortable.SortDesc(ctx)
	}

	return ProxyError(p.target, runtime.TypeSortable)
}

func (p *ProxySlice[T]) Iterate(ctx context.Context) (runtime.Iterator, error) {
	iterable, ok := p.target.(runtime.Iterable)

	if ok {
		return iterable.Iterate(ctx)
	}

	return nil, ProxyError(p.target, runtime.TypeIterable)
}
