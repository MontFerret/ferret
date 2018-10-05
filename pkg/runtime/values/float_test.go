package values_test

import (
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestFloat(t *testing.T) {
	Convey(".Hash", t, func() {
		Convey("It should calculate hash", func() {
			v := values.NewFloat(1.1)

			h := v.Hash()

			So(h, ShouldBeGreaterThan, 0)

			v2 := values.NewFloat(1.2)

			So(h, ShouldNotEqual, v2.Hash())
		})

		Convey("Hash sum should be consistent", func() {
			v := values.NewFloat(1.1)

			So(v.Hash(), ShouldEqual, v.Hash())
		})
	})
}
