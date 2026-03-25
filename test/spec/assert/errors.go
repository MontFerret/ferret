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

func CastExpectedError(err any) (*ExpectedError, bool) {
	switch e := err.(type) {
	case ExpectedError:
		return &e, true
	case *ExpectedError:
		return e, true
	default:
		return nil, false
	}
}

func CastExpectedMultiError(err any) (*ExpectedMultiError, bool) {
	switch e := err.(type) {
	case ExpectedMultiError:
		return &e, true
	case *ExpectedMultiError:
		return e, true
	default:
		return nil, false
	}
}
