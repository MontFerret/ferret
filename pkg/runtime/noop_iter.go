package runtime

import (
	"context"
	"io"
)

type noopIterator struct{}

var NoopIterator = noopIterator{}

func (n noopIterator) Next(_ context.Context) (value Value, key Value, err error) {
	return None, None, io.EOF
}
