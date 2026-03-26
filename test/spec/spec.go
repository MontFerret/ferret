package spec

import (
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type Spec struct {
	Base    BaseSpec
	Compile CompileInfo
	Exec    ExecInfo
}

func New(expression string, desc ...string) Spec {
	return Spec{
		Base: NewBaseSpec(expression, desc...),
	}
}

func NewWith(input Input, desc ...string) Spec {
	return Spec{
		Base: NewBaseSpecWith(input, desc...),
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

func (s Spec) Raw() Spec {
	s.Exec.RawOutput = true

	return s
}

func (s Spec) Env(o ...vm.EnvironmentOption) Spec {
	s.Exec.Env = append(s.Exec.Env, o...)

	return s
}

func (s Spec) VM(o ...vm.Option) Spec {
	s.Exec.VM = append(s.Exec.VM, o...)

	return s
}

func (s Spec) Merge(other Spec) Spec {
	out := s
	out.Base = s.Base.Merge(other.Base)
	out.Compile = s.Compile.Merge(other.Compile)
	out.Exec = s.Exec.Merge(other.Exec)

	return out
}

func (s Spec) SuiteName(suite string) string {
	return s.Base.SuiteName(suite)
}

func (s Spec) String() string {
	return s.Base.String()
}
