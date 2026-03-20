// Package types provides v1-compatible type constants for the Ferret compatibility layer.
// It mirrors the github.com/MontFerret/ferret/pkg/runtime/values/types package from Ferret v1.
package types

import (
	"github.com/MontFerret/ferret/v2/compat/runtime/core"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Type constants wrapping v2 runtime types as compat core.Type values.
var (
	None     = core.WrapType(runtime.TypeNone)
	Boolean  = core.WrapType(runtime.TypeBoolean)
	Int      = core.WrapType(runtime.TypeInt)
	Float    = core.WrapType(runtime.TypeFloat)
	String   = core.WrapType(runtime.TypeString)
	DateTime = core.WrapType(runtime.TypeDateTime)
	Array    = core.WrapType(runtime.TypeArray)
	Object   = core.WrapType(runtime.TypeObject)
	Binary   = core.WrapType(runtime.TypeBinary)
)

// typeRank returns an ordinal rank for the given type name, matching v1 ordering.
func typeRank(t core.Type) int {
	if t == nil {
		return 0
	}

	switch t.String() {
	case "None":
		return 0
	case "Boolean":
		return 1
	case "Int":
		return 2
	case "Float":
		return 3
	case "String":
		return 4
	case "DateTime":
		return 5
	case "Binary":
		return 6
	case "Array":
		return 7
	case "Object":
		return 8
	default:
		return 9
	}
}

// Compare compares two types by their ordinal rank.
// Returns -1 if first < second, 0 if equal, 1 if first > second.
func Compare(first, second core.Type) int64 {
	a := typeRank(first)
	b := typeRank(second)

	if a == b {
		return 0
	}

	if a < b {
		return -1
	}

	return 1
}
