package base

import (
	"strings"

	convey "github.com/smartystreets/goconvey/convey"
)

type TestCase struct {
	Expression   string
	Expected     any
	PreAssertion convey.Assertion
	Assertions   []convey.Assertion
	Description  string
	Skip         bool
	RawOutput    bool
	DebugOutput  bool
}

func NewCase(expression string, expected any, assertion convey.Assertion, desc ...string) TestCase {
	return TestCase{
		Expression:  expression,
		Expected:    expected,
		Assertions:  []convey.Assertion{assertion},
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
	exp = strings.ReplaceAll(exp, "\n", " ")
	exp = strings.ReplaceAll(exp, "\t", " ")
	// Replace multiple spaces with a single space
	exp = strings.Join(strings.Fields(exp), " ")

	return exp
}
