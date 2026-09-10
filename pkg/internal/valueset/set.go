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

	index, err := s.find(ctx, hash, first, value)
	if err != nil {
		return false, err
	}

	if index >= 0 {
		return false, nil
	}

	if s.collisions == nil {
		s.collisions = make(map[uint64][]runtime.Value)
	}

	s.collisions[hash] = append(s.collisions[hash], value)
	s.count++

	return true, nil
}

// Contains verifies membership using runtime equality within a hash bucket.
func (s *Set) Contains(ctx context.Context, value runtime.Value) (bool, error) {
	hash := value.Hash()
	first, exists := s.firstByHash[hash]
	if !exists {
		return false, nil
	}

	index, err := s.find(ctx, hash, first, value)

	return index >= 0, err
}

// Remove deletes an equal value if present. Failed equality leaves the set intact.
func (s *Set) Remove(ctx context.Context, value runtime.Value) (bool, error) {
	hash := value.Hash()
	first, exists := s.firstByHash[hash]
	if !exists {
		return false, nil
	}

	index, err := s.find(ctx, hash, first, value)
	if err != nil || index < 0 {
		return false, err
	}

	bucket := s.collisions[hash]
	if index == 0 {
		if len(bucket) == 0 {
			delete(s.firstByHash, hash)
		} else {
			s.firstByHash[hash] = bucket[len(bucket)-1]
			bucket[len(bucket)-1] = nil
			bucket = bucket[:len(bucket)-1]
		}
	} else {
		copy(bucket[index-1:], bucket[index:])
		bucket[len(bucket)-1] = nil
		bucket = bucket[:len(bucket)-1]
	}

	if len(bucket) == 0 {
		delete(s.collisions, hash)
	} else {
		s.collisions[hash] = bucket
	}

	s.count--

	return true, nil
}

func (s *Set) Len() int {
	return s.count
}

// Index zero denotes the primary value; subsequent indexes select collisions.
func (s *Set) find(ctx context.Context, hash uint64, first, value runtime.Value) (int, error) {
	equal, err := runtime.EqualValues(ctx, first, value)
	if err != nil {
		return -1, err
	}

	if equal {
		return 0, nil
	}

	for index, existing := range s.collisions[hash] {
		equal, err = runtime.EqualValues(ctx, existing, value)
		if err != nil {
			return -1, err
		}

		if equal {
			return index + 1, nil
		}
	}

	return -1, nil
}
