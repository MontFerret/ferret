package sdk_test

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type conversionIteratorValue struct {
	terminalErr error
	closeErr    error
	values      []runtime.Value
	length      runtime.Int
	position    int
	closed      bool
}

func (value *conversionIteratorValue) String() string {
	return "conversionIteratorValue"
}

func (value *conversionIteratorValue) Hash() uint64 {
	return 1
}

func (value *conversionIteratorValue) Copy() runtime.Value {
	if value == nil {
		return runtime.None
	}

	copy := *value
	return &copy
}

func (value *conversionIteratorValue) Iterate(ctx context.Context) (runtime.Iterator, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return value, nil
}

func (value *conversionIteratorValue) Length(ctx context.Context) (runtime.Int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	if value.length != 0 {
		return value.length, nil
	}

	return runtime.Int(len(value.values)), nil
}

func (value *conversionIteratorValue) Next(ctx context.Context) (runtime.Value, runtime.Value, error) {
	if err := ctx.Err(); err != nil {
		return runtime.None, runtime.None, err
	}
	if value.position >= len(value.values) {
		if value.terminalErr != nil {
			return runtime.None, runtime.None, value.terminalErr
		}

		return runtime.None, runtime.None, io.EOF
	}

	position := value.position
	item := value.values[position]
	value.position++

	return item, runtime.NewInt(position), nil
}

func (value *conversionIteratorValue) Close() error {
	value.closed = true
	return value.closeErr
}
