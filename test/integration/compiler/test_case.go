package compiler_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/file"

	"github.com/MontFerret/ferret/v2/pkg/asm"

	"github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/test/integration/base"
)

func Case(expression string, expected *bytecode.Program, desc ...string) UseCase {
	return NewCase(expression, expected, convey.ShouldEqual, desc...)
}

func SkipCase(expression string, expected *bytecode.Program, desc ...string) UseCase {
	return Skip(Case(expression, expected, desc...))
}

func ByteCodeCase(expression string, expected []bytecode.Instruction, desc ...string) UseCase {
	return NewCase(expression, &bytecode.Program{
		Bytecode: expected,
	}, ShouldEqualBytecode, desc...)
}

func ErrorCase(expression string, expected base.ExpectedError, desc ...string) UseCase {
	uc := NewCase(expression, &expected, nil, desc...)
	uc.PreAssertion = base.ShouldBeCompilationError
	uc.Assertions = []convey.Assertion{
		func(actual any, expected ...any) string {
			return "expected compilation error"
		},
	}

	return uc
}

func SkipErrorCase(expression string, expected base.ExpectedError, desc ...string) UseCase {
	return Skip(ErrorCase(expression, expected, desc...))
}

func MultiErrorCase(expression string, expected base.ExpectedMultiError, desc ...string) UseCase {
	uc := NewCase(expression, &expected, nil, desc...)
	uc.PreAssertion = base.ShouldBeCompilationError
	uc.Assertions = []convey.Assertion{
		func(actual any, expected ...any) string {
			return "expected compilation error"
		},
	}

	return uc
}

func SkipMultiErrorCase(expression string, expected base.ExpectedMultiError, desc ...string) UseCase {
	return Skip(MultiErrorCase(expression, expected, desc...))
}

func SkipByteCodeCase(expression string, expected []bytecode.Instruction, desc ...string) UseCase {
	return Skip(ByteCodeCase(expression, expected, desc...))
}

func RunUseCasesWith(t *testing.T, c *compiler.Compiler, useCases []UseCase) {
	for _, useCase := range useCases {
		name := useCase.Description

		if useCase.Description == "" {
			name = strings.TrimSpace(useCase.Expression)
		}

		name = strings.Replace(name, "\n", " ", -1)
		name = strings.Replace(name, "\t", " ", -1)
		// Replace multiple spaces with a single space
		name = strings.Join(strings.Fields(name), " ")
		skip := useCase.Skip

		t.Run("Bytecode Test: "+name, func(t *testing.T) {
			if skip {
				t.Skip()

				return
			}

			convey.Convey(useCase.Expression, t, func() {
				actual, err := c.Compile(file.NewSource(name, useCase.Expression))

				if !base.ArePtrsEqual(useCase.PreAssertion, base.ShouldBeCompilationError) {
					convey.So(err, convey.ShouldBeNil)
				} else {
					convey.So(err, base.ShouldBeCompilationError, useCase.Expected)

					return
				}

				println("")
				println("Actual:")
				println(asm.Disassemble(actual))

				if p, ok := useCase.Expected.(*bytecode.Program); ok {
					println("Expected:")
					println(asm.Disassemble(p))
				}

				convey.So(err, convey.ShouldBeNil)

				for _, assertion := range useCase.Assertions {
					convey.So(actual, assertion, useCase.Expected)
				}
			})
		})
	}
}

func RunUseCases(t *testing.T, useCases []UseCase) {
	RunUseCasesWith(t, compiler.New(compiler.WithOptimizationLevel(compiler.O0)), useCases)
}

func Disassembly(instr []string, opcodes ...bytecode.Opcode) string {
	var disassembly string

	for i := 0; i < len(instr); i++ {
		disassembly += fmt.Sprintf("%d: [%d] %s\n", i, opcodes[i], instr[i])
	}

	return disassembly
}
