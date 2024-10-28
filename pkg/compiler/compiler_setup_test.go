package compiler_test

import (
	"context"
	j "encoding/json"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime"
)

type UseCase struct {
	Expression string
	Expected   any
	Assertion  Assertion
}

func Run(p *runtime.Program, opts ...runtime.EnvironmentOption) ([]byte, error) {
	vm := runtime.NewVM(p)

	out, err := vm.Run(context.Background(), opts...)

	if err != nil {
		return nil, err
	}

	return out.MarshalJSON()
}

func Exec(p *runtime.Program, raw bool, opts ...runtime.EnvironmentOption) (any, error) {
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
	expectedArr := expected[0].([]any)

	for _, item := range expectedArr {
		if err := ShouldContain(actual, item); err != "" {
			return err
		}
	}

	return ""
}

func RunUseCasesWith(t *testing.T, c *compiler.Compiler, useCases []UseCase, opts ...runtime.EnvironmentOption) {
	for _, useCase := range useCases {
		Convey(useCase.Expression, t, func() {
			// catch panic
			//defer func() {
			//	if r := recover(); r != nil {
			//		panic(fmt.Sprintf("%v,\nUse Case %d: - %s", r, idx+1, useCase.Expression))
			//	}
			//}()

			prog, err := c.Compile(useCase.Expression)

			So(err, ShouldBeNil)

			options := []runtime.EnvironmentOption{
				runtime.WithFunctions(c.Functions().Unwrap()),
			}
			options = append(options, opts...)

			out, err := Exec(prog, ArePtrsEqual(useCase.Assertion, ShouldEqualJSON), options...)

			if !ArePtrsEqual(useCase.Assertion, ShouldBeError) {
				So(err, ShouldBeNil)
			}

			if ArePtrsEqual(useCase.Assertion, ShouldEqualJSON) {
				expected, err := j.Marshal(useCase.Expected)
				So(err, ShouldBeNil)
				So(out, ShouldEqualJSON, string(expected))
			} else if ArePtrsEqual(useCase.Assertion, ShouldBeError) {
				if useCase.Expected != nil {
					So(err, ShouldBeError, useCase.Expected)
				} else {
					So(err, ShouldBeError)
				}
			} else if ArePtrsEqual(useCase.Assertion, ShouldHaveSameItems) {
				So(out, ShouldHaveSameItems, useCase.Expected)
			} else {
				So(out, ShouldEqual, useCase.Expected)
			}
		})
	}
}

func RunUseCases(t *testing.T, useCases []UseCase, opts ...runtime.EnvironmentOption) {
	RunUseCasesWith(t, compiler.New(), useCases, opts...)
}

func RunBenchmarkWith(b *testing.B, c *compiler.Compiler, expression string, opts ...runtime.EnvironmentOption) {
	prog, err := c.Compile(expression)

	if err != nil {
		panic(err)
	}

	options := []runtime.EnvironmentOption{
		runtime.WithFunctions(c.Functions().Unwrap()),
	}
	options = append(options, opts...)

	for n := 0; n < b.N; n++ {
		vm := runtime.NewVM(prog)

		_, err := vm.Run(context.Background(), opts...)

		if err != nil {
			panic(err)
		}
	}
}

func RunBenchmark(b *testing.B, expression string, opts ...runtime.EnvironmentOption) {
	RunBenchmarkWith(b, compiler.New(), expression, opts...)
}
