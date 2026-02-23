package runtime_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDateTime(t *testing.T) {
	Convey(".Hash", t, func() {
		Convey("It should calculate hash", func() {
			d := runtime.NewCurrentDateTime()

			h := d.Hash()

			So(h, ShouldBeGreaterThan, 0)
		})

		Convey("Hash sum should be consistent", func() {
			d := runtime.NewCurrentDateTime()

			So(d.Hash(), ShouldEqual, d.Hash())
		})
	})
}
