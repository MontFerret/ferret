package encoding_test

import (
	"bytes"
	"context"
	"io"
)

type (
	failingWriter struct {
		failure error
		bytes.Buffer
		limit  int
		calls  int
		closed bool
	}

	chunkReader struct {
		failure error
		cancel  context.CancelFunc
		data    []byte
		chunk   int
		calls   int
		closed  bool
	}

	cancelWriter struct {
		cancel context.CancelFunc
		bytes.Buffer
	}
)

func (w *failingWriter) Write(data []byte) (int, error) {
	w.calls++
	size := min(len(data), max(0, w.limit-w.Len()))
	n, _ := w.Buffer.Write(data[:size])

	if size < len(data) {
		return n, w.failure
	}

	return n, nil
}

func (w *failingWriter) Close() error {
	w.closed = true

	return nil
}

func (r *chunkReader) Read(dst []byte) (int, error) {
	r.calls++

	if r.cancel != nil {
		r.cancel()
		r.cancel = nil
	}

	size := min(len(dst), r.chunk, len(r.data))
	copy(dst, r.data[:size])
	r.data = r.data[size:]
	if len(r.data) == 0 {
		if r.failure != nil {
			return size, r.failure
		}

		return size, io.EOF
	}

	return size, nil
}

func (r *chunkReader) Close() error {
	r.closed = true

	return nil
}

func (w *cancelWriter) Write(data []byte) (int, error) {
	n, err := w.Buffer.Write(data)

	if w.cancel != nil {
		w.cancel()
		w.cancel = nil
	}

	return n, err
}
