package core

import "math/rand"

type (
	Type interface {
		ID() int64
		Namespace() string
		Name() string
		String() string
	}

	GenericType struct {
		id        int64
		namespace string
		name      string
	}
)

func NewType(namespace, name string) Type {
	return &GenericType{rand.Int63(), namespace, name}
}

func (b *GenericType) ID() int64 {
	return b.id
}

func (b *GenericType) Namespace() string {
	return b.namespace
}

func (b *GenericType) Name() string {
	return b.name
}

func (b *GenericType) String() string {
	return b.namespace + "." + b.name
}

// IsTypeOf return true when value's type
// is equal to check type.
// Returns false, otherwise.
func IsTypeOf(value Value, check Type) bool {
	return Reflect(value).ID() == check.ID()
}

// ValidateType checks the match of
// value's type and required types.
func ValidateType(value Value, required ...Type) error {
	var valid bool
	tid := Reflect(value).ID()

	for _, t := range required {
		if tid == t.ID() {
			valid = true
			break
		}
	}

	if !valid {
		return TypeError(value, required...)
	}

	return nil
}
