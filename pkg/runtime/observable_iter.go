package runtime

import (
	"context"
	"io"
	"time"
)

type StreamIterator struct {
	stream      Stream
	channel     <-chan Message
	timeout     time.Duration
	initialized bool
	closed      bool
}

func NewIterator(stream Stream) Iterator {
	return NewIteratorWithTimeout(stream, DefaultStreamTimeout)
}

func NewIteratorWithTimeout(stream Stream, timeout time.Duration) Iterator {
	return &StreamIterator{
		stream:  stream,
		timeout: timeout * time.Millisecond,
	}
}

func (s *StreamIterator) Next(ctx context.Context) (value Value, key Value, err error) {
	if !s.initialized {
		s.channel = s.stream.Read(ctx)
		s.initialized = true
	}

	if s.closed {
		return None, None, io.EOF
	}

	var message Message
	var isOpen bool

	select {
	case message, isOpen = <-s.channel:
	case <-time.After(s.timeout):
		return None, None, ErrTimeout
	}

	if !isOpen {
		s.closed = true

		return None, None, io.EOF
	}

	return message.Value(), None, message.Err()
}

func (s *StreamIterator) Close() error {
	if s.closed {
		return nil
	}

	s.closed = true
	return s.stream.Close()
}
