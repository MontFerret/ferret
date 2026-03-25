package format

import (
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/formatter"
	"github.com/MontFerret/ferret/v2/test/spec"
)

type Spec struct {
	FormatOptions []formatter.Option
	Base          spec.BaseSpec
	Output        spec.Outcomes
}

func NewSpec(expression string, desc ...string) Spec {
	return Spec{
		Base: spec.New(expression, desc...).Base,
	}
}

func (s Spec) Expect() ExpectationBuilder {
	return NewExpectationBuilder(&s)
}

func (s Spec) Suffix(suffix string) Spec {
	s.Base = s.Base.Suffix(suffix)
	return s
}

func (s Spec) Debug() Spec {
	s.Base = s.Base.Debug()
	return s
}

func (s Spec) Skip(reason ...string) Spec {
	s.Base = s.Base.Skip(reason...)
	return s
}

func (s Spec) Options(options ...formatter.Option) Spec {
	s.FormatOptions = options

	return s
}

func (s Spec) SuiteName(suite string) string {
	return s.Base.SuiteName(suite)
}

func (s Spec) String() string {
	return s.Base.String()
}

func S(expression, expected string, desc ...string) Spec {
	normalized := expected

	if strings.HasSuffix(normalized, "\n") {
		if !strings.HasSuffix(normalized, "\n\n") {
			// ends with a single newline; add one more to make it two
			normalized += "\n"
		}
	} else {
		// no trailing newline; add two
		normalized += "\n\n"
	}

	return NewSpec(expression, desc...)
}
