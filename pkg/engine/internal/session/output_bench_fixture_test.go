package session

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	adoptionBenchResource struct{ runtime.Value }
	adoptionBenchIterable struct{ items []runtime.Value }
	adoptionBenchIterator struct {
		items []runtime.Value
		index int
	}
)

func (*adoptionBenchResource) Close() error { return nil }

func (*adoptionBenchResource) MarshalJSON() ([]byte, error) { return []byte(`"value"`), nil }

func (*adoptionBenchResource) MarshalMsgpack() ([]byte, error) {
	return []byte{0xa5, 'v', 'a', 'l', 'u', 'e'}, nil
}

func (v *adoptionBenchIterable) String() string { return "iterable" }

func (v *adoptionBenchIterable) Hash() uint64 { return 0 }

func (v *adoptionBenchIterable) Copy() runtime.Value { return v }

func (v *adoptionBenchIterable) Iterate(context.Context) (runtime.Iterator, error) {
	return &adoptionBenchIterator{items: v.items}, nil
}

func (v *adoptionBenchIterator) Next(context.Context) (runtime.Value, runtime.Value, error) {
	if v.index == len(v.items) {
		return runtime.None, runtime.None, io.EOF
	}

	value := v.items[v.index]
	v.index++

	return value, runtime.None, nil
}
