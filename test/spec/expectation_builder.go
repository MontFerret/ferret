package spec

import "github.com/MontFerret/ferret/v2/test/spec/assert"

type ExpectationBuilder[T any] struct {
	ret  *T
	base func(*T) *Spec
}

func NewExpectationBuilder[T any](ret *T, base func(*T) *Spec) ExpectationBuilder[T] {
	return ExpectationBuilder[T]{
		ret:  ret,
		base: base,
	}
}

func (b ExpectationBuilder[T]) Exec(fn assert.Assertion, val ...any) T {
	spec := b.base(b.ret)

	if len(val) > 0 {
		spec.Exec.Result.Value = val[0]
	}

	spec.Exec.Result.Assertion = fn
	spec.Exec.Error = Expectation{}

	return *b.ret
}

func (b ExpectationBuilder[T]) ExecError(fn assert.Assertion, val ...any) T {
	spec := b.base(b.ret)

	if len(val) > 0 {
		spec.Exec.Error.Value = val[0]
	}

	spec.Exec.Error.Assertion = fn
	spec.Exec.Result = Expectation{}

	return *b.ret
}

func (b ExpectationBuilder[T]) Compile(fn assert.Assertion, val ...any) T {
	spec := b.base(b.ret)

	if len(val) > 0 {
		spec.Compile.Result.Value = val[0]
	}

	spec.Compile.Result.Assertion = fn
	spec.Compile.Error = Expectation{}

	return *b.ret
}

func (b ExpectationBuilder[T]) CompileError(fn assert.Assertion, val ...any) T {
	spec := b.base(b.ret)

	if len(val) > 0 {
		spec.Compile.Error.Value = val[0]
	}

	spec.Compile.Error.Assertion = fn
	spec.Compile.Result = Expectation{}

	return *b.ret
}
