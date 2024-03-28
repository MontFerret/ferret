package compiler_test

import (
	"testing"
)

func TestFunctionCall(t *testing.T) {
	RunUseCases(t, []UseCase{
		//{
		//	"RETURN TYPENAME(1)",
		//	"int",
		//	nil,
		//},
		//{
		//	"WAIT(10) RETURN 1",
		//	1,
		//	nil,
		//},
		//{
		//	"LET duration = 10 WAIT(duration) RETURN 1",
		//	1,
		//	nil,
		//},
		{
			"RETURN (FALSE OR T::FAIL())?",
			nil,
			nil,
		},
		{
			"RETURN T::FAIL()?",
			nil,
			nil,
		},
		//{
		//	`FOR i IN [1, 2, 3, 4]
		//		LET duration = 10
		//
		//		WAIT(duration)
		//
		//		RETURN i * 2`,
		//	[]int{2, 4, 6, 8},
		//	ShouldEqualJSON,
		//},
	})

	//
	//Convey("Should handle errors when ? is used", t, func() {
	//	c := compiler.New()
	//	c.RegisterFunction("ERROR", func(ctx context.Context, args ...core.Value) (core.Value, error) {
	//		return values.None, errors.New("test error")
	//	})
	//
	//	p, err := c.Compile(`
	//		RETURN ERROR()?
	//	`)
	//
	//	So(err, ShouldBeNil)
	//
	//	out, err := p.Run(context.Background())
	//
	//	So(err, ShouldBeNil)
	//
	//	So(string(out), ShouldEqual, `null`)
	//})
	//
	//Convey("Should handle errors when ? is used within a group", t, func() {
	//	c := compiler.New()
	//
	//	p, err := c.Compile(`
	//		RETURN (FALSE OR T::FAIL())?
	//	`)
	//
	//	So(err, ShouldBeNil)
	//
	//	out, err := p.Run(context.Background())
	//
	//	So(err, ShouldBeNil)
	//
	//	So(string(out), ShouldEqual, `null`)
	//})
	//
	//Convey("Should return NONE when error is handled", t, func() {
	//	c := compiler.New()
	//	c.RegisterFunction("ERROR", func(ctx context.Context, args ...core.Value) (core.Value, error) {
	//		return values.NewString("booo"), errors.New("test error")
	//	})
	//
	//	p, err := c.Compile(`
	//		RETURN ERROR()?
	//	`)
	//
	//	So(err, ShouldBeNil)
	//
	//	out, err := p.Run(context.Background())
	//
	//	So(err, ShouldBeNil)
	//
	//	So(string(out), ShouldEqual, `null`)
	//})
	//
	//Convey("Should be able to use FOR as an argument", t, func() {
	//	c := compiler.New()
	//
	//	p, err := c.Compile(`
	//		RETURN FIRST((FOR i IN 1..10 RETURN i * 2))
	//	`)
	//
	//	So(err, ShouldBeNil)
	//
	//	out, err := p.Run(context.Background())
	//
	//	So(err, ShouldBeNil)
	//
	//	So(string(out), ShouldEqual, `2`)
	//})
	//
	//Convey("Should be able to use FOR as arguments", t, func() {
	//	c := compiler.New()
	//
	//	p, err := c.Compile(`
	//		RETURN UNION((FOR i IN 0..5 RETURN i), (FOR i IN 6..10 RETURN i))
	//	`)
	//
	//	So(err, ShouldBeNil)
	//
	//	out, err := p.Run(context.Background())
	//
	//	So(err, ShouldBeNil)
	//
	//	So(string(out), ShouldEqual, `[0,1,2,3,4,5,6,7,8,9,10]`)
	//})
}
