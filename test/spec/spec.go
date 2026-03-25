package spec

import (
	"strings"
)

type (
	ExecInfo struct {
		Outcomes
		RawOutput bool
	}

	Spec struct {
		Base    BaseSpec
		Compile Outcomes
		Exec    ExecInfo
	}
)

func New(expression string, desc ...string) Spec {
	return Spec{
		Base: BaseSpec{
			Expression:  expression,
			Description: strings.Join(desc, " "),
		},
	}
}

func (s Spec) Expect() ExpectationBuilder[Spec] {
	return NewExpectationBuilder(&s, func(s *Spec) *Spec {
		return s
	})
}

func (s Spec) Suffix(suffix string) Spec {
	s.Base = s.Base.Suffix(suffix)

	return s
}

func (s Spec) Skip(reason ...string) Spec {
	s.Base = s.Base.Skip(reason...)

	return s
}

func (s Spec) Debug() Spec {
	s.Base.DebugOutput = true

	return s
}

func (s Spec) ExecRaw() Spec {
	s.Exec.RawOutput = true

	return s
}

func (s Spec) SuiteName(suite string) string {
	return s.Base.SuiteName(suite)
}

func (s Spec) String() string {
	return s.Base.String()
}
