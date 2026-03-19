package mem

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// CellHandle uniquely identifies a cell in a value store using token, generation, and slot differentiation.
type CellHandle struct {
	token      uint64
	generation uint64
	slot       uint64
}

func (h CellHandle) ID() uint64 {
	return h.token
}

func (h CellHandle) String() string {
	return fmt.Sprintf("<cell:%d>", h.token)
}

func (h CellHandle) Hash() uint64 {
	return h.token
}

func (h CellHandle) Copy() runtime.Value {
	return h
}

func (h CellHandle) ResourceID() uint64 {
	return h.token
}

func (CellHandle) Close() error {
	return nil
}

func (CellHandle) VMUntracked() {}

func (h CellHandle) valid() bool {
	return h.token != 0 && h.generation != 0 && h.slot != 0
}
