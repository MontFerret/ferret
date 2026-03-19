package mem

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type CellHandle struct {
	id uint64
}

func NewCellHandle(id uint64) CellHandle {
	return CellHandle{id: id}
}

func (h CellHandle) ID() uint64 {
	return h.id
}

func (h CellHandle) String() string {
	return fmt.Sprintf("<cell:%d>", h.id)
}

func (h CellHandle) Hash() uint64 {
	return h.id
}

func (h CellHandle) Copy() runtime.Value {
	return h
}

func (h CellHandle) ResourceID() uint64 {
	return h.id
}

func (CellHandle) Close() error {
	return nil
}

func (CellHandle) VMUntracked() {}

type CellStore struct {
	values map[uint64]runtime.Value
	nextID uint64
}

func (s *CellStore) New(val runtime.Value) CellHandle {
	if s.nextID == 0 {
		s.nextID = 1
	}

	id := s.nextID
	s.nextID++

	if s.values == nil {
		s.values = make(map[uint64]runtime.Value)
	}

	s.values[id] = val

	return NewCellHandle(id)
}

func (s *CellStore) Get(handle CellHandle) (runtime.Value, bool) {
	if s.values == nil || handle.id == 0 {
		return nil, false
	}

	val, ok := s.values[handle.id]

	return val, ok
}

func (s *CellStore) Set(handle CellHandle, val runtime.Value) bool {
	if s.values == nil || handle.id == 0 {
		return false
	}

	if _, ok := s.values[handle.id]; !ok {
		return false
	}

	s.values[handle.id] = val

	return true
}

func (s *CellStore) Delete(handle CellHandle) (runtime.Value, bool) {
	if s.values == nil || handle.id == 0 {
		return nil, false
	}

	val, ok := s.values[handle.id]
	if !ok {
		return nil, false
	}

	delete(s.values, handle.id)

	return val, true
}

func (s *CellStore) Reset() {
	s.nextID = 0

	if s.values != nil {
		clear(s.values)
	}
}
