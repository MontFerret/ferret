package spec

import "github.com/MontFerret/ferret/v2/test/spec/assert"

type (
	Expectation struct {
		Value          any
		ValueAssertion assert.Assertion
		Error          any
		ErrorAssertion assert.Assertion
	}

	ExpectationBuilder[T any] struct {
		ret  *T
		base func(*T) *Spec
	}
)

func NewExpectationBuilder[T any](ret *T, base func(*T) *Spec) ExpectationBuilder[T] {
	return ExpectationBuilder[T]{
		ret:  ret,
		base: base,
	}
}

func (b ExpectationBuilder[T]) Run(fn assert.Assertion, val ...any) T {
	spec := b.base(b.ret)

	if len(val) > 0 {
		spec.Run.Value = val[0]
	}

	spec.Run.ValueAssertion = fn
	spec.Run.Error = nil
	spec.Run.ErrorAssertion = nil
	return *b.ret
}

func (b ExpectationBuilder[T]) RunError(fn assert.Assertion, val ...any) T {
	spec := b.base(b.ret)

	if len(val) > 0 {
		spec.Run.Error = val[0]
	}

	spec.Run.ErrorAssertion = fn
	spec.Run.Value = nil
	spec.Run.ValueAssertion = nil
	return *b.ret
}

func (b ExpectationBuilder[T]) Compile(fn assert.Assertion, val ...any) T {
	spec := b.base(b.ret)

	if len(val) > 0 {
		spec.Run.Value = val[0]
	}

	spec.Compile.ValueAssertion = fn
	spec.Compile.Error = nil
	spec.Compile.ErrorAssertion = nil
	return *b.ret
}

func (b ExpectationBuilder[T]) CompileError(fn assert.Assertion, val ...any) T {
	spec := b.base(b.ret)

	if len(val) > 0 {
		spec.Run.Error = val[0]
	}

	spec.Compile.ErrorAssertion = fn
	spec.Compile.Value = nil
	spec.Compile.ValueAssertion = nil
	return *b.ret
}
