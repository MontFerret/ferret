package data

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
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

func (s *Sorter) Iterate(ctx context.Context) (runtime.Iterator, error) {
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

func (s *Sorter) Set(ctx context.Context, key, value runtime.Value) error {
	if err := s.Value.Append(ctx, NewKV(key, value)); err != nil {
		return err
	}

	s.sorted = false

	return nil
}

func (s *Sorter) sort(ctx context.Context) error {
	return runtime.SortListWith(ctx, s.Value, func(ctx context.Context, first, second runtime.Value) (runtime.Ordering, error) {
		firstKV := first.(*KV)
		secondKV := second.(*KV)

		comp, err := runtime.CompareValues(ctx, firstKV.Key, secondKV.Key)
		if err != nil {
			return runtime.Equal, err
		}

		if s.direction == runtime.SortDirectionAsc {
			return comp, nil
		}

		return -comp, nil
	})
}

func (s *Sorter) Get(_ context.Context, _ runtime.Value) (runtime.Value, error) {
	return runtime.None, runtime.ErrNotSupported
}

func (s *Sorter) Length(ctx context.Context) (runtime.Int, error) {
	return s.Value.Length(ctx)
}

func (s *Sorter) Close() error {
	val := s.Value
	s.Value = nil

	if closer, ok := val.(io.Closer); ok {
		return closer.Close()
	}

	return nil
}
