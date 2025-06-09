package base

import (
	"context"
	j "encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"

	. "github.com/smartystreets/goconvey/convey"
)

type UseCase struct {
	Expression  string
	Expected    any
	Assertion   Assertion
	Description string
	Skip        bool
}

func NewCase(expression string, expected any, assertion Assertion, desc ...string) UseCase {
	return UseCase{
		Expression:  expression,
		Expected:    expected,
		Assertion:   assertion,
		Description: strings.TrimSpace(strings.Join(desc, " ")),
	}
}

func Skip(uc UseCase) UseCase {
	uc.Skip = true
	return uc
}

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

func CaseRuntimeError(expression string, desc ...string) UseCase {
	return NewCase(expression, nil, ShouldBeError, desc...)
}

func CaseRuntimeErrorAs(expression string, expected error, desc ...string) UseCase {
	return NewCase(expression, expected, ShouldBeError, desc...)
}

func SkipCaseRuntimeError(expression string, desc ...string) UseCase {
	return Skip(CaseRuntimeError(expression, desc...))
}

func SkipCaseRuntimeErrorAs(expression string, expected error, desc ...string) UseCase {
	return Skip(CaseRuntimeErrorAs(expression, expected, desc...))
}

func CaseCompilationError(expression string, desc ...string) UseCase {
	return NewCase(expression, nil, ShouldBeCompilationError, desc...)
}

func SkipCaseCompilationError(expression string, desc ...string) UseCase {
	return Skip(CaseCompilationError(expression, desc...))
}

func CaseObject(expression string, expected map[string]any, desc ...string) UseCase {
	return NewCase(expression, expected, ShouldEqualJSON, desc...)
}

func SkipCaseObject(expression string, expected map[string]any, desc ...string) UseCase {
	return Skip(CaseObject(expression, expected, desc...))
}

func CaseArray(expression string, expected []any, desc ...string) UseCase {
	return NewCase(expression, expected, ShouldEqualJSON, desc...)
}

func SkipCaseArray(expression string, expected []any, desc ...string) UseCase {
	return Skip(CaseArray(expression, expected, desc...))
}

func CaseItems(expression string, expected ...any) UseCase {
	return NewCase(expression, expected, ShouldHaveSameItems)
}

func SkipCaseItems(expression string, expected ...any) UseCase {
	return Skip(CaseItems(expression, expected...))
}

func CaseJSON(expression string, expected string, desc ...string) UseCase {
	return NewCase(expression, expected, ShouldEqualJSON, desc...)
}

func SkipCaseJSON(expression string, expected string, desc ...string) UseCase {
	return Skip(CaseJSON(expression, expected, desc...))
}

type ExpectedProgram struct {
	Disassembly string
	Constants   []runtime.Value
	Registers   int
}

type ByteCodeUseCase struct {
	Expression string
	Expected   ExpectedProgram
}

func Compile(expression string) (*vm.Program, error) {
	c := compiler.New()

	return c.Compile(expression)
}

func Run(p *vm.Program, opts ...vm.EnvironmentOption) ([]byte, error) {
	instance := vm.New(p)

	out, err := instance.Run(context.Background(), opts)

	if err != nil {
		return nil, err
	}

	return out.MarshalJSON()
}

func Exec(p *vm.Program, raw bool, opts ...vm.EnvironmentOption) (any, error) {
	out, err := Run(p, opts...)

	if err != nil {
		return 0, err
	}

	if raw {
		return string(out), nil
	}

	var res any

	err = j.Unmarshal(out, &res)

	if err != nil {
		return nil, err
	}

	return res, err
}

func ArePtrsEqual(expected, actual any) bool {
	if expected == nil || actual == nil {
		return false
	}

	p1 := fmt.Sprintf("%v", expected)
	p2 := fmt.Sprintf("%v", actual)

	return p1 == p2
}

func ShouldHaveSameItems(actual any, expected ...any) string {
	wapper := expected[0].([]any)
	expectedArr := wapper[0].([]any)

	for _, item := range expectedArr {
		if err := ShouldContain(actual, item); err != "" {
			return err
		}
	}

	return ""
}

func ShouldBeCompilationError(actual any, _ ...any) string {
	// TODO: Expect a particular error message

	So(actual, ShouldBeError)

	return ""
}

func RunAsmUseCases(t *testing.T, useCases []ByteCodeUseCase) {
	c := compiler.New()
	for _, useCase := range useCases {
		t.Run(fmt.Sprintf("Bytecode: %s", useCase.Expression), func(t *testing.T) {
			Convey(useCase.Expression, t, func() {
				assertJSON := func(actual, expected interface{}) {
					actualJ, err := j.Marshal(actual)
					So(err, ShouldBeNil)

					expectedJ, err := j.Marshal(expected)
					So(err, ShouldBeNil)

					So(string(actualJ), ShouldEqualJSON, string(expectedJ))
				}

				prog, err := c.Compile(useCase.Expression)

				So(err, ShouldBeNil)

				So(strings.TrimSpace(prog.Disassemble()), ShouldEqual, strings.TrimSpace(useCase.Expected.Disassembly))

				assertJSON(prog.Constants, useCase.Expected.Constants)
				//assertJSON(prog.CatchTable, useCase.Expected.CatchTable)
				//So(prog.Registers, ShouldEqual, useCase.Expected.Registers)
			})
		})
	}
}

func RunUseCasesWith(t *testing.T, c *compiler.Compiler, useCases []UseCase, opts ...vm.EnvironmentOption) {
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

		t.Run(name, func(t *testing.T) {
			if skip {
				t.Skip()

				return
			}

			Convey(useCase.Expression, t, func() {
				prog, err := c.Compile(useCase.Expression)

				if !ArePtrsEqual(useCase.Assertion, ShouldBeCompilationError) {
					So(err, ShouldBeNil)
				} else {
					So(err, ShouldBeError)

					return
				}

				options := []vm.EnvironmentOption{
					vm.WithFunctions(c.Functions().Unwrap()),
				}
				options = append(options, opts...)

				expected := useCase.Expected
				actual, err := Exec(prog, ArePtrsEqual(useCase.Assertion, ShouldEqualJSON), options...)

				if ArePtrsEqual(useCase.Assertion, ShouldBeError) {
					So(err, ShouldBeError)

					if expected != nil {
						So(err, ShouldBeError, expected)
					} else {
						So(err, ShouldBeError)
					}

					return
				}

				So(err, ShouldBeNil)

				if ArePtrsEqual(useCase.Assertion, ShouldEqualJSON) {
					expectedJ, err := j.Marshal(expected)
					So(err, ShouldBeNil)
					So(actual, ShouldEqualJSON, string(expectedJ))
				} else if ArePtrsEqual(useCase.Assertion, ShouldHaveSameItems) {
					So(actual, ShouldHaveSameItems, expected)
				} else if ArePtrsEqual(useCase.Assertion, ShouldBeNil) {
					So(actual, ShouldBeNil)
				} else if useCase.Assertion == nil {
					So(actual, ShouldEqual, expected)
				} else {
					So(actual, useCase.Assertion, expected)
				}
			})
		})
	}
}

func RunUseCases(t *testing.T, useCases []UseCase, opts ...vm.EnvironmentOption) {
	RunUseCasesWith(t, compiler.New(), useCases, opts...)
}

func RunBenchmarkWith(b *testing.B, c *compiler.Compiler, expression string, opts ...vm.EnvironmentOption) {
	prog, err := c.Compile(expression)

	if err != nil {
		panic(err)
	}

	options := []vm.EnvironmentOption{
		vm.WithFunctions(c.Functions().Unwrap()),
	}
	options = append(options, opts...)

	ctx := context.Background()
	vm := vm.New(prog)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_, err := vm.Run(ctx, opts)

		if err != nil {
			panic(err)
		}
	}
}

func RunBenchmark(b *testing.B, expression string, opts ...vm.EnvironmentOption) {
	RunBenchmarkWith(b, compiler.New(), expression, opts...)
}
