package runtime_test

import (
	"encoding/json"
	"github.com/MontFerret/ferret/pkg/runtime"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInt(t *testing.T) {
	Convey(".Hash", t, func() {
		Convey("It should calculate hash", func() {
			v := runtime.NewInt(1)

			h := v.Hash()

			So(h, ShouldBeGreaterThan, 0)

			v2 := runtime.NewInt(2)

			So(h, ShouldNotEqual, v2.Hash())
		})

		Convey("Hash sum should be consistent", func() {
			v := runtime.NewInt(1)

			So(v.Hash(), ShouldEqual, v.Hash())
		})
	})

	Convey(".MarshalJSON", t, func() {
		Convey("It should correctly serialize Value", func() {
			value := 10

			json1, err := json.Marshal(value)
			So(err, ShouldBeNil)

			json2, err := runtime.NewInt(value).MarshalJSON()
			So(err, ShouldBeNil)

			So(json1, ShouldResemble, json2)
		})
	})
}
