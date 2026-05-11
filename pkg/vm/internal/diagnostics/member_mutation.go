package diagnostics

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type MemberMutationError struct {
	Target runtime.Type
	Access MemberAccessKind
	Member string
}

func MemberMutationErrorOf(target runtime.Value, access MemberAccessKind, member runtime.Value) error {
	return &MemberMutationError{
		Target: runtime.TypeOf(target),
		Access: access,
		Member: memberAccessName(member),
	}
}

func (e *MemberMutationError) Error() string {
	if e == nil {
		return ""
	}

	switch e.Access {
	case MemberAccessIndex:
		if e.Member == "" {
			return fmt.Sprintf("cannot write index of %s", e.Target)
		}

		return fmt.Sprintf("cannot write index %s of %s", e.Member, e.Target)
	case MemberAccessProperty:
		if e.Member == "" {
			return fmt.Sprintf("cannot write property of %s", e.Target)
		}

		return fmt.Sprintf("cannot write property %q of %s", e.Member, e.Target)
	default:
		if e.Member == "" {
			return fmt.Sprintf("cannot write member of %s", e.Target)
		}

		return fmt.Sprintf("cannot write member %q of %s", e.Member, e.Target)
	}
}

func (e *MemberMutationError) Unwrap() error {
	return runtime.ErrInvalidType
}

func (e *MemberMutationError) Label() string {
	if e == nil {
		return ""
	}

	switch e.Access {
	case MemberAccessIndex:
		if e.Member == "" {
			return "index cannot be written to this value"
		}

		return fmt.Sprintf("index %s cannot be written to this value", e.Member)
	case MemberAccessProperty:
		if e.Member == "" {
			return "property cannot be written to this value"
		}

		return fmt.Sprintf("property %q cannot be written to this value", e.Member)
	default:
		if e.Member == "" {
			return "member cannot be written to this value"
		}

		return fmt.Sprintf("member %q cannot be written to this value", e.Member)
	}
}

func (e *MemberMutationError) Note() string {
	return e.Error()
}

func (e *MemberMutationError) Hint() string {
	if e == nil {
		return ""
	}

	switch e.Access {
	case MemberAccessProperty:
		return "Ensure the value supports property writes (for example, a mutable object)"
	case MemberAccessIndex:
		return "Ensure the value supports index writes (for example, a mutable array)"
	default:
		return "Ensure the value supports member writes"
	}
}
