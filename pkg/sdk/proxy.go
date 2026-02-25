package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type Proxy[T any] struct {
	target   any
	typeName runtime.Type
}

// NewProxy creates a new Proxy for the given target.
// The target can be of any type and can partially implement the runtime interfaces.
// The Proxy will attempt to delegate method calls to the target if it implements the corresponding interfaces,
// otherwise it will return an error indicating that the target does not support the required interface.
func NewProxy[T any](target T) *Proxy[T] {
	return &Proxy[T]{
		target: target,
	}
}

// NewProxyWithType creates a new Proxy for the given target and type.
// This is useful when the target does not implement the runtime.Typed interface, or when you want to override the type information for the target.
// The provided type will be used for type assertions and method dispatching, instead of the type of the target.
func NewProxyWithType[T any](typeName runtime.Type, target T) *Proxy[T] {
	return &Proxy[T]{
		target:   target,
		typeName: typeName,
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

func (p *Proxy[T]) Type() runtime.Type {
	if p.typeName != nil {
		return p.typeName
	}

	typed, ok := p.target.(runtime.Typed)

	if ok {
		return typed.Type()
	}

	return runtime.HostTypeOf(p.target)
}

func (p *Proxy[T]) Unwrap() any {
	unwrappable, ok := p.target.(runtime.Unwrappable)

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
	hashable, ok := p.target.(runtime.Hashable)

	if ok {
		return hashable.Hash()
	}

	return uint64(reflect.ValueOf(p.target).Pointer())
}

func (p *Proxy[T]) Compare(other runtime.Value) int {
	comp, ok := p.target.(runtime.Comparable)

	if ok {
		return comp.Compare(other)
	}

	return -1
}

func (p *Proxy[T]) Copy() runtime.Value {
	return NewProxy[T](p.Target())
}

func (p *Proxy[T]) Clone(ctx context.Context) (runtime.Cloneable, error) {
	clonable, ok := p.target.(runtime.Cloneable)

	if ok {
		cloned, err := clonable.Clone(ctx)

		if err != nil {
			return nil, err
		}

		return NewProxy[T](cloned.(T)), nil
	}

	return nil, ProxyError(p.target, runtime.TypeCloneable)
}

func (p *Proxy[T]) Length(ctx context.Context) (runtime.Int, error) {
	measurable, ok := p.target.(runtime.Measurable)

	if ok {
		return measurable.Length(ctx)
	}

	return -1, ProxyError(p.target, runtime.TypeMeasurable)
}

func (p *Proxy[T]) Set(ctx context.Context, key, value runtime.Value) error {
	keyWritable, ok := p.target.(runtime.KeyWritable)

	if ok {
		return keyWritable.Set(ctx, key, value)
	}

	return ProxyError(p.target, runtime.TypeKeyWritable)
}

func (p *Proxy[T]) Get(ctx context.Context, key runtime.Value) (runtime.Value, error) {
	keyReadable, ok := p.target.(runtime.KeyReadable)

	if ok {
		return keyReadable.Get(ctx, key)
	}

	return runtime.None, ProxyError(p.target, runtime.TypeKeyReadable)
}

func (p *Proxy[T]) Iterate(ctx context.Context) (runtime.Iterator, error) {
	iterable, ok := p.target.(runtime.Iterable)

	if ok {
		return iterable.Iterate(ctx)
	}

	return nil, ProxyError(p.target, runtime.TypeIterable)
}

func (p *Proxy[T]) Dispatch(ctx context.Context, event runtime.DispatchEvent) (runtime.Value, error) {
	dispatchable, ok := p.target.(runtime.Dispatchable)

	if ok {
		return dispatchable.Dispatch(ctx, event)
	}

	return runtime.None, ProxyError(p.target, runtime.TypeDispatchable)
}

func (p *Proxy[T]) Subscribe(ctx context.Context, subscription runtime.Subscription) (runtime.Stream, error) {
	observable, ok := p.target.(runtime.Observable)

	if ok {
		return observable.Subscribe(ctx, subscription)
	}

	return nil, ProxyError(p.target, runtime.TypeObservable)
}

func (p *Proxy[T]) Query(ctx context.Context, q runtime.Query) (runtime.List, error) {
	queryable, ok := p.target.(runtime.Queryable)

	if ok {
		return queryable.Query(ctx, q)
	}

	return nil, ProxyError(p.target, runtime.TypeQueryable)
}
