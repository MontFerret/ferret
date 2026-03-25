package exec

import (
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	"github.com/MontFerret/ferret/v2/test/spec/assert"
)

type Spec struct {
	Base       spec.Spec
	EnvOptions []vm.EnvironmentOption
}

func NewSpec(expression string, desc ...string) Spec {
	return Spec{
		Base: spec.New(expression, desc...),
	}
}

func (s Spec) Expect() spec.ExpectationBuilder[Spec] {
	return spec.NewExpectationBuilder(&s, func(s *Spec) *spec.Spec {
		return &s.Base
	})
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

func (s Spec) Options(options ...vm.EnvironmentOption) Spec {
	s.EnvOptions = options

	return s
}

func (s Spec) String() string {
	return s.Base.String()
}

func S(expression string, expected any, desc ...string) Spec {
	return NewSpec(expression, desc...).Expect().Run(assert.ShouldEqual, expected)
}

func Nil(expression string, desc ...string) Spec {
	return NewSpec(expression, desc...).Expect().Run(assert.ShouldBeNil)
}

func Error(expression string, desc ...string) Spec {
	return NewSpec(expression, desc...).Expect().RunError(assert.ShouldNotBeNil)
}

func ErrorAs(expression string, expected error, desc ...string) Spec {
	return NewSpec(expression, desc...).Expect().RunError(assert.ShouldBeError, expected)
}

func ErrorStr(expression string, expected string, desc ...string) Spec {
	return NewSpec(expression, desc...).Expect().RunError(assert.ShouldBeError, expected)
}

func Object(expression string, expected map[string]any, desc ...string) Spec {
	s := NewSpec(expression, desc...).Expect().Run(assert.ShouldEqualJSON, expected)
	s.Base.RawOutput = true
	return s
}

func Array(expression string, expected []any, desc ...string) Spec {
	s := NewSpec(expression, desc...).Expect().Run(assert.ShouldEqualJSON, expected)
	s.Base.RawOutput = true
	return s
}

func Items(expression string, expected ...any) Spec {
	return NewSpec(expression).Expect().Run(assert.ShouldHaveSameItems, expected)
}

func Fn(expression string, assertion assert.Unary, desc ...string) Spec {
	return NewSpec(expression, desc...).Expect().Run(assert.NewUnaryAssertion(assertion))
}

func JSON(expression string, expected string, desc ...string) Spec {
	s := NewSpec(expression, desc...).Expect().Run(assert.ShouldEqualJSON, expected)
	s.Base.RawOutput = true
	return s
}
