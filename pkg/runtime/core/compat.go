package core

import "github.com/MontFerret/ferret/pkg/runtime"

// This file is used to provide backward compatibility for the Ferret runtime.
type (
	Value     = runtime.Value
	Boolean   = runtime.Boolean
	String    = runtime.String
	Int       = runtime.Int
	Float     = runtime.Float
	Array     = runtime.Array
	Object    = runtime.Object
	DateTime  = runtime.DateTime
	Binary    = runtime.Binary
	Iterable  = runtime.Iterable
	List      = runtime.List
	Map       = runtime.Map
	Function  = runtime.Function
	Namespace = runtime.Namespace
)

var (
	None               = runtime.None
	EmptyString        = runtime.EmptyString
	False              = runtime.False
	True               = runtime.True
	NewString          = runtime.NewString
	NewStringFromRunes = runtime.NewStringFromRunes
	NewInt             = runtime.NewInt
	NewFloat           = runtime.NewFloat
	NewArray           = runtime.NewArray
	NewArrayWith       = runtime.NewArrayWith
	NewObject          = runtime.NewObject
	NewDateTime        = runtime.NewDateTime
	NewCurrentDateTime = runtime.NewCurrentDateTime
	NewBinary          = runtime.NewBinary
	NewBoolean         = runtime.NewBoolean

	ForEach = runtime.ForEach

	NewFunctions        = runtime.NewFunctions
	NewFunctionsFromMap = runtime.NewFunctionsFromMap

	ValidateArgs  = runtime.ValidateArgs
	Error         = runtime.Error
	TypeError     = runtime.TypeError
	CompareValues = runtime.CompareValues
	Reflect       = runtime.Reflect

	AssertBinary   = runtime.AssertBinary
	AssertDateTime = runtime.AssertDateTime
	AssertInt      = runtime.AssertInt
	AssertFloat    = runtime.AssertFloat
	AssertList     = runtime.AssertList
	AssertMap      = runtime.AssertMap
	AssertString   = runtime.AssertString

	Parse            = runtime.Parse
	ToBoolean        = runtime.ToBoolean
	Random           = runtime.Random
	RandomDefault    = runtime.RandomDefault
	NumberBoundaries = runtime.NumberBoundaries

	MaxArgs = runtime.MaxArgs
)

func ValidateType(value Value, expectedType string) error {
	if runtime.Reflect(value) != expectedType {
		return runtime.TypeError(value, expectedType)
	}

	return nil
}
