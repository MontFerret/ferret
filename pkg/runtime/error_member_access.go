package runtime

import "fmt"

type MemberAccessKind string

const (
	MemberAccessProperty MemberAccessKind = "property"
	MemberAccessIndex    MemberAccessKind = "index"
)

type MemberAccessError struct {
	Target Type
	Access MemberAccessKind
	Member string
}

func MemberAccessErrorOf(target Value, access MemberAccessKind, member Value) error {
	return &MemberAccessError{
		Target: TypeOf(target),
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
	return ErrInvalidType
}

func (e *MemberAccessError) Label() string {
	if e == nil {
		return ""
	}

	access := string(e.Access)
	if access == "" {
		access = "member"
	}

	if e.Target == "" {
		return fmt.Sprintf("%s access", access)
	}

	return fmt.Sprintf("%s access on %s", access, e.Target)
}

func (e *MemberAccessError) Hint() string {
	if e == nil {
		return ""
	}

	if e.Target == TypeNone {
		return "Use optional chaining (?.) or check for none before accessing a member"
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

func memberAccessName(member Value) string {
	if member == nil || member == None {
		return ""
	}

	switch v := member.(type) {
	case String:
		return string(v)
	default:
		return v.String()
	}
}
