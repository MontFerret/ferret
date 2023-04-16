package compiler_v2_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	compiler "github.com/MontFerret/ferret/pkg/compiler_v2"
	runtime "github.com/MontFerret/ferret/pkg/runtime_v2"
)

func TestEqualityOperators(t *testing.T) {
	Convey("Equality operators", t, func() {
		run := func(p *runtime.Program) (string, error) {
			vm := runtime.NewVM()

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
				So(p, ShouldHaveSameTypeAs, &runtime.Program{})

				out, err := run(p)

				So(err, ShouldBeNil)
				So(out == "true", ShouldEqual, useCase.Expected)
			})
		}
	})
}
