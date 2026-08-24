package test

import (
	"io"
	"reflect"
)

type Testing[B io.Closer] struct {
	Benchmark B
	Options   Options
}

// NewTesting creates VM testing state from the supplied options.
func NewTesting[B io.Closer](opts []Option) (Testing[B], error) {
	cfg, err := NewOptions(opts)
	if err != nil {
		return Testing[B]{}, err
	}

	return Testing[B]{Options: cfg}, nil
}

func (t *Testing[B]) SetBenchmark(benchmark B) {
	t.Benchmark = benchmark
}

func (t *Testing[B]) Close() {
	if !isZero(t.Benchmark) {
		t.Benchmark.Close()
	}

	var zero B
	t.Benchmark = zero
}

func isZero[T any](value T) bool {
	rv := reflect.ValueOf(value)
	if !rv.IsValid() {
		return true
	}

	return rv.IsZero()
}
