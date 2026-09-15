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
	entropy func() [16]byte
	pcg     rand.PCG
}

// New creates an independent source that reads operating-system entropy on its
// first actual draw. Construction and equal integer bounds do not read entropy.
func New() *Source {
	return newSource(func() [16]byte {
		var seed [16]byte
		// crypto/rand.Read fills the buffer and never returns an error.
		_, _ = cryptorand.Read(seed[:])

		return seed
	})
}

// NewSeed creates a reproducible sequence. Exact sequences are stable within a
// Ferret version, not a compatibility guarantee between versions.
func NewSeed(seed int64) *Source {
	src := &Source{}
	src.pcg.Seed(uint64(seed), 0)

	return src
}

// The private entropy boundary allows deterministic initialization tests without
// replacing crypto/rand.Reader or exposing host RNG injection. Returning an array
// avoids passing a buffer pointer through the callback. The entropy reader may
// still cause the buffer to escape, depending on the platform and instrumentation.
func newSource(entropy func() [16]byte) *Source {
	return &Source{entropy: entropy}
}

// Float64 draws a value in [0, 1).
func (s *Source) Float64() float64 {
	random := s.generator()

	return random.Float64()
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
	random := s.generator()
	var offset uint64
	if width == 0 {
		offset = random.Uint64()
	} else {
		offset = random.Uint64N(width)
	}

	return int64(uint64(min) + offset)
}

// Bool draws a uniformly distributed Boolean using one source value.
func (s *Source) Bool() bool {
	random := s.generator()

	return random.Uint64()&1 != 0
}

// LegacyFloat64 preserves the historical rand(max, min) rounded calculation,
// including reversed, fractional, and non-finite bounds. It always draws once.
func (s *Source) LegacyFloat64(max, min float64) float64 {
	return math.Floor(s.Float64()*(max-min+1)) + min
}

func (s *Source) generator() rand.Rand {
	if s.entropy != nil {
		s.initialize()
	}

	return *rand.New(&s.pcg)
}

// Keep the one-time entropy work separate from the initialized draw path.
func (s *Source) initialize() {
	seed := s.entropy()
	s.pcg.Seed(binary.LittleEndian.Uint64(seed[:8]), binary.LittleEndian.Uint64(seed[8:]))
	s.entropy = nil
}
