package vm_test

import (
	j "encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/pkg/asm"

	"github.com/MontFerret/ferret/test/integration/base"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/vm"
)

func Case(expression string, expected any, desc ...string) UseCase {
	return NewCase(expression, expected, ShouldEqual, desc...)
}

func SkipCase(expression string, expected any, desc ...string) UseCase {
	return Skip(Case(expression, expected, desc...))
}

func CaseNil(expression string, desc ...string) UseCase {
	return NewCase(expression, nil, ShouldBeNil, desc...)
}

func SkipCaseNil(expression string, desc ...string) UseCase {
	return Skip(CaseNil(expression, desc...))
}

func DebugCaseNil(expression string, desc ...string) UseCase {
	return Debug(CaseNil(expression, desc...))
}

func CaseRuntimeError(expression string, desc ...string) UseCase {
	return NewCase(expression, nil, ShouldBeError, desc...)
}

func CaseRuntimeErrorAs(expression string, expected error, desc ...string) UseCase {
	return NewCase(expression, expected, ShouldBeError, desc...)
}

func SkipCaseRuntimeError(expression string, desc ...string) UseCase {
	return Skip(CaseRuntimeError(expression, desc...))
}

func DebugCaseRuntimeError(expression string, desc ...string) UseCase {
	return Debug(CaseRuntimeError(expression, desc...))
}

func SkipCaseRuntimeErrorAs(expression string, expected error, desc ...string) UseCase {
	return Skip(CaseRuntimeErrorAs(expression, expected, desc...))
}

func CaseCompilationError(expression string, desc ...string) UseCase {
	return NewCase(expression, nil, base.ShouldBeCompilationError, desc...)
}

func SkipCaseCompilationError(expression string, desc ...string) UseCase {
	return Skip(CaseCompilationError(expression, desc...))
}

func CaseObject(expression string, expected map[string]any, desc ...string) UseCase {
	uc := NewCase(expression, expected, ShouldEqualJSON, desc...)
	uc.RawOutput = true
	return uc
}

func SkipCaseObject(expression string, expected map[string]any, desc ...string) UseCase {
	return Skip(CaseObject(expression, expected, desc...))
}

func DebugCaseObject(expression string, expected map[string]any, desc ...string) UseCase {
	return Debug(CaseObject(expression, expected, desc...))
}

func CaseArray(expression string, expected []any, desc ...string) UseCase {
	uc := NewCase(expression, expected, ShouldEqualJSON, desc...)
	uc.RawOutput = true
	return uc
}

func SkipCaseArray(expression string, expected []any, desc ...string) UseCase {
	return Skip(CaseArray(expression, expected, desc...))
}

func DebugCaseArray(expression string, expected []any, desc ...string) UseCase {
	return Debug(CaseArray(expression, expected, desc...))
}

func CaseItems(expression string, expected ...any) UseCase {
	return NewCase(expression, expected, base.ShouldHaveSameItems)
}

func CaseFn(expression string, assertion func(actual any, expected ...any) string) UseCase {
	return NewCase(expression, nil, assertion)
}

func SkipCaseFn(expression string, assertion func(actual any, expected ...any) string) UseCase {
	return Skip(CaseFn(expression, assertion))
}

func DebugCaseFn(expression string, assertion func(actual any, expected ...any) string) UseCase {
	return Debug(CaseFn(expression, assertion))
}

func SkipCaseItems(expression string, expected ...any) UseCase {
	return Skip(CaseItems(expression, expected...))
}

func DebugCaseItems(expression string, expected ...any) UseCase {
	return Debug(CaseItems(expression, expected...))
}

func CaseJSON(expression string, expected string, desc ...string) UseCase {
	uc := NewCase(expression, expected, ShouldEqualJSON, desc...)
	uc.RawOutput = true
	return uc
}

func SkipCaseJSON(expression string, expected string, desc ...string) UseCase {
	return Skip(CaseJSON(expression, expected, desc...))
}

func DebugCaseJSON(expression string, expected string, desc ...string) UseCase {
	return Debug(CaseJSON(expression, expected, desc...))
}

func printDebugInfo(name string, uc UseCase, prog *vm.Program) {
	fmt.Println("")
	fmt.Println("VM Test:", name)
	fmt.Println("Expression:", uc.Expression)
	fmt.Println("")
	fmt.Println("Bytecode:")

	out, e := asm.Disassemble(prog, asm.WithDebug())

	if e == nil {
		fmt.Println(out)
	}
}

func RunUseCasesWith(t *testing.T, c *compiler.Compiler, useCases []UseCase, opts ...vm.EnvironmentOption) {
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

		t.Run("VM Test: "+name, func(t *testing.T) {
			if skip {
				t.Skip()

				return
			}

			Convey(useCase.Expression, t, func() {
				prog, err := c.Compile(useCase.Expression)

				defer func() {
					if recovered := recover(); recovered != nil {
						printDebugInfo(name, useCase, prog)

						t.Error(recovered)
					}
				}()

				if !base.ArePtrsEqual(useCase.PreAssertion, base.ShouldBeCompilationError) {
					So(err, ShouldBeNil)
				} else {
					So(err, ShouldBeError)

					return
				}

				options := []vm.EnvironmentOption{
					vm.WithFunctions(std),
				}
				options = append(options, opts...)

				expected := useCase.Expected
				actual, err := base.Exec(prog, useCase.RawOutput, options...)

				for _, assertion := range useCase.Assertions {
					if base.ArePtrsEqual(assertion, ShouldBeError) {
						So(err, ShouldBeError)

						if expected != nil {
							So(err, ShouldBeError, expected)
						} else {
							So(err, ShouldBeError)
						}

						return
					}

					So(err, ShouldBeNil)

					if base.ArePtrsEqual(assertion, ShouldEqualJSON) {
						expectedJ, err := j.Marshal(expected)
						So(err, ShouldBeNil)
						So(actual, ShouldEqualJSON, string(expectedJ))
					} else if base.ArePtrsEqual(assertion, base.ShouldHaveSameItems) {
						So(actual, base.ShouldHaveSameItems, expected)
					} else if base.ArePtrsEqual(assertion, ShouldBeNil) {
						So(actual, ShouldBeNil)
					} else {
						So(actual, assertion, expected)
					}
				}

				// If no assertions are provided, we check the expected value directly
				if len(useCase.Assertions) == 0 {
					So(actual, ShouldEqual, expected)
				}

				if useCase.DebugOutput {
					printDebugInfo(name, useCase, prog)
				}
			})
		})
	}
}

func RunUseCases(t *testing.T, useCases []UseCase, opts ...vm.EnvironmentOption) {
	RunUseCasesWith(t, compiler.New(), useCases, opts...)
}
