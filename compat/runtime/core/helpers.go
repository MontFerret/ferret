package core

import (
	"fmt"
	"hash/fnv"
	"io"
)

// simpleType is a pure compat Type created by NewType.
type simpleType struct {
	name string
}

func (t *simpleType) ID() int64 {
	h := fnv.New64a()
	h.Write([]byte(t.name))
	return int64(h.Sum64())
}

func (t *simpleType) String() string { return t.name }

func (t *simpleType) Equals(other Type) bool {
	if other == nil {
		return false
	}

	return t.name == other.String()
}

// NewType creates a new compat Type with the given name.
func NewType(name string) Type {
	return &simpleType{name: name}
}

// ValidateType checks that value matches at least one of the required types.
// Returns an error if it does not.
func ValidateType(value Value, required ...Type) error {
	if value == nil {
		if len(required) == 0 {
			return nil
		}

		return fmt.Errorf("expected %v, got nil", required)
	}

	actual := value.Type()

	for _, t := range required {
		if t == nil {
			continue
		}

		if actual.Equals(t) {
			return nil
		}
	}

	return fmt.Errorf("expected one of %v, got %v", required, actual)
}

// ValidateArgs checks that the number of arguments is within [minimum, maximum].
// Pass -1 for maximum to allow unlimited arguments.
func ValidateArgs(args []Value, minimum, maximum int) error {
	n := len(args)

	if n < minimum {
		return fmt.Errorf("expected at least %d argument(s), got %d", minimum, n)
	}

	if maximum >= 0 && n > maximum {
		return fmt.Errorf("expected at most %d argument(s), got %d", maximum, n)
	}

	return nil
}

// IsTypeOf reports whether value matches the given type.
func IsTypeOf(value Value, check Type) bool {
	if value == nil || check == nil {
		return false
	}

	return value.Type().Equals(check)
}

// NewRootScope creates a root-level execution scope (standalone, not connected to v2 VM).
// Provided for compilation compatibility with v1 code that uses core.NewRootScope().
func NewRootScope() (*Scope, CloseFunc) {
	root := &RootScope{
		disposables: make([]io.Closer, 0, 10),
	}
	scope := newScope(root, nil)

	return scope, root.Close
}

// PairValueType pairs a Value with a list of acceptable Types for batch validation.
type PairValueType struct {
	Value Value
	Types []Type
}

// NewPairValueType creates a PairValueType.
func NewPairValueType(value Value, types ...Type) PairValueType {
	return PairValueType{Value: value, Types: types}
}

// ValidateValueTypePairs validates each pair, returning the first error encountered.
func ValidateValueTypePairs(pairs ...PairValueType) error {
	for _, p := range pairs {
		if err := ValidateType(p.Value, p.Types...); err != nil {
			return err
		}
	}

	return nil
}
