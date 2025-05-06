package types

import "github.com/MontFerret/ferret/pkg/runtime"

var (
	None     = runtime.TypeNone
	Boolean  = runtime.TypeBoolean
	Int      = runtime.TypeInt
	Float    = runtime.TypeFloat
	String   = runtime.TypeString
	DateTime = runtime.TypeDateTime
	Array    = runtime.TypeList
	Object   = runtime.TypeMap
	Binary   = runtime.TypeBinary
)
