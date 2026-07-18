package http

import (
	"io"
	"sync/atomic"
)

type trackedResponseBody struct {
	reader io.Reader
	closed atomic.Bool
}

func newTrackedResponseBody(reader io.Reader) *trackedResponseBody {
	return &trackedResponseBody{reader: reader}
}

func (b *trackedResponseBody) Read(p []byte) (int, error) {
	return b.reader.Read(p)
}

func (b *trackedResponseBody) Close() error {
	b.closed.Store(true)

	return nil
}
