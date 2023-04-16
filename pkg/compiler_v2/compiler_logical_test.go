package compiler_v2_test

import (
	"context"
	"testing"

	j "encoding/json"

	. "github.com/smartystreets/goconvey/convey"

	compiler "github.com/MontFerret/ferret/pkg/compiler_v2"
	runtime "github.com/MontFerret/ferret/pkg/runtime_v2"
)

func TestLogicalOperators(t *testing.T) {
	run := func(p *runtime.Program) (any, error) {
		vm := runtime.NewVM()

		out, err := vm.Run(context.Background(), p)

		if err != nil {
			return 0, err
		}

		var res interface{}

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
		{"RETURN 1 AND 1", 1},
		{"RETURN 1 AND 0", 0},
		{"RETURN 2 > 1 AND 1 > 0", true},
		{"RETURN NONE && true", nil},
		{"RETURN '' && true", ""},
		{"RETURN true && 23", 23},
		{"RETURN 2 > 1 OR 1 < 0", true},
		{"RETURN 1 || 7", 1},
		{"RETURN 0 || 7", 7},
		{"RETURN NONE || 'foo'", "foo"},
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

	//
	//Convey("ERROR()? || 'boo'  should return 'boo'", t, func() {
	//	c := compiler.New()
	//	c.RegisterFunction("ERROR", func(ctx context.Context, args ...core.Value) (core.Value, error) {
	//		return nil, errors.New("test")
	//	})
	//
	//	p, err := c.Compile(`
	//		RETURN ERROR()? || 'boo'
	//	`)
	//
	//	So(err, ShouldBeNil)
	//
	//	out, err := p.Run(context.Background())
	//
	//	So(err, ShouldBeNil)
	//	So(string(out), ShouldEqual, `"boo"`)
	//})
	//
	//Convey("!ERROR()? && TRUE should return false", t, func() {
	//	c := compiler.New()
	//	c.RegisterFunction("ERROR", func(ctx context.Context, args ...core.Value) (core.Value, error) {
	//		return nil, errors.New("test")
	//	})
	//
	//	p, err := c.Compile(`
	//		RETURN !ERROR()? && TRUE
	//	`)
	//
	//	So(err, ShouldBeNil)
	//
	//	out, err := p.Run(context.Background())
	//
	//	So(err, ShouldBeNil)
	//	So(string(out), ShouldEqual, `true`)
	//})
	//
	//

	//
	//Convey("NOT TRUE should return false", t, func() {
	//	c := compiler.New()
	//
	//	p, err := c.Compile(`
	//		RETURN NOT TRUE
	//	`)
	//
	//	So(err, ShouldBeNil)
	//
	//	out, err := p.Run(context.Background())
	//
	//	So(err, ShouldBeNil)
	//	So(string(out), ShouldEqual, `false`)
	//})
	//
	//Convey("NOT u.valid should return true", t, func() {
	//	c := compiler.New()
	//
	//	p, err := c.Compile(`
	//		LET u = { valid: false }
	//
	//		RETURN NOT u.valid
	//	`)
	//
	//	So(err, ShouldBeNil)
	//
	//	out, err := p.Run(context.Background())
	//
	//	So(err, ShouldBeNil)
	//	So(string(out), ShouldEqual, `true`)
	//})
}
