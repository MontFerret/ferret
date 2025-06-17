package bytecode_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/vm"
	"github.com/MontFerret/ferret/test/integration/base"
)

func Case(expression string, expected *vm.Program, desc ...string) UseCase {
	return NewCase(expression, expected, convey.ShouldEqual, desc...)
}

func SkipCase(expression string, expected *vm.Program, desc ...string) UseCase {
	return Skip(Case(expression, expected, desc...))
}

func ByteCodeCase(expression string, expected []vm.Instruction, desc ...string) UseCase {
	return NewCase(expression, &vm.Program{
		Bytecode: expected,
	}, ShouldEqualBytecode, desc...)
}

func SkipByteCodeCase(expression string, expected []vm.Instruction, desc ...string) UseCase {
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
				actual, err := c.Compile(useCase.Expression)

				if !base.ArePtrsEqual(useCase.PreAssertion, base.ShouldBeCompilationError) {
					convey.So(err, convey.ShouldBeNil)
				} else {
					convey.So(err, convey.ShouldBeError)

					return
				}

				println("")
				println("Actual:")
				println(actual.String())

				convey.So(err, convey.ShouldBeNil)

				for _, assertion := range useCase.Assertions {
					convey.So(actual, assertion, useCase.Expected)
				}
			})
		})
	}
}

func RunUseCases(t *testing.T, useCases []UseCase) {
	RunUseCasesWith(t, compiler.New(), useCases)
}

func Disassembly(instr []string, opcodes ...vm.Opcode) string {
	var disassembly string

	for i := 0; i < len(instr); i++ {
		disassembly += fmt.Sprintf("%d: [%d] %s\n", i, opcodes[i], instr[i])
	}

	return disassembly
}
