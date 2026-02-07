package base

import (
	"strings"

	. "github.com/smartystreets/goconvey/convey"
)

type TestCase struct {
	Expression   string
	Expected     any
	PreAssertion Assertion
	Assertions   []Assertion
	Description  string
	Skip         bool
	RawOutput    bool
	DebugOutput  bool
}

func NewCase(expression string, expected any, assertion Assertion, desc ...string) TestCase {
	return TestCase{
		Expression:  expression,
		Expected:    expected,
		Assertions:  []Assertion{assertion},
		Description: strings.TrimSpace(strings.Join(desc, " ")),
	}
}

func Skip(uc TestCase) TestCase {
	uc.Skip = true
	return uc
}

func Debug(useCase TestCase) TestCase {
	useCase.DebugOutput = true

	return useCase
}

func (tc TestCase) String() string {
	if tc.Description != "" {
		return strings.TrimSpace(tc.Description)
	}

	exp := strings.TrimSpace(tc.Expression)
	exp = strings.Replace(exp, "\n", " ", -1)
	exp = strings.Replace(exp, "\t", " ", -1)
	// Replace multiple spaces with a single space
	exp = strings.Join(strings.Fields(exp), " ")

	return exp
}
