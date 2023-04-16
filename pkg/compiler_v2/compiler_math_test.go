package compiler_v2_test

import (
	"context"
	j "encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	compiler "github.com/MontFerret/ferret/pkg/compiler_v2"
	runtime "github.com/MontFerret/ferret/pkg/runtime_v2"
)

func TestMathOperators(t *testing.T) {
	Convey("Math operators", t, func() {
		run := func(p *runtime.Program) (int, error) {
			vm := runtime.NewVM()

			out, err := vm.Run(context.Background(), p)

			if err != nil {
				return 0, err
			}

			var res int

			err = j.Unmarshal(out, &res)

			if err != nil {
				return 0, err
			}

			return res, err
		}

		type UseCase struct {
			Operator string
			Expected int
		}

		useCases := []UseCase{
			{"+", 6},
			{"-", 2},
			{"*", 8},
			{"/", 2},
			{"%", 0},
		}

		for _, useCase := range useCases {
			exp := "RETURN 4 " + useCase.Operator + " 2"

			Convey("Should compile "+exp, func() {
				c := compiler.New()

				p, err := c.Compile(exp)

				So(err, ShouldBeNil)
				So(p, ShouldHaveSameTypeAs, &runtime.Program{})

				out, err := run(p)

				So(err, ShouldBeNil)
				So(out, ShouldEqual, useCase.Expected)
			})
		}
	})
}
