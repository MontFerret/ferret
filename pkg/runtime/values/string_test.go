package values_test

import (
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestString(t *testing.T) {
	Convey(".Hash", t, func() {
		Convey("It should calculate hash", func() {
			v := values.NewString("a")

			h := v.Hash()

			So(h, ShouldBeGreaterThan, 0)

			v2 := values.NewString("b")

			So(h, ShouldNotEqual, v2.Hash())
		})

		Convey("Hash sum should be consistent", func() {
			v := values.NewString("foobar")

			So(v.Hash(), ShouldEqual, v.Hash())
		})
	})
}
