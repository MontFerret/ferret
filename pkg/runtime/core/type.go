package core

import (
	"github.com/pkg/errors"
	"math/rand"
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
