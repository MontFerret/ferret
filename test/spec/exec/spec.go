package exec

import (
	"github.com/MontFerret/ferret/v2/test/spec"
	"github.com/MontFerret/ferret/v2/test/spec/assert"
)

func S(expression string, expected any, desc ...string) spec.Spec {
	return spec.New(expression, desc...).Expect().Exec(assert.ShouldEqual, expected)
}

func Nil(expression string, desc ...string) spec.Spec {
	return spec.New(expression, desc...).Expect().Exec(assert.ShouldBeNil)
}

func Error(expression string, desc ...string) spec.Spec {
	return spec.New(expression, desc...).Expect().ExecError(assert.ShouldNotBeNil)
}

func ErrorAs(expression string, expected error, desc ...string) spec.Spec {
	return spec.New(expression, desc...).Expect().ExecError(assert.ShouldBeError, expected)
}

func ErrorStr(expression string, expected string, desc ...string) spec.Spec {
	return spec.New(expression, desc...).Expect().ExecError(assert.ShouldBeError, expected)
}

func Object(expression string, expected map[string]any, desc ...string) spec.Spec {
	s := spec.New(expression, desc...).Expect().Exec(assert.ShouldEqualJSON, expected)
	s.Exec.RawOutput = true
	return s
}

func Array(expression string, expected []any, desc ...string) spec.Spec {
	s := spec.New(expression, desc...).Expect().Exec(assert.ShouldEqualJSON, expected)
	s.Exec.RawOutput = true
	return s
}

func Items(expression string, expected ...any) spec.Spec {
	return spec.New(expression).Expect().Exec(assert.ShouldHaveSameItems, expected)
}

func Fn(expression string, assertion assert.Unary, desc ...string) spec.Spec {
	return spec.New(expression, desc...).Expect().Exec(assert.NewUnaryAssertion(assertion))
}

func JSON(expression string, expected string, desc ...string) spec.Spec {
	s := spec.New(expression, desc...).Expect().Exec(assert.ShouldEqualJSON, expected)
	s.Exec.RawOutput = true
	return s
}
