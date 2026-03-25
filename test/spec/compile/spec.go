package compile

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/test/spec"
	"github.com/MontFerret/ferret/v2/test/spec/assert"
)

type (
	OpcodeExpectation interface {
		OpcodeExistence | OpcodeCount
	}

	OpcodeExistence struct {
		Exists    []bytecode.Opcode
		NotExists []bytecode.Opcode
	}

	OpcodeCount struct {
		Count map[bytecode.Opcode]int
	}
)

func ByteCode(expression string, expected []bytecode.Instruction, desc ...string) spec.Spec {
	return spec.New(expression, desc...).Expect().Compile(ShouldHaveSameBytecode, &bytecode.Program{
		Bytecode: expected,
	})
}

func Opcode[T OpcodeExpectation](expression string, expectation T, desc ...string) spec.Spec {
	return spec.New(expression, desc...).Expect().Compile(ShouldContainOpcode, expectation)
}

func Registers(expression string, num int, desc ...string) spec.Spec {
	return spec.New(expression, desc...).Expect().Compile(ShouldHaveRegisters, num)
}

func Error(expression string, expected assert.ExpectedError, desc ...string) spec.Spec {
	return spec.New(expression, desc...).Expect().Compile(ShouldBeCompilationError, expected)
}

func MultiErrorSpec(expression string, expected assert.ExpectedMultiError, desc ...string) spec.Spec {
	return spec.New(expression, desc...).Expect().Compile(ShouldBeCompilationError, expected)
}
