package vm

import (
	"context"
	"sync/atomic"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type streamGroupSetupStream struct {
	closeErr   error
	closePanic any
	closeCount atomic.Int32
}

func (s *streamGroupSetupStream) Read(context.Context) <-chan runtime.Message {
	return nil
}

func (s *streamGroupSetupStream) Close() error {
	s.closeCount.Add(1)

	if s.closePanic != nil {
		panic(s.closePanic)
	}

	return s.closeErr
}
