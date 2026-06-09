package valueset

import "github.com/MontFerret/ferret/v2/pkg/runtime"

// Set tracks distinct runtime values using Ferret equality semantics.
// Hashes select candidate buckets; equality always confirms membership.
// The set retains value references without copying, owning, or closing them.
type Set struct {
	firstByHash map[uint64]runtime.Value
	collisions  map[uint64][]runtime.Value
	count       int
}

func New(capacity int) *Set {
	if capacity < 0 {
		capacity = 0
	}

	return &Set{
		firstByHash: make(map[uint64]runtime.Value, capacity),
	}
}

func (s *Set) Add(value runtime.Value) bool {
	hash := value.Hash()
	first, exists := s.firstByHash[hash]

	if !exists {
		s.firstByHash[hash] = value
		s.count++

		return true
	}

	if runtime.CompareValues(first, value) == 0 {
		return false
	}

	for _, existing := range s.collisions[hash] {
		if runtime.CompareValues(existing, value) == 0 {
			return false
		}
	}

	if s.collisions == nil {
		s.collisions = make(map[uint64][]runtime.Value)
	}

	s.collisions[hash] = append(s.collisions[hash], value)
	s.count++

	return true
}

func (s *Set) Len() int {
	return s.count
}
