package base

import (
	"fmt"

	"github.com/MontFerret/ferret/pkg/compiler"

	. "github.com/smartystreets/goconvey/convey"
)

type ExpectedError struct {
	Message string
	Kind    compiler.ErrorKind
}

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

func ShouldBeCompilationError(actual any, expected ...any) string {
	err, ok := actual.(*compiler.CompilationError)

	if !ok {
		return "expected a compilation error"
	}

	var msg string

	switch ex := expected[0].(type) {
	case *ExpectedError:
		if ex.Kind != "" {
			msg = ShouldEqual(err.Kind, ex.Kind)
		}

		if msg == "" {
			msg = ShouldEqual(err.Message, ex.Message)
		}

		break
	case string:
		msg = ShouldEqual(err.Message, ex)
	default:
		msg = "expected a compilation error"
	}

	if msg != "" {
		fmt.Println(err.Format())
	}

	return msg
}
