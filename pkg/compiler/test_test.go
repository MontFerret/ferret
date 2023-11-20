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

func Exec(p *runtime.Program, raw bool, opts ...runtime.EnvironmentOption) (any, error) {
	vm := runtime.NewVM(opts...)

	out, err := vm.Run(context.Background(), p)

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

func RunUseCases(t *testing.T, c *compiler.Compiler, useCases []UseCase, opts ...runtime.EnvironmentOption) {
	for _, useCase := range useCases {
		Convey(useCase.Expression, t, func() {
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
			} else {
				So(out, ShouldEqual, useCase.Expected)
			}
		})
	}
}
