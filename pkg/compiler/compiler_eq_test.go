package compiler_test

import (
	"context"
	runtime2 "github.com/MontFerret/ferret/pkg/runtime"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/compiler"
)

func TestEqualityOperators(t *testing.T) {
	Convey("Equality operators", t, func() {
		run := func(p *runtime2.Program) (string, error) {
			vm := runtime2.NewVM()

			out, err := vm.Run(context.Background(), p)

			return string(out), err
		}

		type UseCase struct {
			Operator string
			Expected bool
		}

		useCases := []UseCase{
			{">", true},
			{"==", false},
			{">=", true},
			{"<", false},
			{"!=", true},
			{"<=", false},
		}

		for _, useCase := range useCases {
			Convey("Should compile RETURN 2 "+useCase.Operator+" 1", func() {
				c := compiler.New()

				p, err := c.Compile(`
				RETURN 2 ` + useCase.Operator + ` 1
			`)

				So(err, ShouldBeNil)
				So(p, ShouldHaveSameTypeAs, &runtime2.Program{})

				out, err := run(p)

				So(err, ShouldBeNil)
				So(out == "true", ShouldEqual, useCase.Expected)
			})
		}
	})
}
