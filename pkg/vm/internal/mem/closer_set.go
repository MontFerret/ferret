package mem

import (
	"errors"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// CloserSet is an ordered, deduplicated collection of io.Closers.
// It is the shared mechanical primitive underlying DeferredClosers and Result.
// Deduplication uses ResourceKey: runtime.Resource values are deduplicated by
// ResourceID; plain io.Closer values are deduplicated by pointer identity.
type CloserSet struct {
	seen    map[ResourceKey]struct{}
	closers []io.Closer
}

// Add inserts a closer into the set if it is non-nil and not already present.
// If closer implements runtime.Resource its ResourceID is used for dedup so
// two distinct objects with the same ID are treated as the same resource.
func (s *CloserSet) Add(closer io.Closer) bool {
	if closer == nil {
		return false
	}

	var key ResourceKey
	if res, ok := closer.(runtime.Resource); ok {
		key = ResourceKey{ID: res.ResourceID()}
	} else {
		key = ResourceKey{Closer: closer}
	}

	if s.seen == nil {
		s.seen = make(map[ResourceKey]struct{})
	}

	if _, exists := s.seen[key]; exists {
		return false
	}

	s.seen[key] = struct{}{}
	s.closers = append(s.closers, closer)

	return true
}

// CloseAll closes every closer in insertion order, joining any errors, then
// resets the set.
func (s *CloserSet) CloseAll() error {
	var err error

	for _, closer := range s.closers {
		if closeErr := closer.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}

	s.Reset()

	return err
}

// ForEach iterates over closers in insertion order.
func (s *CloserSet) ForEach(fn func(io.Closer)) {
	if fn == nil {
		return
	}

	for _, closer := range s.closers {
		fn(closer)
	}
}

// Merge adds all closers from other into s (deduplicating), then resets other.
func (s *CloserSet) Merge(other *CloserSet) {
	if other == nil {
		return
	}

	for _, closer := range other.closers {
		s.Add(closer)
	}

	other.Reset()
}

// Reset clears the set.
func (s *CloserSet) Reset() {
	s.closers = nil
	s.seen = nil
}

// Len returns the number of closers in the set.
func (s *CloserSet) Len() int {
	return len(s.closers)
}
