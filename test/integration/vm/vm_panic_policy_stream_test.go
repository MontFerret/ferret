package vm_test

import (
	"context"
	"sync/atomic"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type panicPolicyStream struct {
	messages   <-chan runtime.Message
	readPanic  any
	closePanic any
	closeCount *atomic.Int32
}

func (s *panicPolicyStream) Read(context.Context) <-chan runtime.Message {
	if s.readPanic != nil {
		panic(s.readPanic)
	}

	return s.messages
}

func (s *panicPolicyStream) Close() error {
	s.closeCount.Add(1)

	if s.closePanic != nil {
		panic(s.closePanic)
	}

	return nil
}
