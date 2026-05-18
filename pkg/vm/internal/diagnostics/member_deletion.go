package diagnostics

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type MemberDeletionError struct {
	Target runtime.Type
	Access MemberAccessKind
	Member string
}

func MemberDeletionErrorOf(target runtime.Value, access MemberAccessKind, member runtime.Value) error {
	return &MemberDeletionError{
		Target: runtime.TypeOf(target),
		Access: access,
		Member: memberAccessName(member),
	}
}

func (e *MemberDeletionError) Error() string {
	if e == nil {
		return ""
	}

	switch e.Access {
	case MemberAccessIndex:
		if e.Member == "" {
			return fmt.Sprintf("cannot delete index from %s", e.Target)
		}

		return fmt.Sprintf("cannot delete index %s from %s", e.Member, e.Target)
	case MemberAccessProperty:
		if e.Member == "" {
			return fmt.Sprintf("cannot delete property from %s", e.Target)
		}

		return fmt.Sprintf("cannot delete property %q from %s", e.Member, e.Target)
	default:
		if e.Member == "" {
			return fmt.Sprintf("cannot delete member from %s", e.Target)
		}

		return fmt.Sprintf("cannot delete member %q from %s", e.Member, e.Target)
	}
}

func (e *MemberDeletionError) Unwrap() error {
	return runtime.ErrInvalidType
}

func (e *MemberDeletionError) Label() string {
	if e == nil {
		return ""
	}

	switch e.Access {
	case MemberAccessIndex:
		if e.Member == "" {
			return "index cannot be deleted from this value"
		}

		return fmt.Sprintf("index %s cannot be deleted from this value", e.Member)
	case MemberAccessProperty:
		if e.Member == "" {
			return "property cannot be deleted from this value"
		}

		return fmt.Sprintf("property %q cannot be deleted from this value", e.Member)
	default:
		if e.Member == "" {
			return "member cannot be deleted from this value"
		}

		return fmt.Sprintf("member %q cannot be deleted from this value", e.Member)
	}
}

func (e *MemberDeletionError) Note() string {
	return e.Error()
}

func (e *MemberDeletionError) Hint() string {
	if e == nil {
		return ""
	}

	if runtime.IsSameType(e.Target, runtime.TypeNone) {
		return "Use optional chaining (?.) or check for None before deleting a member"
	}

	switch e.Access {
	case MemberAccessProperty:
		return "Ensure the value supports property deletion (for example, a mutable object)"
	case MemberAccessIndex:
		return "Ensure the value supports index deletion"
	default:
		return "Ensure the value supports member deletion"
	}
}
