package runtime

import (
	"encoding/binary"
	"hash/fnv"
	"time"
)

// Duration is a signed span of time with nanosecond precision.
type Duration time.Duration

// ZeroDuration is the zero-length duration value.
const ZeroDuration Duration = 0

// NewDuration creates a Ferret duration from a Go duration.
func NewDuration(input time.Duration) Duration {
	return Duration(input)
}

func (d Duration) Type() Type {
	return TypeDuration
}

func (d Duration) String() string {
	return time.Duration(d).String()
}

func (d Duration) Unwrap() any {
	return time.Duration(d)
}

func (d Duration) Hash() uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(TypeDuration.Name()))

	var bytes [8]byte
	binary.LittleEndian.PutUint64(bytes[:], uint64(d))
	_, _ = h.Write(bytes[:])

	return h.Sum64()
}

func (d Duration) Copy() Value {
	return d
}
