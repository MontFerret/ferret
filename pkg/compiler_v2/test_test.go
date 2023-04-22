package compiler_v2_test

import (
	"context"
	j "encoding/json"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	runtime "github.com/MontFerret/ferret/pkg/runtime_v2"

	compiler "github.com/MontFerret/ferret/pkg/compiler_v2"
)

type UseCase struct {
	Expression string
	Expected   any
	Assertion  Assertion
}

func Exec(p *runtime.Program, raw bool) (any, error) {
	vm := runtime.NewVM()

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

func RunUseCases(t *testing.T, c *compiler.Compiler, useCases []UseCase) {
	for _, useCase := range useCases {
		Convey(useCase.Expression, t, func() {
			prog, err := c.Compile(useCase.Expression)

			So(err, ShouldBeNil)

			out, err := Exec(prog, ArePtrsEqual(useCase.Assertion, ShouldEqualJSON))

			So(err, ShouldBeNil)

			if useCase.Assertion != nil {
				if ArePtrsEqual(useCase.Assertion, ShouldEqualJSON) {
					expected, err := j.Marshal(useCase.Expected)
					So(err, ShouldBeNil)
					So(out, ShouldEqualJSON, string(expected))
				} else {
					So(out, useCase.Assertion, useCase.Expected)
				}
			} else {
				So(out, ShouldEqual, useCase.Expected)
			}
		})
	}
}
