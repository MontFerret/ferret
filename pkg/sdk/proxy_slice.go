package sdk

import (
	"context"
	"sort"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// ProxySlice is a proxy for a slice of type T that implements indexed access, index removal, value removal, sorting, and iteration interfaces.
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

func (p *ProxySlice[T]) RemoveAt(ctx context.Context, idx runtime.Int) (runtime.Value, error) {
	indexRemovable, ok := p.target.(runtime.IndexRemovable)

	if ok {
		return indexRemovable.RemoveAt(ctx, idx)
	}

	target := p.Target()
	edge := runtime.Int(len(target) - 1)

	if idx > edge {
		return runtime.None, nil
	}

	if idx < 0 {
		return runtime.None, runtime.Error(runtime.ErrInvalidOperation, "out of bounds")
	}

	item := target[idx]
	p.target = append(target[:idx], target[idx+1:]...)

	return p.itemToProxy(item), nil
}

func (p *ProxySlice[T]) Remove(ctx context.Context, value runtime.Value) error {
	valueRemovable, ok := p.target.(runtime.ValueRemovable)

	if ok {
		return valueRemovable.Remove(ctx, value)
	}

	target := p.Target()

	for idx, item := range target {
		if proxyValueEqual(value, item, p.itemToProxy(item)) {
			_, err := p.RemoveAt(ctx, runtime.Int(idx))

			return err
		}
	}

	return nil
}

func (p *ProxySlice[T]) RemoveKey(ctx context.Context, key runtime.Value) error {
	switch idx := key.(type) {
	case runtime.Int:
		_, err := p.RemoveAt(ctx, idx)

		return err
	case runtime.Float:
		return runtime.TypeErrorOf(idx, runtime.TypeInt)
	default:
		return ProxyError(p.target, runtime.TypeIndexRemovable)
	}
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

func (p *ProxySlice[T]) compare(a, b T) int {
	if comparableA, ok := any(a).(runtime.Comparable); ok {
		return comparableA.Compare(p.itemToProxy(b))
	}

	return 0
}
