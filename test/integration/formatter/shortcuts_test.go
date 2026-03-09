package formatter_test

import (
	"github.com/MontFerret/ferret/v2/pkg/formatter"
	"github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/test/integration/base"
)

var (
	WithParam = base.WithParam
)

type UseCase struct {
	base.TestCase

	Options []formatter.Option
}

func NewCase(expression string, expected any, assertion convey.Assertion, desc ...string) UseCase {
	return UseCase{
		TestCase: base.NewCase(expression, expected, assertion, desc...),
	}
}

func Debug(uc UseCase) UseCase {
	uc.DebugOutput = true

	return uc
}

func Skip(uc UseCase) UseCase {
	uc.Skip = true

	return uc
}

func Options(uc UseCase, options ...formatter.Option) UseCase {
	uc.Options = options

	return uc
}
