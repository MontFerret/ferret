package core

import (
	"encoding/json"
)

// Value represents an interface of
// any type that needs to be used during runtime
type Value interface {
	json.Marshaler
	String() string
	Unwrap() interface{}
	Hash() uint64
	Copy() Value
}
