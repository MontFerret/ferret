package compiler_test

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestFunctionCall(t *testing.T) {
	RunUseCases(t, []UseCase{
		{
			"RETURN TYPENAME(1)",
			"int",
			nil,
		},
		{
			"WAIT(10) RETURN 1",
			1,
			nil,
		},
		{
			"LET duration = 10 WAIT(duration) RETURN 1",
			1,
			nil,
		},
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
		{
			`FOR i IN [1, 2, 3, 4]
				LET duration = 10
		
				WAIT(duration)
		
				RETURN i * 2`,
			[]int{2, 4, 6, 8},
			ShouldEqualJSON,
		},
		//{
		//	`RETURN FIRST((FOR i IN 1..10 RETURN i * 2))`,
		//	2,
		//	nil,
		//},
	})

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
