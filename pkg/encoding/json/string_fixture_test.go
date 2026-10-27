package json_test

import (
	"bytes"
	"context"
)

type stringChunkWriter struct {
	failure error
	cancel  context.CancelFunc
	sizes   []int
	bytes.Buffer
	failAt int
}

func (w *stringChunkWriter) Write(data []byte) (int, error) {
	w.sizes = append(w.sizes, len(data))
	if len(w.sizes) == w.failAt {
		return 0, w.failure
	}

	n, err := w.Buffer.Write(data)
	if w.cancel != nil {
		w.cancel()
	}

	return n, err
}
