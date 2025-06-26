package core

import "github.com/MontFerret/ferret/pkg/runtime"

// This file is used to provide backward compatibility for the Ferret runtime.
type (
	Type       = runtime.Type
	Value      = runtime.Value
	Boolean    = runtime.Boolean
	String     = runtime.String
	Int        = runtime.Int
	Float      = runtime.Float
	Array      = runtime.Array
	Object     = runtime.Object
	DateTime   = runtime.DateTime
	Binary     = runtime.Binary
	Iterable   = runtime.Iterable
	Iterator   = runtime.Iterator
	List       = runtime.List
	Map        = runtime.Map
	Keyed      = runtime.Keyed
	Indexed    = runtime.Indexed
	Cloneable  = runtime.Cloneable
	Measurable = runtime.Measurable
	Function   = runtime.Function
	Namespace  = runtime.Namespace

	PairValueType = runtime.PairValueType
)

var (
	NewType = runtime.NewType

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

	DefaultTimeLayout = runtime.DefaultTimeLayout
	IsNaN             = runtime.IsNaN
	IsInf             = runtime.IsInf
	ForEach           = runtime.ForEachIter

	NewFunctions        = runtime.NewFunctions
	NewFunctionsFromMap = runtime.NewFunctionsFromMap

	Error               = runtime.Error
	Errorf              = runtime.Errorf
	ErrInvalidOperation = runtime.ErrInvalidOperation
	ErrMissedArgument   = runtime.ErrMissedArgument
	ErrNotUnique        = runtime.ErrNotUnique
	ErrNotFound         = runtime.ErrNotFound
	ErrInvalidArgument  = runtime.ErrInvalidArgument

	ValidateType  = runtime.ValidateType
	ValidateArgs  = runtime.ValidateArgs
	TypeError     = runtime.TypeErrorOf
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
