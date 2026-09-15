// Package rnd supplies Ferret's execution-scoped, non-cryptographic randomness.
// Sources must be used serially and are not suitable for security-sensitive data.
package rnd

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"math"
	"math/rand/v2"
)

// Source holds one pseudo-random sequence. It requires no cleanup or locking.
// Construct sources with New or NewSeed; do not copy a source after construction.
type Source struct {
	random rand.Rand
}

// New creates an independently seeded source using operating-system entropy.
func New() *Source {
	return newSource(func(seed []byte) {
		// crypto/rand.Read fills the buffer and never returns an error.
		_, _ = cryptorand.Read(seed)
	})
}

// NewSeed creates a reproducible sequence. Exact sequences are stable within a
// Ferret version, not a compatibility guarantee between versions.
func NewSeed(seed int64) *Source {
	return &Source{random: *rand.New(rand.NewPCG(uint64(seed), 0))}
}

// The private entropy boundary allows deterministic construction tests without
// replacing crypto/rand.Reader or exposing host RNG injection.
func newSource(fill func([]byte)) *Source {
	var seed [16]byte
	fill(seed[:])

	return &Source{random: *rand.New(rand.NewPCG(
		binary.LittleEndian.Uint64(seed[:8]),
		binary.LittleEndian.Uint64(seed[8:]),
	))}
}

// Float64 draws a value in [0, 1).
func (s *Source) Float64() float64 {
	return s.random.Float64()
}

// Int64 draws uniformly from [min, max]. Bounds must already satisfy min <= max.
// Equal bounds do not advance the sequence.
func (s *Source) Int64(min, max int64) int64 {
	if min == max {
		return min
	}

	// Unsigned subtraction includes intervals crossing zero without signed
	// overflow. A wrapped width of zero represents the complete int64 domain.
	width := uint64(max) - uint64(min) + 1
	var offset uint64
	if width == 0 {
		offset = s.random.Uint64()
	} else {
		offset = s.random.Uint64N(width)
	}

	return int64(uint64(min) + offset)
}

// Bool draws a uniformly distributed Boolean using one source value.
func (s *Source) Bool() bool {
	return s.random.Uint64()&1 != 0
}

// LegacyFloat64 preserves the historical rand(max, min) rounded calculation,
// including reversed, fractional, and non-finite bounds. It always draws once.
func (s *Source) LegacyFloat64(max, min float64) float64 {
	return math.Floor(s.Float64()*(max-min+1)) + min
}
