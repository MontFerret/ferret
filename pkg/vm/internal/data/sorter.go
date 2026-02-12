package data

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/pkg/runtime"
)

type Sorter struct {
	*runtime.Box[runtime.List]
	direction runtime.SortDirection
	sorted    bool
}

func NewSorter(direction runtime.SortDirection) Transformer {
	return &Sorter{
		Box: &runtime.Box[runtime.List]{
			Value: runtime.NewArray(8),
		},
		direction: direction,
	}
}

func (s *Sorter) Iterate(ctx runtime.Context) (runtime.Iterator, error) {
	if !s.sorted {
		if err := s.sort(ctx); err != nil {
			return nil, err
		}

		s.sorted = true
	}

	iter, err := s.Value.Iterate(ctx)

	if err != nil {
		return nil, err
	}

	return NewKVIterator(iter), nil
}

func (s *Sorter) Set(ctx runtime.Context, key, value runtime.Value) error {
	return s.Value.Append(ctx, NewKV(key, value))
}

func (s *Sorter) sort(ctx context.Context) error {
	return runtime.SortListWith(ctx, s.Value, func(first, second runtime.Value) int64 {
		firstKV := first.(*KV)
		secondKV := second.(*KV)

		comp := runtime.CompareValues(nil, firstKV.Key, secondKV.Key)

		if s.direction == runtime.SortDirectionAsc {
			return comp
		}

		return -comp
	})
}

func (s *Sorter) Get(ctx runtime.Context, key runtime.Value) (runtime.Value, error) {
	return runtime.None, runtime.ErrNotSupported
}

func (s *Sorter) Length(ctx runtime.Context) (runtime.Int, error) {
	return s.Value.Length(ctx)
}

func (s *Sorter) Close() error {
	val := s.Value
	s.Value = nil

	if closer := val.(io.Closer); closer != nil {
		return closer.Close()
	}

	return nil
}
