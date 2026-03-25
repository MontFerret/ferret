package optimization_test

import (
	j "encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/asm"
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/base/assert"
	"github.com/MontFerret/ferret/v2/test/base/compilation"
)

func Case(expression string, expected []bytecode.Instruction, desc ...string) UseCase {
	return UseCase{
		Spec: compilation.ByteCodeSpec(expression, expected, desc...),
	}
}

func SkipCase(expression string, expected []bytecode.Instruction, desc ...string) UseCase {
	return Skip(Case(expression, expected, desc...))
}

func OpcodeCase[T compilation.OpcodeExpectation](expression string, expectation T, output any, desc ...string) UseCase {
	return UseCase{
		Spec: compilation.OpcodeCase(expression, expectation, desc...),
		Execution: Execution{
			Run:       true,
			Expected:  output,
			Assertion: assert.ShouldEqual,
		},
	}
}

func SkipOpcodeCase[T compilation.OpcodeExpectation](expression string, expectation T, output any, desc ...string) UseCase {
	return Skip(OpcodeCase(expression, expectation, output, desc...))
}

func RegistersCase(expression string, expectation int, output any, desc ...string) UseCase {
	return UseCase{
		Spec: compilation.RegistersCase(expression, expectation, desc...),
		Execution: Execution{
			Run:       true,
			Expected:  output,
			Assertion: assert.ShouldEqual,
		},
	}
}

func SkipRegistersCase(expression string, expectation int, output any, desc ...string) UseCase {
	return Skip(RegistersCase(expression, expectation, output, desc...))
}

func RegistersCaseWith(expression string, expectation int, output any, params map[string]runtime.Value, desc ...string) UseCase {
	return UseCase{
		Spec: compilation.RegistersCase(expression, expectation, desc...),
		Execution: Execution{
			Run:       true,
			Expected:  output,
			Assertion: assert.ShouldEqual,
			Options:   []vm.EnvironmentOption{vm.WithParams(params)},
		},
	}
}

func SkipRegistersCaseWith(expression string, expectation int, output any, desc ...string) UseCase {
	return Skip(RegistersCase(expression, expectation, output, desc...))
}

func RegistersArrayCase(expression string, num int, output []any, desc ...string) UseCase {
	uc := RegistersCase(expression, num, output, desc...)
	uc.Execution.Assertion = assert.ShouldEqualJSON
	uc.RawOutput = true
	return uc
}

func SkipRegistersArrayCase(expression string, num int, output []any, desc ...string) UseCase {
	return Skip(RegistersArrayCase(expression, num, output, desc...))
}

func RegistersObjectCase(expression string, num int, output map[string]any, desc ...string) UseCase {
	uc := RegistersCase(expression, num, output, desc...)
	uc.Execution.Assertion = assert.ShouldEqualJSON
	uc.RawOutput = true
	return uc
}

func SkipRegistersObjectCase(expression string, num int, output map[string]any, desc ...string) UseCase {
	return Skip(RegistersObjectCase(expression, num, output, desc...))
}

func RunUseCasesWith(t *testing.T, c *compiler.Compiler, useCases []UseCase) {
	// Add standard library functions
	std := spec.Stdlib()

	for _, useCase := range useCases {
		name := useCase.Description

		if useCase.Description == "" {
			name = strings.TrimSpace(useCase.Expression)
		}

		name = strings.ReplaceAll(name, "\n", " ")
		name = strings.ReplaceAll(name, "\t", " ")
		// Replace multiple spaces with a single space
		name = strings.Join(strings.Fields(name), " ")
		skip := useCase.Skip

		t.Run("Bytecode Test: "+name, func(t *testing.T) {
			if skip {
				t.Skip()

				return
			}

			actual, err := c.Compile(file.NewSource(name, useCase.Expression))

			if !assert.IsSameAssertion(useCase.PostCompilation, assert.ShouldBeCompilationError) {
				assert.ShouldBeNil(t, err)
			} else {
				assert.ShouldBeCompilationError(t, err, useCase.Expected)

				return
			}

			println("")
			println("Actual:")
			println(asm.Disassemble(actual))

			if p, ok := useCase.Expected.(*bytecode.Program); ok {
				println("Expected:")
				println(asm.Disassemble(p))
			}

			assert.ShouldBeNil(t, err)

			useCase.PostRun(t, actual, useCase.Expected)

			if useCase.Execution.Run {
				options := []vm.EnvironmentOption{
					vm.WithNamespace(std),
					vm.WithFunctionsBuilder(spec.ForWhileHelpers()),
				}

				if len(useCase.Execution.Options) > 0 {
					options = append(options, useCase.Execution.Options...)
				}

				assertion := useCase.Execution.Assertion
				expected := useCase.Execution.Expected
				out, err := spec.Exec(actual, useCase.RawOutput, options...)

				_, ok := expected.(error)

				if ok {
					assert.ShouldNotBeNil(t, err)
					return
				}

				assert.ShouldBeNil(t, err)

				if assert.IsSameAssertion(assertion, assert.ShouldEqualJSON) || assert.IsSameAssertion(assertion, assert.ShouldEqualJSON) {
					expectedJ, err := j.Marshal(expected)
					assert.ShouldBeNil(t, err)
					assert.ShouldEqualJSON(t, out, string(expectedJ))
				} else if assert.IsSameAssertion(assertion, assert.ShouldHaveSameItems) {
					assert.ShouldHaveSameItems(t, out, expected)
				} else if assert.IsSameAssertion(assertion, assert.ShouldBeNil) {
					assert.ShouldBeNil(t, out)
				} else {
					assertion(t, out, expected)
				}
			}
		})
	}
}

func RunUseCases(t *testing.T, level compiler.OptimizationLevel, useCases []UseCase) {
	RunUseCasesWith(t, compiler.New(compiler.WithOptimizationLevel(level)), useCases)
}

func Disassembly(instr []string, opcodes ...bytecode.Opcode) string {
	var disassembly string

	for i := 0; i < len(instr); i++ {
		disassembly += fmt.Sprintf("%d: [%d] %s\n", i, opcodes[i], instr[i])
	}

	return disassembly
}
