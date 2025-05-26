package runtime

import (
	"encoding/json"
)

// Value represents an interface of
// any type that needs to be used during runtime
type Value interface {
	// TODO: Remove Marshaler and introduce a runtime Serializer
	json.Marshaler
	String() string
	Unwrap() interface{}
	Hash() uint64
	// TODO: Add context and return error. It needs to support values that rely on external storages.
	Copy() Value
}
