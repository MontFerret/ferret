package mock

import (
	"sync/atomic"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type Message struct {
	value runtime.Value
	obs   *Observable
	err   error
}

func (m *Message) Value() runtime.Value {
	if m.obs != nil {
		atomic.AddInt32(&m.obs.readCount, 1)
	}

	return m.value
}

func (m *Message) Err() error {
	return m.err
}
