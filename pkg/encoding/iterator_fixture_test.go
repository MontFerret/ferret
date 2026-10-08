package encoding_test

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	codecIterable struct {
		iter runtime.Iterator
	}

	codecIterator struct {
		closeErr error
		nextErr  error
		onNext   func()
		events   *[]string
		items    []runtime.Value
		index    int
		closes   int
	}

	selfCodecIterator struct {
		codecIterator
	}

	resourceCodecIterable struct {
		*codecIterable
		id uint64
	}

	resourceCodecIterator struct {
		*codecIterator
		id uint64
	}
)

func (v *codecIterable) String() string {
	return "iterable"
}

func (v *codecIterable) Hash() uint64 {
	return 0
}

func (v *codecIterable) Copy() runtime.Value {
	return v
}

func (v *codecIterable) Iterate(context.Context) (runtime.Iterator, error) {
	return v.iter, nil
}

func (v *codecIterator) Next(context.Context) (runtime.Value, runtime.Value, error) {
	if v.onNext != nil {
		v.onNext()
	}

	if v.index >= len(v.items) {
		if v.nextErr != nil {
			return runtime.None, runtime.None, v.nextErr
		}

		return runtime.None, runtime.None, io.EOF
	}

	item := v.items[v.index]
	key := runtime.NewInt(v.index)
	v.index++

	return item, key, nil
}

func (v *codecIterator) Close() error {
	v.closes++
	if v.events != nil {
		*v.events = append(*v.events, "close iterator")
	}

	return v.closeErr
}

func (v *selfCodecIterator) String() string {
	return "self iterator"
}

func (v *selfCodecIterator) Hash() uint64 {
	return 0
}

func (v *selfCodecIterator) Copy() runtime.Value {
	return v
}

func (v *selfCodecIterator) Iterate(context.Context) (runtime.Iterator, error) {
	return v, nil
}

func (v *resourceCodecIterable) Close() error {
	return nil
}

func (v *resourceCodecIterable) ResourceID() uint64 {
	return v.id
}

func (v *resourceCodecIterator) ResourceID() uint64 {
	return v.id
}
