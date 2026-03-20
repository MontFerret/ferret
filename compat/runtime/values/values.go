// Package values provides v1-compatible concrete value types and helpers for the Ferret
// compatibility layer. It mirrors the github.com/MontFerret/ferret/pkg/runtime/values
// package from Ferret v1.
//
// The concrete value types (Int, Float, String, etc.) are type aliases for their v2
// runtime equivalents. This means values.Int IS runtime.Int — they share the same
// underlying type and can be used interchangeably.
//
// NOTE: Because these are type aliases for runtime.Value (not core.Value), they do NOT
// implement core.Value directly. To convert a concrete value to a core.Value, use
// core.WrapValue(). For example:
//
//	myInt := values.NewInt(42)
//	var cv core.Value = core.WrapValue(myInt) // wraps into compat Value
package values

import (
	"context"

	"github.com/MontFerret/ferret/v2/compat/runtime/core"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// --- Type aliases for concrete v2 types ---

type (
	Int      = runtime.Int
	Float    = runtime.Float
	String   = runtime.String
	Boolean  = runtime.Boolean
	DateTime = runtime.DateTime
	Binary   = runtime.Binary
	Array    = runtime.Array
	Object   = runtime.Object
)

// --- Constants and singletons ---

// None is the singleton nil/null value.
var (
	None = core.WrapValue(runtime.None)
)

const (
	True        = runtime.True
	False       = runtime.False
	ZeroInt     = runtime.ZeroInt
	ZeroFloat   = runtime.ZeroFloat
	EmptyString = runtime.EmptyString
	SpaceString = runtime.SpaceString
)

var (
	ZeroDateTime      = runtime.ZeroDateTime
	DefaultTimeLayout = runtime.DefaultTimeLayout
)

// --- Constructor functions ---

var (
	NewInt             = runtime.NewInt
	NewInt64           = runtime.NewInt64
	NewFloat           = runtime.NewFloat
	NaN                = runtime.NaN
	NewString          = runtime.NewString
	NewStringFromRunes = runtime.NewStringFromRunes
	NewBoolean         = runtime.NewBoolean
	NewBinary          = runtime.NewBinary
	NewBinaryFrom      = runtime.NewBinaryFrom
	NewDateTime        = runtime.NewDateTime
	NewCurrentDateTime = runtime.NewCurrentDateTime
	NewArray           = runtime.NewArray
	NewArray64         = runtime.NewArray64
	NewSizedArray      = runtime.NewSizedArray
	NewArrayWith       = runtime.NewArrayWith
	NewArrayOf         = runtime.NewArrayOf
	EmptyArray         = runtime.EmptyArray
	NewObject          = runtime.NewObject
	NewObjectOf        = runtime.NewObjectOf
)

// --- Parse functions ---

var (
	ParseInt          = runtime.ParseInt
	ParseFloat        = runtime.ParseFloat
	MustParseFloat    = runtime.MustParseFloat
	ParseString       = runtime.ParseString
	MustParseString   = runtime.MustParseString
	ParseBoolean      = runtime.ParseBoolean
	MustParseBoolean  = runtime.MustParseBoolean
	ParseDateTime     = runtime.ParseDateTime
	ParseDateTimeWith = runtime.ParseDateTimeWith
	MustParseDateTime = runtime.MustParseDateTime
)

// --- Utility functions (re-exports) ---

var (
	IsNaN    = runtime.IsNaN
	IsInf    = runtime.IsInf
	IsNumber = runtime.IsNumber
	IsScalar = runtime.IsScalar
	ValueOf  = runtime.ValueOf
	Hash     = runtime.Hash
	MapHash  = runtime.MapHash
)

// --- Helper functions returning core.Value ---

// Parse converts a Go value to a core.Value.
// This mirrors the v1 values.Parse function.
func Parse(input interface{}) core.Value {
	v, err := runtime.ValueOf(input)

	if err != nil {
		return None
	}

	return core.WrapValue(v)
}

// ToBoolean converts a core.Value to a Boolean.
func ToBoolean(input core.Value) Boolean {
	return runtime.ToBoolean(core.UnwrapValue(input))
}

// ToFloat converts a core.Value to a Float.
func ToFloat(input core.Value) Float {
	f, _ := runtime.ToFloat(context.Background(), core.UnwrapValue(input))
	return f
}

// ToInt converts a core.Value to an Int.
func ToInt(input core.Value) Int {
	i, _ := runtime.ToInt(context.Background(), core.UnwrapValue(input))
	return i
}

// ToIntDefault converts a core.Value to an Int, returning defaultValue on failure
// or if the result is not positive.
func ToIntDefault(input core.Value, defaultValue Int) Int {
	i, _ := runtime.ToIntDefault(context.Background(), core.UnwrapValue(input), defaultValue)
	return i
}

// ToString converts a core.Value to a String.
func ToString(input core.Value) String {
	return runtime.ToString(core.UnwrapValue(input))
}

// ToBinary converts a core.Value to a Binary.
func ToBinary(input core.Value) Binary {
	return runtime.ToBinary(core.UnwrapValue(input))
}

// ToArray converts a core.Value to an Array.
func ToArray(input core.Value) *Array {
	list, _ := runtime.ToList(context.Background(), core.UnwrapValue(input))
	if arr, ok := list.(*Array); ok {
		return arr
	}

	return EmptyArray()
}

// --- ObjectProperty and NewObjectWith compat wrapper ---

// ObjectProperty pairs a name with a runtime.Value, mirroring v1's values.ObjectProperty.
type ObjectProperty struct {
	Name  string
	Value runtime.Value
}

// NewObjectProperty creates a new ObjectProperty.
func NewObjectProperty(name string, value runtime.Value) *ObjectProperty {
	return &ObjectProperty{Name: name, Value: value}
}

// NewObjectWith creates a new Object from a list of ObjectProperty values.
// This mirrors the v1 API where NewObjectWith accepted variadic *ObjectProperty.
// The v2 runtime.NewObjectWith accepts a map instead.
func NewObjectWith(props ...*ObjectProperty) *Object {
	m := make(map[string]runtime.Value, len(props))
	for _, p := range props {
		m[p.Name] = p.Value
	}

	return runtime.NewObjectWith(m)
}
