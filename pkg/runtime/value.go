package runtime

import (
	"fmt"
)

// Value represents an interface of
// any type that needs to be used during runtime
type Value interface {
	fmt.Stringer
	Hashable
	// Copy returns a shallow copy of the value.
	Copy() Value
}
