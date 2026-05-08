package diagnostics

import "errors"

type memberDiagnostic interface {
	Label() string
	Note() string
	Hint() string
}

func memberRuntimeDiagnostic(err error) (memberDiagnostic, bool) {
	var accessErr *MemberAccessError
	if errors.As(err, &accessErr) {
		return accessErr, true
	}

	var mutationErr *MemberMutationError
	if errors.As(err, &mutationErr) {
		return mutationErr, true
	}

	return nil, false
}
