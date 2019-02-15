package types

import "github.com/MontFerret/ferret/pkg/runtime/core"

var (
	None     = core.NewType("none")
	Boolean  = core.NewType("boolean")
	Int      = core.NewType("int")
	Float    = core.NewType("float")
	String   = core.NewType("string")
	DateTime = core.NewType("date_time")
	Array    = core.NewType("array")
	Object   = core.NewType("object")
	Binary   = core.NewType("binary")
)
