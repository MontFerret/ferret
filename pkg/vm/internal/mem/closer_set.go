package mem

import (
	"errors"
	"io"
)

// CloserSet is an ordered, deduplicated collection of comparable io.Closers.
// It is the shared mechanical primitive underlying DeferredClosers and Result.
type CloserSet struct {
	seen    map[io.Closer]struct{}
	closers []io.Closer
}

// Add inserts a closer into the set if it is non-nil, comparable, and not
// already present. Returns true if the closer was actually added.
// Add inserts a closer into the set if it is non-nil and not already present.
// All closers are pointer-comparable by construction (either native pointer
// types or *ManagedResource), so no reflect-based comparability check is needed.
func (s *CloserSet) Add(closer io.Closer) bool {
	if closer == nil {
		return false
	}

	if s.seen == nil {
		s.seen = make(map[io.Closer]struct{})
	}

	if _, exists := s.seen[closer]; exists {
		return false
	}

	s.seen[closer] = struct{}{}
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
