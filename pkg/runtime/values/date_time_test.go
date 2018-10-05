package values_test

import (
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDateTime(t *testing.T) {
	Convey(".Hash", t, func() {
		Convey("It should calculate hash", func() {
			d := values.NewCurrentDateTime()

			h := d.Hash()

			So(h, ShouldBeGreaterThan, 0)
		})

		Convey("Hash sum should be consistent", func() {
			d := values.NewCurrentDateTime()

			So(d.Hash(), ShouldEqual, d.Hash())
		})
	})
}
