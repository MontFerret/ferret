package types

import "github.com/MontFerret/ferret/pkg/runtime/core"

var (
	None    = core.NewType("none")
	Boolean = core.NewType("boolean")
	Int     = core.NewType("int")
	Float   = core.NewType("float")
	String  = core.NewType("string")
	Date    = core.NewType("date")
	Array   = core.NewType("array")
	Object  = core.NewType("object")
	Binary  = core.NewType("binary")
)
