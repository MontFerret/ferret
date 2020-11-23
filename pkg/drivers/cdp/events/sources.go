package events

import (
	"errors"
	"sync"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type SourceCollection struct {
	mu     sync.RWMutex
	values []Source
}

func NewSourceCollection() *SourceCollection {
	sc := new(SourceCollection)
	sc.values = make([]Source, 0, 10)

	return sc
}

func (sc *SourceCollection) Close() error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if sc.values == nil {
		return errors.New("sources are already closed")
	}

	errs := make([]error, 0, len(sc.values))

	for _, e := range sc.values {
		if err := e.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	sc.values = nil

	if len(errs) > 0 {
		return core.Errors(errs...)
	}

	return nil
}

func (sc *SourceCollection) Size() int {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return len(sc.values)
}

func (sc *SourceCollection) Get(idx int) (Source, error) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	if len(sc.values) <= idx {
		return nil, core.ErrNotFound
	}

	return sc.values[idx], nil
}

func (sc *SourceCollection) Add(source Source) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.values = append(sc.values, source)
}

func (sc *SourceCollection) Remove(source Source) bool {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	idx := -1

	for i, current := range sc.values {
		if current == source {
			idx = i
			break
		}
	}

	if idx > -1 {
		sc.values = append(sc.values[:idx], sc.values[idx+1:]...)
	}

	return idx > -1
}
