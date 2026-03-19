package mem

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

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

type CellStore struct {
	values     map[CellHandle]runtime.Value
	generation uint64
	nextSlot   uint64
	nextToken  uint64
}

func (s *CellStore) New(val runtime.Value) CellHandle {
	if s.generation == 0 {
		s.generation = 1
	}
	if s.nextSlot == 0 {
		s.nextSlot = 1
	}
	if s.nextToken == 0 {
		s.nextToken = 1
	}

	handle := CellHandle{
		token:      s.nextToken,
		generation: s.generation,
		slot:       s.nextSlot,
	}
	s.nextSlot++
	s.nextToken++

	if s.values == nil {
		s.values = make(map[CellHandle]runtime.Value)
	}

	s.values[handle] = val

	return handle
}

func (s *CellStore) Get(handle CellHandle) (runtime.Value, bool) {
	if s.values == nil || !handle.valid() {
		return nil, false
	}

	val, ok := s.values[handle]

	return val, ok
}

func (s *CellStore) Set(handle CellHandle, val runtime.Value) bool {
	if s.values == nil || !handle.valid() {
		return false
	}

	if _, ok := s.values[handle]; !ok {
		return false
	}

	s.values[handle] = val

	return true
}

func (s *CellStore) Delete(handle CellHandle) (runtime.Value, bool) {
	if s.values == nil || !handle.valid() {
		return nil, false
	}

	val, ok := s.values[handle]
	if !ok {
		return nil, false
	}

	delete(s.values, handle)

	return val, true
}

func (s *CellStore) Reset() {
	if s.generation == 0 {
		s.generation = 1
	} else {
		s.generation++
		if s.generation == 0 {
			s.generation = 1
		}
	}
	s.nextSlot = 0

	if s.values != nil {
		clear(s.values)
	}
}

func (h CellHandle) valid() bool {
	return h.token != 0 && h.generation != 0 && h.slot != 0
}
