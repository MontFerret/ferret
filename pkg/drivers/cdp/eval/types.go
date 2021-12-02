package eval

import (
	"github.com/mafredri/cdp/protocol/runtime"
)

type (
	RemoteType string

	RemoteObjectType string

	RemoteClassName string
)

// List of possible remote types
// "object", "function", "undefined", "string", "number", "boolean", "symbol", "bigint"
const (
	UnknownType   RemoteType = ""
	UndefinedType RemoteType = "undefined"
	StringType    RemoteType = "string"
	NumberType    RemoteType = "number"
	BooleanType   RemoteType = "boolean"
	SymbolType    RemoteType = "symbol"
	BigintType    RemoteType = "bigint"
	ObjectType    RemoteType = "object"
)

var remoteTypeMap = map[string]RemoteType{
	string(UndefinedType): UndefinedType,
	string(StringType):    StringType,
	string(NumberType):    NumberType,
	string(BooleanType):   BooleanType,
	string(SymbolType):    SymbolType,
	string(BigintType):    BigintType,
	string(ObjectType):    ObjectType,
}

// List of possible remote object types
const (
	UnknownObjectType     RemoteObjectType = ""
	NullObjectType        RemoteObjectType = "null"
	UndefinedObjectType   RemoteObjectType = "undefined"
	ArrayObjectType       RemoteObjectType = "array"
	NodeObjectType        RemoteObjectType = "node"
	RegexpObjectType      RemoteObjectType = "regexp"
	DateObjectType        RemoteObjectType = "date"
	MapObjectType         RemoteObjectType = "map"
	SetObjectType         RemoteObjectType = "set"
	WeakMapObjectType     RemoteObjectType = "weakmap"
	WeakSetObjectType     RemoteObjectType = "weakset"
	IteratorObjectType    RemoteObjectType = "iterator"
	GeneratorObjectType   RemoteObjectType = "generator"
	ErrorObjectType       RemoteObjectType = "error"
	ProxyObjectType       RemoteObjectType = "proxy"
	PromiseObjectType     RemoteObjectType = "promise"
	TypedArrayObjectType  RemoteObjectType = "typedarray"
	ArrayBufferObjectType RemoteObjectType = "arraybuffer"
	DataViewObjectType    RemoteObjectType = "dataview"
)

var remoteObjectTypeMap = map[string]RemoteObjectType{
	string(NullObjectType):        NullObjectType,
	string(UndefinedObjectType):   UndefinedObjectType,
	string(ArrayObjectType):       ArrayObjectType,
	string(NodeObjectType):        NodeObjectType,
	string(RegexpObjectType):      RegexpObjectType,
	string(DateObjectType):        DateObjectType,
	string(MapObjectType):         MapObjectType,
	string(SetObjectType):         SetObjectType,
	string(WeakMapObjectType):     WeakMapObjectType,
	string(WeakSetObjectType):     WeakSetObjectType,
	string(IteratorObjectType):    IteratorObjectType,
	string(GeneratorObjectType):   GeneratorObjectType,
	string(ErrorObjectType):       ErrorObjectType,
	string(ProxyObjectType):       ProxyObjectType,
	string(PromiseObjectType):     PromiseObjectType,
	string(TypedArrayObjectType):  TypedArrayObjectType,
	string(ArrayBufferObjectType): ArrayBufferObjectType,
	string(DataViewObjectType):    DataViewObjectType,
}

// List of supported remote classses
const (
	UnknownClassName  RemoteClassName = ""
	DocumentClassName RemoteClassName = "HTMLDocument"
)

var remoteClassNameMap = map[string]RemoteClassName{
	string(DocumentClassName): DocumentClassName,
}

func ToRemoteType(ref runtime.RemoteObject) RemoteType {
	remoteType, found := remoteTypeMap[ref.Type]

	if found {
		return remoteType
	}

	return UnknownType
}

func ToRemoteObjectType(ref runtime.RemoteObject) RemoteObjectType {
	if ref.Subtype != nil {
		remoteObjectType, found := remoteObjectTypeMap[*ref.Subtype]

		if found {
			return remoteObjectType
		}
	}

	return UnknownObjectType
}

func ToRemoteClassName(ref runtime.RemoteObject) RemoteClassName {
	if ref.ClassName != nil {
		remoteClassName, found := remoteClassNameMap[*ref.ClassName]

		if found {
			return remoteClassName
		}
	}

	return UnknownClassName
}
