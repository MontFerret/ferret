package vm_test

import (
	"context"
	"errors"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	stdlibCancellationSource struct {
		runtime.Value
		expected context.Context
		cancel   context.CancelFunc
		visits   int
		closes   int
	}

	stdlibCancellationIterator struct {
		source *stdlibCancellationSource
	}
)

func (s *stdlibCancellationSource) Iterate(ctx context.Context) (runtime.Iterator, error) {
	if ctx != s.expected {
		return nil, errors.New("iterator context changed")
	}

	return &stdlibCancellationIterator{source: s}, nil
}

func (i *stdlibCancellationIterator) Next(ctx context.Context) (runtime.Value, runtime.Value, error) {
	s := i.source
	if ctx != s.expected {
		return runtime.None, runtime.None, errors.New("next context changed")
	}

	if s.visits == 3 {
		return runtime.None, runtime.None, io.EOF
	}

	s.visits++
	s.cancel()

	return runtime.Int(s.visits), runtime.Int(7), nil
}

func (i *stdlibCancellationIterator) Close() error {
	i.source.closes++

	return nil
}
