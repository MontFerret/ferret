package collections_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	benchmarkList struct {
		*runtime.Array
		backend string
	}

	benchmarkIterable struct {
		values *runtime.Array
	}
)

func (l *benchmarkList) New(ctx context.Context) (runtime.List, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	size, err := l.Length(ctx)
	if err != nil {
		return nil, err
	}

	return &benchmarkList{Array: runtime.NewArray64(size), backend: l.backend}, nil
}

func (s *benchmarkIterable) String() string      { return "benchmark iterable" }
func (s *benchmarkIterable) Hash() uint64        { return 1 }
func (s *benchmarkIterable) Copy() runtime.Value { return s }
func (s *benchmarkIterable) Iterate(ctx context.Context) (runtime.Iterator, error) {
	return s.values.Iterate(ctx)
}
