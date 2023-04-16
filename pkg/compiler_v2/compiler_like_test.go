package compiler_v2_test

import (
	"context"
	"testing"

	j "encoding/json"

	. "github.com/smartystreets/goconvey/convey"

	compiler "github.com/MontFerret/ferret/pkg/compiler_v2"
	runtime "github.com/MontFerret/ferret/pkg/runtime_v2"
)

func TestLikeOperator(t *testing.T) {
	run := func(p *runtime.Program) (interface{}, error) {
		vm := runtime.NewVM()

		out, err := vm.Run(context.Background(), p)

		if err != nil {
			return false, err
		}

		var i interface{}

		err = j.Unmarshal(out, &i)

		if err != nil {
			return nil, err
		}

		return i, nil
	}

	type UseCase struct {
		Expression string
		Expected   interface{}
	}

	useCases := []UseCase{
		{`RETURN "foo" LIKE "f*"`, true},
		{`RETURN "foo" LIKE "b*"`, false},
		{`RETURN "foo" NOT LIKE "f*"`, false},
		{`RETURN "foo" NOT LIKE "b*"`, true},
		{`LET res = "foo" LIKE  "f*"
			RETURN res`, true},
		{`RETURN ("foo" LIKE  "b*") ? "foo" : "bar"`, `bar`},
		{`RETURN ("foo" NOT LIKE  "b*") ? "foo" : "bar"`, `foo`},
		{`RETURN true ? ("foo" NOT LIKE  "b*") : false`, true},
		{`RETURN true ? false : ("foo" NOT LIKE  "b*")`, false},
		{`RETURN false ? false : ("foo" NOT LIKE  "b*")`, true},
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

	//Convey("FOR IN LIKE", t, func() {
	//	c := compiler.New()
	//
	//	out1, err := c.MustCompile(`
	//		FOR str IN ["foo", "bar", "qaz"]
	//			FILTER str LIKE "*a*"
	//			RETURN str
	//	`).Run(context.Background())
	//
	//	So(err, ShouldBeNil)
	//	So(string(out1), ShouldEqual, `["bar","qaz"]`)
	//})
	//
	//Convey("FOR IN LIKE 2", t, func() {
	//	c := compiler.New()
	//
	//	out1, err := c.MustCompile(`
	//		FOR str IN ["foo", "bar", "qaz"]
	//			FILTER str LIKE "*a*"
	//			RETURN str
	//	`).Run(context.Background())
	//
	//	So(err, ShouldBeNil)
	//	So(string(out1), ShouldEqual, `["bar","qaz"]`)
	//})

}
