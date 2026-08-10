package sdk

import (
	"sync/atomic"
)

// IDGenerator allocates monotonically increasing opaque resource IDs.
type IDGenerator struct {
	next atomic.Uint64
}

// Next returns the next resource ID.
func (g *IDGenerator) Next() uint64 {
	return g.next.Add(1)
}
