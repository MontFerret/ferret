package base

import (
	. "github.com/smartystreets/goconvey/convey"
	"strings"
)

type UseCase struct {
	Expression   string
	Expected     any
	PreAssertion Assertion
	Assertions   []Assertion
	Description  string
	Skip         bool
	RawOutput    bool
}

func NewCase(expression string, expected any, assertion Assertion, desc ...string) UseCase {
	return UseCase{
		Expression:  expression,
		Expected:    expected,
		Assertions:  []Assertion{assertion},
		Description: strings.TrimSpace(strings.Join(desc, " ")),
	}
}

func Skip(uc UseCase) UseCase {
	uc.Skip = true
	return uc
}
