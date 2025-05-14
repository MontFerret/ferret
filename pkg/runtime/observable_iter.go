package runtime

import (
	"context"
	"time"
)

type StreamIterator struct {
	stream      Stream
	channel     <-chan Message
	message     Message
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

func (s *StreamIterator) HasNext(ctx context.Context) (bool, error) {
	if !s.initialized {
		s.channel = s.stream.Read(ctx)
		s.initialized = true
	}

	if s.closed {
		return false, nil
	}

	var message Message
	var isOpen bool

	select {
	case message, isOpen = <-s.channel:
	case <-time.After(s.timeout):
		return false, ErrTimeout
	}

	if !isOpen {
		s.closed = true

		return false, nil
	}

	s.message = message

	return true, nil
}

func (s *StreamIterator) Next(_ context.Context) (value Value, key Value, err error) {
	return s.message.Value(), None, s.message.Err()
}

func (s *StreamIterator) Close() error {
	if s.closed {
		return nil
	}

	s.closed = true
	return s.stream.Close()
}
