package eval

import (
	"github.com/mafredri/cdp/protocol/runtime"
)

type RemoteType string

// List of possible remote types
const (
	UnknownType     RemoteType = ""
	NullType        RemoteType = "null"
	UndefinedType   RemoteType = "undefined"
	ArrayType       RemoteType = "array"
	NodeType        RemoteType = "node"
	RegexpType      RemoteType = "regexp"
	DateType        RemoteType = "date"
	MapType         RemoteType = "map"
	SetType         RemoteType = "set"
	WeakMapType     RemoteType = "weakmap"
	WeakSetType     RemoteType = "weakset"
	IteratorType    RemoteType = "iterator"
	GeneratorType   RemoteType = "generator"
	ErrorType       RemoteType = "error"
	ProxyType       RemoteType = "proxy"
	PromiseType     RemoteType = "promise"
	TypedArrayType  RemoteType = "typedarray"
	ArrayBufferType RemoteType = "arraybuffer"
	DataViewType    RemoteType = "dataview"
)

func ToRemoteType(ref runtime.RemoteObject) RemoteType {
	var subtype string

	if ref.Subtype != nil {
		subtype = *ref.Subtype
	}

	return RemoteType(subtype)
}
