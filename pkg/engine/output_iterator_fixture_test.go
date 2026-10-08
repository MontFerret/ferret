package engine

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/vmihailenco/msgpack/v5"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	outputIterator struct {
		closeErr     error
		duplicateErr error
		nextErr      error
		onClose      func()
		onNext       func(context.Context, int)
		items        []runtime.Value
		index        int
		closes       int
	}

	outputIterable struct {
		iter *outputIterator
	}

	outputIteratorLeaf struct {
		closeErr error
		onClose  func()
		text     string
		closes   int
		marshals int
	}

	outputIteratorResource struct {
		*outputIteratorLeaf
		id uint64
	}

	outputIteratorHostList struct {
		*runtime.Array
		reads  int
		closes int
	}
)

func (v *outputIterator) String() string {
	return "output iterator"
}

func (v *outputIterator) Hash() uint64 {
	return 0
}

func (v *outputIterator) Copy() runtime.Value {
	return v
}

func (v *outputIterator) Iterate(context.Context) (runtime.Iterator, error) {
	return v, nil
}

func (v *outputIterator) Next(ctx context.Context) (runtime.Value, runtime.Value, error) {
	if v.closes != 0 {
		return runtime.None, runtime.None, errors.New("iterating a closed resource")
	}

	if v.onNext != nil {
		v.onNext(ctx, v.index)
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

func (v *outputIterator) Close() error {
	v.closes++
	if v.onClose != nil {
		v.onClose()
	}

	if v.closes > 1 {
		return v.duplicateErr
	}

	return v.closeErr
}

func (v *outputIterable) String() string {
	return "output iterable"
}

func (v *outputIterable) Hash() uint64 {
	return 0
}

func (v *outputIterable) Copy() runtime.Value {
	return v
}

func (v *outputIterable) Iterate(context.Context) (runtime.Iterator, error) {
	return v.iter, nil
}

func (v *outputIteratorLeaf) String() string {
	return v.text
}

func (v *outputIteratorLeaf) Hash() uint64 {
	return runtime.NewString(v.text).Hash()
}

func (v *outputIteratorLeaf) Copy() runtime.Value {
	return v
}

func (v *outputIteratorLeaf) MarshalJSON() ([]byte, error) {
	v.marshals++
	if v.closes != 0 {
		return nil, errors.New("encoding a closed yielded resource")
	}

	return json.Marshal(v.text)
}

func (v *outputIteratorLeaf) MarshalMsgpack() ([]byte, error) {
	v.marshals++
	if v.closes != 0 {
		return nil, errors.New("encoding a closed yielded resource")
	}

	return msgpack.Marshal(v.text)
}

func (v *outputIteratorLeaf) Close() error {
	v.closes++
	if v.onClose != nil {
		v.onClose()
	}

	return v.closeErr
}

func (v *outputIteratorResource) ResourceID() uint64 {
	return v.id
}

func (v *outputIteratorHostList) Length(context.Context) (runtime.Int, error) {
	v.reads++

	return 0, errors.New("host list was inspected")
}

func (v *outputIteratorHostList) Close() error {
	v.closes++

	return nil
}
