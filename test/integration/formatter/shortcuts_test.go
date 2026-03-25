package formatter_test

import (
	"github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/formatter"

	"github.com/MontFerret/ferret/v2/test/integration/base"
)

var (
	WithParam = spec.WithParam
)

type UseCase struct {
	Options []formatter.Option
	spec.Spec
}

func NewCase(expression string, expected any, assertion convey.Assertion, desc ...string) UseCase {
	return UseCase{
		Spec: spec.NewSpec(expression, expected, assertion, desc...),
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
