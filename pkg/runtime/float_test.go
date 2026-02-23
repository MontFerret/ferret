package runtime_test

import (
	"encoding/json"
	"testing"

	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFloat(t *testing.T) {
	Convey(".Hash", t, func() {
		Convey("It should calculate hash", func() {
			v := runtime.NewFloat(1.1)

			h := v.Hash()

			So(h, ShouldBeGreaterThan, 0)

			v2 := runtime.NewFloat(1.2)

			So(h, ShouldNotEqual, v2.Hash())
		})

		Convey("Hash sum should be consistent", func() {
			v := runtime.NewFloat(1.1)

			So(v.Hash(), ShouldEqual, v.Hash())
		})
	})

	Convey(".MarshalJSON", t, func() {
		Convey("It should correctly serialize Value", func() {
			value := float64(10)

			json1, err := json.Marshal(value)
			So(err, ShouldBeNil)

			json2, err := encodingjson.Default.Encode(runtime.NewFloat(value))
			So(err, ShouldBeNil)

			So(json1, ShouldResemble, json2)
		})
	})
}
