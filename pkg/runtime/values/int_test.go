package values_test

import (
	"encoding/json"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestInt(t *testing.T) {
	Convey(".Hash", t, func() {
		Convey("It should calculate hash", func() {
			v := values.NewInt(1)

			h := v.Hash()

			So(h, ShouldBeGreaterThan, 0)

			v2 := values.NewInt(2)

			So(h, ShouldNotEqual, v2.Hash())
		})

		Convey("Hash sum should be consistent", func() {
			v := values.NewInt(1)

			So(v.Hash(), ShouldEqual, v.Hash())
		})
	})

	Convey(".MarshalJSON", t, func() {
		Convey("It should correctly serialize value", func() {
			value := 10

			json1, err := json.Marshal(value)
			So(err, ShouldBeNil)

			json2, err := values.NewInt(value).MarshalJSON()
			So(err, ShouldBeNil)

			So(json1, ShouldResemble, json2)
		})
	})
}
