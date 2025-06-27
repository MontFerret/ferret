package internal

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/pkg/runtime"
)

type MultiSorter struct {
	*runtime.Box[runtime.List]
	directions []runtime.SortDirection
	sorted     bool
}

func NewMultiSorter(directions []runtime.SortDirection) Transformer {
	return &MultiSorter{
		Box:        &runtime.Box[runtime.List]{Value: runtime.NewArray(8)},
		directions: directions,
		sorted:     false,
	}
}

func (s *MultiSorter) Iterate(ctx context.Context) (runtime.Iterator, error) {
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

func (s *MultiSorter) Add(ctx context.Context, key, value runtime.Value) error {
	return s.Value.Add(ctx, NewKV(key, value))
}

func (s *MultiSorter) sort(ctx context.Context) error {
	return runtime.SortListWith(ctx, s.Value, func(first, second runtime.Value) int64 {
		firstKV := first.(*KV)
		secondKV := second.(*KV)

		firstKVKey := firstKV.Key.(runtime.List)
		secondKVKey := secondKV.Key.(runtime.List)

		for idx, direction := range s.directions {
			firstKey, _ := firstKVKey.Get(ctx, runtime.NewInt(idx))
			secondKey, _ := secondKVKey.Get(ctx, runtime.NewInt(idx))
			comp := runtime.CompareValues(firstKey, secondKey)

			if comp != 0 {
				if direction == runtime.SortDirectionAsc {
					return comp
				}

				return -comp
			}
		}

		return 0
	})
}

func (s *MultiSorter) Get(_ context.Context, _ runtime.Value) (runtime.Value, error) {
	return runtime.None, runtime.ErrNotSupported
}

func (s *MultiSorter) Length(ctx context.Context) (runtime.Int, error) {
	return s.Value.Length(ctx)
}

func (s *MultiSorter) Close() error {
	val := s.Value
	s.Value = nil

	if closer := val.(io.Closer); closer != nil {
		return closer.Close()
	}

	return nil
}
