package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/compiler"
)

func TestUnaryOperator(t *testing.T) {
	RunUseCases(t, compiler.New(), []UseCase{
		{"RETURN !TRUE", false, nil},
		{"RETURN !FALSE", true, nil},
		{"RETURN -1", -1, nil},
		{"RETURN -1.1", -1.1, nil},
		{"RETURN +1", 1, nil},
		{"RETURN +1.1", 1.1, nil},
		{`LET v = 1 RETURN -v`, -1, nil},
		{`LET v = 1.1 RETURN -v`, -1.1, nil},
		{`LET v = -1 RETURN -v`, 1, nil},
		{`LET v = -1.1 RETURN -v`, 1.1, nil},
		{`LET v = -1 RETURN +v`, -1, nil},
		{`LET v = -1.1 RETURN +v`, -1.1, nil},
	})

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
