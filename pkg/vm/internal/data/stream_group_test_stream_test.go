package data

import (
	"context"
	"sync/atomic"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type streamGroupTestStream struct {
	messages   chan runtime.Message
	closeCount atomic.Int32
}

func newStreamGroupTestStream(capacity int) *streamGroupTestStream {
	return &streamGroupTestStream{messages: make(chan runtime.Message, capacity)}
}

func (s *streamGroupTestStream) Read(context.Context) <-chan runtime.Message {
	return s.messages
}

func (s *streamGroupTestStream) Close() error {
	s.closeCount.Add(1)
	return nil
}
