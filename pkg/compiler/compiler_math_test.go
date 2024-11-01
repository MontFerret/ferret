package compiler_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"testing"

	runtime2 "github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/compiler"
)

func TestMathOperators(t *testing.T) {
	Convey("Math operators", t, func() {
		run := func(p *runtime2.Program) (int, error) {
			vm := runtime2.NewVM()

			out, err := vm.Run(context.Background(), p)

			if err != nil {
				return 0, err
			}

			return int(values.ToInt(out)), nil
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
				So(p, ShouldHaveSameTypeAs, &runtime2.Program{})

				out, err := run(p)

				So(err, ShouldBeNil)
				So(out, ShouldEqual, useCase.Expected)
			})
		}
	})
}
