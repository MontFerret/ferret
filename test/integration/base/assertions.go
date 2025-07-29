package base

import (
	"fmt"

	"github.com/MontFerret/ferret/pkg/compiler"

	. "github.com/smartystreets/goconvey/convey"
)

type (
	ExpectedError struct {
		Message string
		Kind    compiler.ErrorKind
		Hint    string
	}

	ExpectedMultiError struct {
		Number int
		Errors []*ExpectedError
	}
)

func ArePtrsEqual(expected, actual any) bool {
	if expected == nil || actual == nil {
		return false
	}

	p1 := fmt.Sprintf("%v", expected)
	p2 := fmt.Sprintf("%v", actual)

	return p1 == p2
}

func ShouldHaveSameItems(actual any, expected ...any) string {
	wapper := expected[0].([]any)
	expectedArr := wapper[0].([]any)

	for _, item := range expectedArr {
		if err := ShouldContain(actual, item); err != "" {
			return err
		}
	}

	return ""
}

func assertExpectedError(actual *compiler.CompilationError, expected *ExpectedError) string {
	if actual == nil {
		return "expected a compilation error"
	}

	if expected.Kind != "" && actual.Kind != expected.Kind {
		return fmt.Sprintf("expected error kind %s, got %s", expected.Kind, actual.Kind)
	}

	if expected.Message != "" && actual.Message != expected.Message {
		return fmt.Sprintf("expected error message '%s', got '%s'", expected.Message, actual.Message)
	}

	if expected.Hint != "" && actual.Hint != expected.Hint {
		return fmt.Sprintf("expected error hint '%s', got '%s'", expected.Hint, actual.Hint)
	}

	return ""
}

func assertExpectedErrors(actual *compiler.MultiCompilationError, expected *ExpectedMultiError) string {
	if actual == nil {
		return "expected a multi compilation error"
	}

	if expected.Number > 0 && len(actual.Errors) != expected.Number {
		return fmt.Sprintf("expected %d errors, got %d", expected.Number, len(actual.Errors))
	}

	if len(expected.Errors) > 0 {
		for i, err := range actual.Errors {
			if i >= len(expected.Errors) {
				break
			}

			msg := assertExpectedError(err, expected.Errors[i])

			if msg != "" {
				return msg
			}
		}
	}

	return ""
}

func ShouldBeCompilationError(actual any, expected ...any) string {
	var msg string

	switch ex := expected[0].(type) {
	case *ExpectedError:
		err, ok := actual.(*compiler.CompilationError)

		if !ok {
			return "expected a compilation error"
		}

		msg = assertExpectedError(err, ex)

		if msg != "" {
			fmt.Println(err.Format())
		}

		break
	case *ExpectedMultiError:
		err, ok := actual.(*compiler.MultiCompilationError)

		if !ok {
			return "expected a multi compilation error"
		}

		msg = assertExpectedErrors(err, ex)

		if msg != "" {
			fmt.Println(err.Format())
		}
	default:
		msg = "expected a compilation error"
	}

	return msg
}
