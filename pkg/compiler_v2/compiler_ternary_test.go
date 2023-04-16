package compiler_v2_test

import (
	"context"
	j "encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	compiler "github.com/MontFerret/ferret/pkg/compiler_v2"
	runtime "github.com/MontFerret/ferret/pkg/runtime_v2"
)

func TestTernaryOperator(t *testing.T) {
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
		{"RETURN 1 < 2 ? 3 : 4", 3},
		{"RETURN 1 > 2 ? 3 : 4", 4},
		{"RETURN 2 ? : 4", 2},
		{`
LET foo = TRUE
RETURN foo ? TRUE : FALSE
`, true},
		{`
LET foo = FALSE
RETURN foo ? TRUE : FALSE
`, false},
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

	//Convey("Should compile ternary operator", t, func() {
	//	c := compiler.New()
	//	p, err := c.Compile(`
	//		FOR i IN [1, 2, 3, 4, 5, 6]
	//			RETURN i < 3 ? i * 3 : i * 2
	//	`)
	//
	//	So(err, ShouldBeNil)
	//
	//	out, err := p.Run(context.Background())
	//
	//	So(err, ShouldBeNil)
	//
	//	So(string(out), ShouldEqual, `[3,6,6,8,10,12]`)
	//})
	//
	//Convey("Should compile ternary operator with shortcut", t, func() {
	//	c := compiler.New()
	//	p, err := c.Compile(`
	//		FOR i IN [1, 2, 3, 4, 5, 6]
	//			RETURN i < 3 ? : i * 2
	//	`)
	//
	//	So(err, ShouldBeNil)
	//
	//	out, err := p.Run(context.Background())
	//
	//	So(err, ShouldBeNil)
	//
	//	So(string(out), ShouldEqual, `[true,true,6,8,10,12]`)
	//})
	//
	//Convey("Should compile ternary operator with shortcut with nones", t, func() {
	//	c := compiler.New()
	//	p, err := c.Compile(`
	//		FOR i IN [NONE, 2, 3, 4, 5, 6]
	//			RETURN i ? : i
	//	`)
	//
	//	So(err, ShouldBeNil)
	//
	//	out, err := p.Run(context.Background())
	//
	//	So(err, ShouldBeNil)
	//
	//	So(string(out), ShouldEqual, `[null,2,3,4,5,6]`)
	//})
	//
	//Convey("Should compile ternary operator with default values", t, func() {
	//	vals := []string{
	//		"0",
	//		"0.0",
	//		"''",
	//		"NONE",
	//		"FALSE",
	//	}
	//
	//	c := compiler.New()
	//
	//	for _, val := range vals {
	//		p, err := c.Compile(fmt.Sprintf(`
	//		FOR i IN [%s, 1, 2, 3]
	//			RETURN i ? i * 2 : 'no value'
	//	`, val))
	//
	//		So(err, ShouldBeNil)
	//
	//		out, err := p.Run(context.Background())
	//
	//		So(err, ShouldBeNil)
	//
	//		So(string(out), ShouldEqual, `["no value",2,4,6]`)
	//	}
	//})
	//
	//Convey("Multi expression", t, func() {
	//	out := compiler.New().MustCompile(`
	//		RETURN 0 && true ? "1" : "some"
	//	`).MustRun(context.Background())
	//
	//	So(string(out), ShouldEqual, `"some"`)
	//
	//	out = compiler.New().MustCompile(`
	//		RETURN length([]) > 0 && true ? "1" : "some"
	//	`).MustRun(context.Background())
	//
	//	So(string(out), ShouldEqual, `"some"`)
	//})
}
