package core

import (
	"math/rand"

	"github.com/pkg/errors"
)

// Type represents runtime type with id for quick type check
// and GetName for error messages

//revive:disable-next-line:redefines-builtin-id
type (
	Type interface {
		ID() int64
		String() string
		Equals(other Type) bool
	}

	BaseType struct {
		id   int64
		name string
	}
)

func NewType(name string) Type {
	return BaseType{rand.Int63(), name}
}

func (t BaseType) ID() int64 {
	return t.id
}

func (t BaseType) String() string {
	return t.name
}

func (t BaseType) Equals(other Type) bool {
	return t.id == other.ID()
}

// IsTypeOf return true when value's type
// is equal to check type.
// Returns false, otherwise.
func IsTypeOf(value Value, check Type) bool {
	return value.Type().ID() == check.ID()
}

// ValidateType checks the match of
// value's type and required types.
func ValidateType(value Value, required ...Type) error {
	var valid bool
	tid := value.Type().ID()

	for _, t := range required {
		if tid == t.ID() {
			valid = true
			break
		}
	}

	if !valid {
		return TypeError(value.Type(), required...)
	}

	return nil
}

// ValidateValueTypePairs validate pairs of
// Values and Types.
// Returns error when type didn't match
func ValidateValueTypePairs(pairs ...PairValueType) error {
	var err error

	for idx, pair := range pairs {
		err = ValidateType(pair.Value, pair.Types...)

		if err != nil {
			return errors.Errorf("pair %d: %v", idx, err)
		}
	}

	return nil
}

// PairValueType is a supporting
// structure that used in validateValueTypePairs.
type PairValueType struct {
	Value Value
	Types []Type
}

// NewPairValueType it's a shortcut for creating a new PairValueType.
//
// The common pattern of using PairValueType is:
// ```
//
//	pairs := []core.PairValueType{
//	    core.PairValueType{args[0], []core.Type{types.String}},               // go vet warning
//	    core.PairValueType{Value: args[1], Types: []core.Type{types.Binary}}, // too long
//	}
//
// ```
// With NewPairValueType there is no need to type `[]core.Type{...}` and code becomes
// more readable and maintainable.
//
// That is how the code above looks like with NewPairValueType:
// ```
//
//	pairs := []core.PairValueType{
//	    core.NewPairValueType(args[0], types.String),
//	    core.NewPairValueType(args[1], types.Binary),
//	}
//
// ```
func NewPairValueType(value Value, types ...Type) PairValueType {
	return PairValueType{
		Value: value,
		Types: types,
	}
}
