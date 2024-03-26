package compiler_test

import (
	"testing"
)

func TestLikeOperator(t *testing.T) {
	RunUseCases(t, []UseCase{
		{`RETURN "foo" LIKE "f*"`, true, nil},
		{`RETURN "foo" LIKE "b*"`, false, nil},
		{`RETURN "foo" NOT LIKE "f*"`, false, nil},
		{`RETURN "foo" NOT LIKE "b*"`, true, nil},
		{`LET res = "foo" LIKE  "f*"
			RETURN res`, true, nil},
		{`RETURN ("foo" LIKE  "b*") ? "foo" : "bar"`, `bar`, nil},
		{`RETURN ("foo" NOT LIKE  "b*") ? "foo" : "bar"`, `foo`, nil},
		{`RETURN true ? ("foo" NOT LIKE  "b*") : false`, true, nil},
		{`RETURN true ? false : ("foo" NOT LIKE  "b*")`, false, nil},
		{`RETURN false ? false : ("foo" NOT LIKE  "b*")`, true, nil},
	})

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
