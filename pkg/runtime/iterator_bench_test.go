package runtime

import (
	"context"
	"errors"
	"io"
	"testing"
)

type eofIter struct {
	n int
	i int
}

func (it *eofIter) Next(_ context.Context) (Value, Value, error) {
	if it.i >= it.n {
		return None, None, io.EOF
	}
	it.i++
	return ZeroInt, ZeroInt, nil
}

func BenchmarkIteratorEOF(b *testing.B) {
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		it := &eofIter{n: 1024}
		for {
			_, _, err := it.Next(ctx)
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				break
			}
		}
	}
}
