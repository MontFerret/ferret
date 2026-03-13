package diagnostics

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type MemberAccessKind string

const (
	MemberAccessProperty MemberAccessKind = "property"
	MemberAccessIndex    MemberAccessKind = "index"
)

type MemberAccessError struct {
	Target runtime.Type
	Access MemberAccessKind
	Member string
}

func MemberAccessErrorOf(target runtime.Value, access MemberAccessKind, member runtime.Value) error {
	return &MemberAccessError{
		Target: runtime.TypeOf(target),
		Access: access,
		Member: memberAccessName(member),
	}
}

func (e *MemberAccessError) Error() string {
	if e == nil {
		return ""
	}

	switch e.Access {
	case MemberAccessIndex:
		if e.Member == "" {
			return fmt.Sprintf("cannot read index of %s", e.Target)
		}

		return fmt.Sprintf("cannot read index %s of %s", e.Member, e.Target)
	case MemberAccessProperty:
		if e.Member == "" {
			return fmt.Sprintf("cannot read property of %s", e.Target)
		}

		return fmt.Sprintf("cannot read property %q of %s", e.Member, e.Target)
	default:
		if e.Member == "" {
			return fmt.Sprintf("cannot read member of %s", e.Target)
		}

		return fmt.Sprintf("cannot read member %q of %s", e.Member, e.Target)
	}
}

func (e *MemberAccessError) Unwrap() error {
	return runtime.ErrInvalidType
}

func (e *MemberAccessError) Label() string {
	if e == nil {
		return ""
	}

	access := string(e.Access)
	if access == "" {
		access = "member"
	}

	if runtime.TypeName(e.Target) == "" {
		return fmt.Sprintf("%s access", access)
	}

	return fmt.Sprintf("%s access on %s", access, e.Target)
}

func (e *MemberAccessError) Hint() string {
	if e == nil {
		return ""
	}

	if runtime.IsSameType(e.Target, runtime.TypeNone) {
		return "Use optional chaining (?.) or check for None before accessing a member"
	}

	switch e.Access {
	case MemberAccessProperty:
		return "Ensure the value supports property access (for example, an object)"
	case MemberAccessIndex:
		return "Ensure the value supports index access (for example, an array)"
	default:
		return "Ensure the value supports member access"
	}
}

func memberAccessName(member runtime.Value) string {
	if member == nil || member == runtime.None {
		return ""
	}

	switch v := member.(type) {
	case runtime.String:
		return string(v)
	default:
		return v.String()
	}
}
