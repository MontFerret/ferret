package runtime

import "io"

// Resource represents a runtime-managed entity that participates in VM
// ownership tracking and must be closable. ResourceID must be stable for the
// lifetime of the live resource and unique among concurrently tracked
// resources.
type Resource interface {
	io.Closer
	ResourceID() uint64
}
