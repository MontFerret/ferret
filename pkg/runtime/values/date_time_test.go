package values_test

import (
	"encoding/json"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
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

	Convey(".MarshalJSON", t, func() {
		Convey("It should correctly serialize value", func() {
			value := time.Now()

			json1, err := json.Marshal(value)
			So(err, ShouldBeNil)

			json2, err := values.NewDateTime(value).MarshalJSON()
			So(err, ShouldBeNil)

			So(json1, ShouldResemble, json2)
		})
	})
}
