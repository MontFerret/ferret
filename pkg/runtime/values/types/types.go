package types

import "github.com/MontFerret/ferret/pkg/runtime/core"

const coreNamespace = "runtime"

func newCoreType(name string) core.Type {
	return core.NewType(coreNamespace, name)
}

var (
	None       = newCoreType("none")
	Boolean    = newCoreType("boolean")
	Int        = newCoreType("int")
	Float      = newCoreType("float")
	String     = newCoreType("string")
	Regexp     = newCoreType("regexp")
	Range      = newCoreType("range")
	DateTime   = newCoreType("datetime")
	Array      = newCoreType("array")
	Object     = newCoreType("object")
	Binary     = newCoreType("binary")
	Boxed      = newCoreType("boxed")
	Measurable = newCoreType("measurable")
	Iterable   = newCoreType("iterable")
	Keyed      = newCoreType("keyed")
	Indexed    = newCoreType("indexed")
)
