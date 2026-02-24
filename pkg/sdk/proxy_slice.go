package sdk

import (
	"context"
	"sort"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// ProxySlice is a proxy for a slice of type T that implements runtime.KeyReadable, runtime.KeyWritable, runtime.Sortable, and runtime.Iterable interfaces if the underlying slice supports them.
type ProxySlice[T any] struct {
	*Proxy[[]T]
	itemTypeName runtime.Type
}

// NewProxySlice creates a new ProxySlice for the given target slice.
func NewProxySlice[T any](target []T) *ProxySlice[T] {
	return &ProxySlice[T]{
		Proxy: NewProxy[[]T](target),
	}
}

// NewProxySliceWithType creates a new ProxySlice for the given target slice and type information.
// This is useful when the target slice does not implement the runtime.Typed interface, or when you want to override the type information for the slice and its items.
func NewProxySliceWithType[T any](typeName, itemTypeName runtime.Type, target []T) *ProxySlice[T] {
	return &ProxySlice[T]{
		Proxy:        NewProxyWithType[[]T](typeName, target),
		itemTypeName: itemTypeName,
	}
}

func (p *ProxySlice[T]) At(_ context.Context, index runtime.Int) (runtime.Value, error) {
	target := p.Target()

	return p.itemToProxy(target[int(index)]), nil
}

func (p *ProxySlice[T]) SetAt(_ context.Context, idx runtime.Int, value runtime.Value) error {
	target := p.Target()

	if idx < 0 || int(idx) >= len(target) {
		return runtime.Error(runtime.ErrInvalidOperation, "out of bounds")
	}

	proxy, ok := value.(*Proxy[T])

	if !ok {
		return runtime.Error(runtime.ErrInvalidType, "expected a proxy of type T")
	}

	target[int(idx)] = proxy.Target()

	return nil
}

func (p *ProxySlice[T]) SortAsc(_ context.Context) error {
	p.sort(true)

	return nil
}

func (p *ProxySlice[T]) SortDesc(ctx context.Context) error {
	p.sort(false)

	return nil
}

func (p *ProxySlice[T]) Iterate(ctx context.Context) (runtime.Iterator, error) {
	iterable, ok := p.target.(runtime.Iterable)

	if ok {
		return iterable.Iterate(ctx)
	}

	return nil, ProxyError(p.target, runtime.TypeIterable)
}

func (p *ProxySlice[T]) itemToProxy(item T) runtime.Value {
	if p.itemTypeName != nil {
		return NewProxyWithType(p.itemTypeName, item)
	}

	return NewProxy(item)
}

func (p *ProxySlice[T]) sort(ascending bool) {
	target := p.Target()
	sort.SliceStable(target, func(i, j int) bool {
		comp := p.compare(target[i], target[j])

		if ascending {
			return comp == -1
		}

		return comp == 1
	})
}

func (p *ProxySlice[T]) compare(a, b T) int64 {
	if comparableA, ok := any(a).(runtime.Comparable); ok {
		return comparableA.Compare(p.itemToProxy(b))
	}

	return 0
}
