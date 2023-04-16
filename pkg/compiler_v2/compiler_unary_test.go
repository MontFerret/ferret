package compiler_v2_test

import (
	"context"
	j "encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	compiler "github.com/MontFerret/ferret/pkg/compiler_v2"
	runtime "github.com/MontFerret/ferret/pkg/runtime_v2"
)

func TestUnaryOperator(t *testing.T) {
	run := func(p *runtime.Program) (any, error) {
		vm := runtime.NewVM()

		out, err := vm.Run(context.Background(), p)

		if err != nil {
			return 0, err
		}

		var res any

		err = j.Unmarshal(out, &res)

		if err != nil {
			return nil, err
		}

		return res, err
	}

	type UseCase struct {
		Expression string
		Expected   any
	}

	useCases := []UseCase{
		{"RETURN !TRUE", false},
		{"RETURN !FALSE", true},
		{"RETURN -1", -1},
		{"RETURN -1.1", -1.1},
		{"RETURN +1", 1},
		{"RETURN +1.1", 1.1},
		{`			LET v = 1
			RETURN -v`, -1},
		{`			LET v = 1.1
			RETURN -v`, -1.1},
		{`			LET v = -1
			RETURN -v`, 1},
		{`			LET v = -1.1
			RETURN -v`, 1.1},
		{`			LET v = -1
			RETURN +v`, -1},
		{`			LET v = -1.1
			RETURN +v`, -1.1},
	}

	for _, useCase := range useCases {
		Convey("Should compile "+useCase.Expression, t, func() {
			c := compiler.New()

			p, err := c.Compile(useCase.Expression)

			So(err, ShouldBeNil)
			So(p, ShouldHaveSameTypeAs, &runtime.Program{})

			out, err := run(p)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, useCase.Expected)
		})
	}

	//Convey("RETURN { enabled: !val}", t, func() {
	//	c := compiler.New()
	//
	//	out1, err := c.MustCompile(`
	//		LET val = ""
	//		RETURN { enabled: !val }
	//	`).Run(context.Background())
	//
	//	So(err, ShouldBeNil)
	//	So(string(out1), ShouldEqual, `{"enabled":true}`)
	//
	//	out2, err := c.MustCompile(`
	//		LET val = ""
	//		RETURN { enabled: !!val }
	//	`).Run(context.Background())
	//
	//	So(err, ShouldBeNil)
	//	So(string(out2), ShouldEqual, `{"enabled":false}`)
	//})
	//
}
