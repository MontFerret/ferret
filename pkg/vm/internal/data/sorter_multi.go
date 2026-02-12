package data

import (
	"io"

	"github.com/MontFerret/ferret/pkg/runtime"
)

type MultiSorter struct {
	*runtime.Box[runtime.List]
	alloc      runtime.Allocator
	directions []runtime.SortDirection
	sorted     bool
}

func NewMultiSorter(alloc runtime.Allocator, directions []runtime.SortDirection) Transformer {
	return &MultiSorter{
		Box:        &runtime.Box[runtime.List]{Value: alloc.Array(8)},
		directions: directions,
		sorted:     false,
	}
}

func (s *MultiSorter) Iterate(ctx runtime.Context) (runtime.Iterator, error) {
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

func (s *MultiSorter) Set(ctx runtime.Context, key, value runtime.Value) error {
	return s.Value.Append(ctx, NewKV(key, value))
}

func (s *MultiSorter) sort(ctx runtime.Context) error {
	return runtime.SortListWith(ctx, s.Value, func(ctx runtime.Context, first, second runtime.Value) int64 {
		firstKV := first.(*KV)
		secondKV := second.(*KV)

		firstKVKey := firstKV.Key.(runtime.List)
		secondKVKey := secondKV.Key.(runtime.List)

		for idx, direction := range s.directions {
			firstKey, _ := firstKVKey.Get(ctx, runtime.NewInt(idx))
			secondKey, _ := secondKVKey.Get(ctx, runtime.NewInt(idx))
			comp := runtime.CompareValues(ctx, firstKey, secondKey)

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

func (s *MultiSorter) Get(_ runtime.Context, _ runtime.Value) (runtime.Value, error) {
	return runtime.None, runtime.ErrNotSupported
}

func (s *MultiSorter) Length(ctx runtime.Context) (runtime.Int, error) {
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
