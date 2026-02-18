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

func (p *Proxy) Copy() Value {
	return NewProxy(p.target)
}

func (p *Proxy) Get(ctx context.Context, key Value) (Value, error) {
	keyReadable, ok := key.(KeyReadable)

	if ok {
		return keyReadable.Get(ctx, key)
	}

	return EncodeField(ctx, p.target, key)
}

func (p *Proxy) Iterate(ctx context.Context) (Iterator, error) {
	iterable, ok := p.target.(Iterable)

	if ok {
		return iterable.Iterate(ctx)
	}

	return nil, fmt.Errorf("cannot iterate over %T", p.target)
}
