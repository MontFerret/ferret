package runtime

import (
	"encoding/json"
	"fmt"
)

// Value represents an interface of
// any type that needs to be used during runtime
type Value interface {
	// TODO: Remove Marshaler and introduce a runtime Serializer
	json.Marshaler
	fmt.Stringer
	Hashable
	// Copy returns a shallow copy of the value.
	Copy() Value
}
