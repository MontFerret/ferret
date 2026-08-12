package diagnostics

import "errors"

type structuredRuntimeDiagnostic interface {
	Label() string
	Note() string
	Hint() string
}

func runtimeDiagnosticDetails(err error) (structuredRuntimeDiagnostic, bool) {
	var accessErr *MemberAccessError
	if errors.As(err, &accessErr) {
		return accessErr, true
	}

	var mutationErr *MemberMutationError
	if errors.As(err, &mutationErr) {
		return mutationErr, true
	}

	var deletionErr *MemberDeletionError
	if errors.As(err, &deletionErr) {
		return deletionErr, true
	}

	return nil, false
}
