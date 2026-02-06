package optimization_test

import (
	j "encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/pkg/file"

	"github.com/MontFerret/ferret/pkg/asm"

	"github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/vm"
	"github.com/MontFerret/ferret/test/integration/base"
)

func Case(expression string, expected *vm.Program, desc ...string) UseCase {
	return UseCase{
		TestCase: base.NewCase(expression, expected, ShouldEqualBytecode, desc...),
	}
}

func SkipCase(expression string, expected *vm.Program, desc ...string) UseCase {
	return Skip(Case(expression, expected, desc...))
}

func ByteCodeCase(expression string, expected []vm.Instruction, desc ...string) UseCase {
	return UseCase{
		TestCase: base.NewCase(expression, &vm.Program{
			Bytecode: expected,
		}, ShouldEqualBytecode, desc...),
	}
}

func SkipByteCodeCase(expression string, expected []vm.Instruction, desc ...string) UseCase {
	return Skip(ByteCodeCase(expression, expected, desc...))
}

func RegistersCase(expression string, num int, output any, desc ...string) UseCase {
	return NewCase(expression, num, ShouldUseEqRegisters, Execution{
		Run:       true,
		Expected:  output,
		Assertion: convey.ShouldEqual,
	}, desc...)
}

func SkipAtMostRegistersCase(expression string, num int, output any, desc ...string) UseCase {
	return Skip(RegistersCase(expression, num, output, desc...))
}

func RegistersArrayCase(expression string, num int, output []any, desc ...string) UseCase {
	return NewCase(expression, num, ShouldUseEqRegisters, Execution{
		Run:       true,
		Expected:  output,
		Assertion: ShouldEqualJSONValue,
	}, desc...)
}

func SkipRegistersArrayCase(expression string, num int, output []any, desc ...string) UseCase {
	return Skip(RegistersArrayCase(expression, num, output, desc...))
}

func RegistersObjectCase(expression string, num int, output map[string]any, desc ...string) UseCase {
	return NewCase(expression, num, ShouldUseEqRegisters, Execution{
		Run:       true,
		Expected:  output,
		Assertion: ShouldEqualJSONValue,
	}, desc...)
}

func SkipRegistersObjectCase(expression string, num int, output map[string]any, desc ...string) UseCase {
	return Skip(RegistersObjectCase(expression, num, output, desc...))
}

func RunUseCasesWith(t *testing.T, c *compiler.Compiler, useCases []UseCase) {
	// Register standard library functions
	std := base.Stdlib()

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

				if p, ok := useCase.Expected.(*vm.Program); ok {
					println("Expected:")
					println(asm.Disassemble(p))
				}

				convey.So(err, convey.ShouldBeNil)

				for _, assertion := range useCase.Assertions {
					convey.So(actual, assertion, useCase.Expected)
				}

				if useCase.Execution.Run {
					options := []vm.EnvironmentOption{
						vm.WithFunctions(std),
					}

					assertion := useCase.Execution.Assertion
					expected := useCase.Execution.Expected
					out, err := base.Exec(actual, useCase.RawOutput, options...)

					convey.So(err, convey.ShouldBeNil)

					if base.ArePtrsEqual(assertion, convey.ShouldEqualJSON) {
						expectedJ, err := j.Marshal(expected)
						convey.So(err, convey.ShouldBeNil)
						convey.So(out, convey.ShouldEqualJSON, string(expectedJ))
					} else if base.ArePtrsEqual(assertion, base.ShouldHaveSameItems) {
						convey.So(out, base.ShouldHaveSameItems, expected)
					} else if base.ArePtrsEqual(assertion, convey.ShouldBeNil) {
						convey.So(out, convey.ShouldBeNil)
					} else {
						convey.So(out, assertion, expected)
					}
				}
			})
		})
	}
}

func RunUseCases(t *testing.T, level compiler.OptimizationLevel, useCases []UseCase) {
	RunUseCasesWith(t, compiler.New(compiler.WithOptimizationLevel(level)), useCases)
}

func Disassembly(instr []string, opcodes ...vm.Opcode) string {
	var disassembly string

	for i := 0; i < len(instr); i++ {
		disassembly += fmt.Sprintf("%d: [%d] %s\n", i, opcodes[i], instr[i])
	}

	return disassembly
}
