package data

import (
	"context"
	"io"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type streamIterator struct {
	stream      runtime.Stream
	channel     <-chan runtime.Message
	timeout     time.Duration
	initialized bool
	closed      bool
}

func newStreamIterator(stream runtime.Stream, timeout runtime.Duration) runtime.Iterator {
	return &streamIterator{
		stream:  stream,
		timeout: time.Duration(timeout),
	}
}

func (s *streamIterator) Next(ctx context.Context) (value runtime.Value, key runtime.Value, err error) {
	if !s.initialized {
		s.channel = s.stream.Read(ctx)
		s.initialized = true
	}

	if s.closed {
		return runtime.None, runtime.None, io.EOF
	}

	var message runtime.Message
	var isOpen bool

	select {
	case message, isOpen = <-s.channel:
	case <-time.After(s.timeout):
		return runtime.None, runtime.None, runtime.ErrTimeout
	}

	if !isOpen {
		s.closed = true
		return runtime.None, runtime.None, io.EOF
	}

	return message.Value(), runtime.None, message.Err()
}

func (s *streamIterator) Close() error {
	if s.closed {
		return nil
	}

	s.closed = true
	return s.stream.Close()
}
