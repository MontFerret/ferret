package mock

import (
	"context"
	"sync"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type BlockingStream struct {
	ch   chan runtime.Message
	once sync.Once
}

func NewBlockingStream() *BlockingStream {
	return &BlockingStream{
		ch: make(chan runtime.Message),
	}
}

func (s *BlockingStream) Read(ctx context.Context) <-chan runtime.Message {
	return s.ch
}

func (s *BlockingStream) Close() error {
	s.once.Do(func() {
		close(s.ch)
	})

	return nil
}
