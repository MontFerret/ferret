package eval

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestWrapExp(t *testing.T) {
	Convey("wrapExp", t, func() {
		Convey("When a plain expression is passed", func() {
			exp := "return true"
			So(wrapExp(exp, 0), ShouldEqual, "() => {\n"+exp+"\n}")
		})

		Convey("When a plain expression is passed with args > 0", func() {
			exp := "return true"
			So(wrapExp(exp, 3), ShouldEqual, "(arg1,arg2,arg3) => {\n"+exp+"\n}")
		})
	})
}
