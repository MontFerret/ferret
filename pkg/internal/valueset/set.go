package valueset

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

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

func (s *Set) Add(ctx context.Context, value runtime.Value) (bool, error) {
	hash := value.Hash()
	first, exists := s.firstByHash[hash]

	if !exists {
		s.firstByHash[hash] = value
		s.count++

		return true, nil
	}

	equal, err := runtime.EqualValues(ctx, first, value)
	if err != nil {
		return false, err
	}

	if equal {
		return false, nil
	}

	for _, existing := range s.collisions[hash] {
		equal, err = runtime.EqualValues(ctx, existing, value)
		if err != nil {
			return false, err
		}

		if equal {
			return false, nil
		}
	}

	if s.collisions == nil {
		s.collisions = make(map[uint64][]runtime.Value)
	}

	s.collisions[hash] = append(s.collisions[hash], value)
	s.count++

	return true, nil
}

func (s *Set) Len() int {
	return s.count
}
