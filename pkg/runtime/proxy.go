package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/wI2L/jettison"
)

type Proxy[T any] struct {
	target any
}

func NewProxy[T any](target T) *Proxy[T] {
	return &Proxy[T]{
		target: target,
	}
}

func (p *Proxy[T]) Target() T {
	// Safeguard against nil target.
	// This can happen if the constructor is bypassed or misused, which is a programming error.
	// In such cases, we return the zero value of T to prevent panics in methods that rely on the target.
	if p.target == nil {
		var zero T

		return zero
	}

	// This should not panic because the constructor enforces the type.
	// If the target is not of type T, it means that the constructor was bypassed or misused, which is a programming error.
	return p.target.(T)
}

func (p *Proxy[T]) Type() Type {
	typed, ok := p.target.(Typed)

	if ok {
		return typed.Type()
	}

	return NewReflectType(p.target)
}

func (p *Proxy[T]) Unwrap() any {
	unwrappable, ok := p.target.(Unwrappable)

	if ok {
		return unwrappable.Unwrap()
	}

	return p.target
}

func (p *Proxy[T]) MarshalJSON() ([]byte, error) {
	marshaler, ok := p.target.(json.Marshaler)

	if ok {
		return marshaler.MarshalJSON()
	}

	return jettison.MarshalOpts(p.target, jettison.NoHTMLEscaping())
}

func (p *Proxy[T]) String() string {
	stringer, ok := p.target.(fmt.Stringer)

	if ok {
		return stringer.String()
	}

	return fmt.Sprintf("%v", p.target)
}

func (p *Proxy[T]) Hash() uint64 {
	hashable, ok := p.target.(Hashable)

	if ok {
		return hashable.Hash()
	}

	return uint64(reflect.ValueOf(p.target).Pointer())
}

func (p *Proxy[T]) Compare(other Value) int64 {
	comp, ok := p.target.(Comparable)

	if ok {
		return comp.Compare(other)
	}

	return -1
}

func (p *Proxy[T]) Copy() Value {
	return NewProxy[T](p.target)
}

func (p *Proxy[T]) Length(ctx context.Context) (Int, error) {
	measurable, ok := p.target.(Measurable)

	if ok {
		return measurable.Length(ctx)
	}

	return -1, ProxyError(p.target, TypeMeasurable)
}

func (p *Proxy[T]) Get(ctx context.Context, key Value) (Value, error) {
	keyReadable, ok := p.target.(KeyReadable)

	if ok {
		return keyReadable.Get(ctx, key)
	}

	return None, ProxyError(p.target, TypeKeyReadable)
}

func (p *Proxy[T]) Set(ctx context.Context, key, value Value) error {
	keyWritable, ok := p.target.(KeyWritable)

	if ok {
		return keyWritable.Set(ctx, key, value)
	}

	return ProxyError(p.target, TypeKeyWritable)
}

func (p *Proxy[T]) Iterate(ctx context.Context) (Iterator, error) {
	iterable, ok := p.target.(Iterable)

	if ok {
		return iterable.Iterate(ctx)
	}

	return nil, ProxyError(p.target, TypeIterable)
}

func (p *Proxy[T]) SortAsc(ctx context.Context) error {
	sortable, ok := p.target.(Sortable)

	if ok {
		return sortable.SortAsc(ctx)
	}

	return ProxyError(p.target, TypeSortable)
}

func (p *Proxy[T]) SortDesc(ctx context.Context) error {
	sortable, ok := p.target.(Sortable)

	if ok {
		return sortable.SortDesc(ctx)
	}

	return ProxyError(p.target, TypeSortable)
}

func (p *Proxy[T]) Dispatch(ctx context.Context, event DispatchEvent) (Value, error) {
	dispatchable, ok := p.target.(Dispatchable)

	if ok {
		return dispatchable.Dispatch(ctx, event)
	}

	return None, ProxyError(p.target, TypeDispatchable)
}

func (p *Proxy[T]) Subscribe(ctx context.Context, subscription Subscription) (Stream, error) {
	observable, ok := p.target.(Observable)

	if ok {
		return observable.Subscribe(ctx, subscription)
	}

	return nil, ProxyError(p.target, TypeObservable)
}

func (p *Proxy[T]) Query(ctx context.Context, q Query) (Value, error) {
	queryable, ok := p.target.(Queryable)

	if ok {
		return queryable.Query(ctx, q)
	}

	return None, ProxyError(p.target, TypeQueryable)
}
