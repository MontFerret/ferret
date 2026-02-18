package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/wI2L/jettison"
)

type Proxy struct {
	target any
}

func NewProxy(target any) *Proxy {
	return &Proxy{
		target: target,
	}
}

func (p *Proxy) Type() Type {
	typed, ok := p.target.(Typed)

	if ok {
		return typed.Type()
	}

	return NewReflectType(p.target)
}

func (p *Proxy) Unwrap() any {
	unwrappable, ok := p.target.(Unwrappable)

	if ok {
		return unwrappable.Unwrap()
	}

	return p.target
}

func (p *Proxy) MarshalJSON() ([]byte, error) {
	marshaler, ok := p.target.(json.Marshaler)

	if ok {
		return marshaler.MarshalJSON()
	}

	return jettison.MarshalOpts(p.target, jettison.NoHTMLEscaping())
}

func (p *Proxy) String() string {
	stringer, ok := p.target.(fmt.Stringer)

	if ok {
		return stringer.String()
	}

	return fmt.Sprintf("%v", p.target)
}

func (p *Proxy) Hash() uint64 {
	hashable, ok := p.target.(Hashable)

	if ok {
		return hashable.Hash()
	}

	return uint64(reflect.ValueOf(p.target).Pointer())
}

func (p *Proxy) Compare(other Value) int64 {
	comp, ok := other.(Comparable)

	if ok {
		return comp.Compare(p)
	}

	return -1
}

func (p *Proxy) Copy() Value {
	return NewProxy(p.target)
}

func (p *Proxy) Length(ctx context.Context) (Int, error) {
	measurable, ok := p.target.(Measurable)

	if ok {
		return measurable.Length(ctx)
	}

	return -1, ProxyError(p.target, TypeMeasurable)
}

func (p *Proxy) Get(ctx context.Context, key Value) (Value, error) {
	keyReadable, ok := key.(KeyReadable)

	if ok {
		return keyReadable.Get(ctx, key)
	}

	return None, ProxyError(p.target, TypeKeyReadable)
}

func (p *Proxy) Set(ctx context.Context, key, value Value) error {
	keyWritable, ok := key.(KeyWritable)

	if ok {
		return keyWritable.Set(ctx, key, value)
	}

	return ProxyError(p.target, TypeKeyWritable)
}

func (p *Proxy) Iterate(ctx context.Context) (Iterator, error) {
	iterable, ok := p.target.(Iterable)

	if ok {
		return iterable.Iterate(ctx)
	}

	return nil, ProxyError(p.target, TypeIterable)
}

func (p *Proxy) SortAsc(ctx context.Context) error {
	sortable, ok := p.target.(Sortable)

	if ok {
		return sortable.SortAsc(ctx)
	}

	return ProxyError(p.target, TypeSortable)
}

func (p *Proxy) SortDesc(ctx context.Context) error {
	sortable, ok := p.target.(Sortable)

	if ok {
		return sortable.SortDesc(ctx)
	}

	return ProxyError(p.target, TypeSortable)
}

func (p *Proxy) Dispatch(ctx context.Context, event DispatchEvent) (Value, error) {
	dispatchable, ok := p.target.(Dispatchable)

	if ok {
		return dispatchable.Dispatch(ctx, event)
	}

	return None, ProxyError(p.target, TypeDispatchable)
}

func (p *Proxy) Subscribe(ctx context.Context, subscription Subscription) (Stream, error) {
	observable, ok := p.target.(Observable)

	if ok {
		return observable.Subscribe(ctx, subscription)
	}

	return nil, ProxyError(p.target, TypeObservable)
}

func (p *Proxy) Query(ctx context.Context, q Query) (Value, error) {
	queryable, ok := p.target.(Queryable)

	if ok {
		return queryable.Query(ctx, q)
	}

	return None, ProxyError(p.target, TypeQueryable)
}
