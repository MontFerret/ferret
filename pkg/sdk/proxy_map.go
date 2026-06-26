package sdk

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// ProxyMap is a proxy for map types that implements key access, key removal, value removal, and iteration interfaces.
// It allows you to interact with the underlying map using the runtime interfaces, while also providing type safety for the keys and values.
type ProxyMap[TKey comparable, TValue any] struct {
	*Proxy[map[TKey]TValue]
	itemTypeName runtime.Type
}

// NewProxyMap creates a new ProxyMap for the given map data.
func NewProxyMap[TKey comparable, TValue any](data map[TKey]TValue) *ProxyMap[TKey, TValue] {
	return &ProxyMap[TKey, TValue]{Proxy: NewProxy[map[TKey]TValue](data)}
}

// NewProxyMapWithType creates a new ProxyMap for the given map data and type information.
// This is useful when the target map does not implement the runtime.Typed interface, or when you want to override the type information for the map and its items.
func NewProxyMapWithType[TKey comparable, TValue any](typeName, itemTypeName runtime.Type, data map[TKey]TValue) *ProxyMap[TKey, TValue] {
	return &ProxyMap[TKey, TValue]{
		Proxy:        NewProxyWithType[map[TKey]TValue](typeName, data),
		itemTypeName: itemTypeName,
	}
}

func (p *ProxyMap[TKey, TValue]) Get(ctx context.Context, key runtime.Value) (runtime.Value, error) {
	keyReadable, ok := p.target.(runtime.KeyReadable)

	if ok {
		return keyReadable.Get(ctx, key)
	}

	m, ok := p.target.(map[string]TValue)

	if ok {
		strKey := key.String()

		if value, found := m[strKey]; found {
			return p.itemToProxy(value), nil
		}

		return runtime.None, nil
	}

	return runtime.None, ProxyError(p.target, runtime.TypeKeyReadable)
}

func (p *ProxyMap[TKey, TValue]) Set(ctx context.Context, key, value runtime.Value) error {
	keyWritable, ok := p.target.(runtime.KeyWritable)

	if ok {
		return keyWritable.Set(ctx, key, value)
	}

	m, ok := p.target.(map[string]TValue)

	if ok {
		strKey := key.String()
		m[strKey] = value.(TValue)

		return nil
	}

	return ProxyError(p.target, runtime.TypeKeyWritable)
}

func (p *ProxyMap[TKey, TValue]) RemoveKey(ctx context.Context, key runtime.Value) error {
	keyRemovable, ok := p.target.(runtime.KeyRemovable)

	if ok {
		return keyRemovable.RemoveKey(ctx, key)
	}

	mapKey, ok := p.keyFromValue(key)

	if !ok {
		return ProxyError(p.target, runtime.TypeKeyRemovable)
	}

	delete(p.Target(), mapKey)

	return nil
}

func (p *ProxyMap[TKey, TValue]) Remove(ctx context.Context, value runtime.Value) error {
	valueRemovable, ok := p.target.(runtime.ValueRemovable)

	if ok {
		return valueRemovable.Remove(ctx, value)
	}

	target := p.Target()

	for key, item := range target {
		if proxyValueEqual(value, item, p.itemToProxy(item)) {
			delete(target, key)

			break
		}
	}

	return nil
}

func (p *ProxyMap[TKey, TValue]) Iterate(ctx context.Context) (runtime.Iterator, error) {
	iterable, ok := p.target.(runtime.Iterable)

	if ok {
		return iterable.Iterate(ctx)
	}

	return NewMapIterator[TKey, TValue](p.Target()), nil
}

func (p *ProxyMap[TKey, TValue]) itemToProxy(item TValue) runtime.Value {
	if p.itemTypeName != nil {
		return NewProxyWithType(p.itemTypeName, item)
	}

	return NewProxy(item)
}

func (p *ProxyMap[TKey, TValue]) keyFromValue(key runtime.Value) (TKey, bool) {
	if key == nil {
		var zero TKey

		return zero, false
	}

	if typedKey, ok := any(key).(TKey); ok {
		return typedKey, true
	}

	var out TKey

	if err := Decode(key, &out); err != nil {
		return out, false
	}

	return out, true
}
