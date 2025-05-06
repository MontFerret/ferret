package runtime_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime"

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

	Convey(".MarshalJSON", t, func() {
		Convey("It should correctly serialize Value", func() {
			value := time.Now()

			json1, err := json.Marshal(value)
			So(err, ShouldBeNil)

			json2, err := runtime.NewDateTime(value).MarshalJSON()
			So(err, ShouldBeNil)

			So(json1, ShouldResemble, json2)
		})
	})
}
