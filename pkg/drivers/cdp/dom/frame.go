package dom

import (
	"sync"

	"github.com/mafredri/cdp/protocol/page"
)

type (
	Frame struct {
		tree page.FrameTree
		node *HTMLDocument
	}

	AtomicFrameID struct {
		mu    sync.Mutex
		value page.FrameID
	}

	AtomicFrameCollection struct {
		mu    sync.Mutex
		value map[page.FrameID]Frame
	}
)

func NewAtomicFrameID() *AtomicFrameID {
	return &AtomicFrameID{}
}

func (id *AtomicFrameID) Get() page.FrameID {
	id.mu.Lock()
	defer id.mu.Unlock()

	return id.value
}

func (id *AtomicFrameID) Set(value page.FrameID) {
	id.mu.Lock()
	defer id.mu.Unlock()

	id.value = value
}

func (id *AtomicFrameID) Reset() {
	id.mu.Lock()
	defer id.mu.Unlock()

	id.value = ""
}

func (id *AtomicFrameID) IsEmpty() bool {
	id.mu.Lock()
	defer id.mu.Unlock()

	return id.value == ""
}

func NewAtomicFrameCollection() *AtomicFrameCollection {
	return &AtomicFrameCollection{
		value: make(map[page.FrameID]Frame),
	}
}

func (fc *AtomicFrameCollection) Length() int {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	return len(fc.value)
}

func (fc *AtomicFrameCollection) ForEach(predicate func(value Frame, key page.FrameID) bool) {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	for k, v := range fc.value {
		if predicate(v, k) == false {
			break
		}
	}
}

func (fc *AtomicFrameCollection) Has(key page.FrameID) bool {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	_, ok := fc.value[key]

	return ok
}

func (fc *AtomicFrameCollection) Get(key page.FrameID) (Frame, bool) {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	found, ok := fc.value[key]

	if ok {
		return found, ok
	}

	return Frame{}, false
}

func (fc *AtomicFrameCollection) Set(key page.FrameID, value Frame) {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	fc.value[key] = value
}

func (fc *AtomicFrameCollection) Remove(key page.FrameID) {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	delete(fc.value, key)
}

func (fc *AtomicFrameCollection) ToSlice() []Frame {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	slice := make([]Frame, 0, len(fc.value))

	for _, v := range fc.value {
		slice = append(slice, v)
	}

	return slice
}
