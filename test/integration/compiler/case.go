package compiler_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/asm"
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/test/base/assert"
	"github.com/MontFerret/ferret/v2/test/base/compilation"
	"github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
)

func ByteCodeCase(expression string, expected []bytecode.Instruction, desc ...string) UseCase {
	return compilation.ByteCodeSpec(expression, expected, desc...)
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

			convey.Convey(useCase.Expression, t, func() {
				actual, err := c.Compile(file.NewSource(name, useCase.Expression))

				if !spec.ArePtrsEqual(useCase.PostCompilation, assert.IsCompilationError) {
					convey.So(err, convey.ShouldBeNil)
				} else {
					convey.So(err, assert.IsCompilationError, useCase.Expected)

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

				for _, assertion := range useCase.PostRun {
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
