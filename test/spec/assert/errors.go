package assert

import "github.com/MontFerret/ferret/v2/pkg/diagnostics"

type (
	ExpectedError struct {
		Message string
		Kind    diagnostics.Kind
		Hint    string
		Format  string
	}

	ExpectedMultiError struct {
		Errors []*ExpectedError
		Number int
	}
)
