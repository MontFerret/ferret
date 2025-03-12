package core_test

import (
	"encoding/json"
	"fmt"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestString(t *testing.T) {
	Convey(".Hash", t, func() {
		Convey("It should calculate hash", func() {
			v := core.NewString("a")

			h := v.Hash()

			So(h, ShouldBeGreaterThan, 0)

			v2 := core.NewString("b")

			So(h, ShouldNotEqual, v2.Hash())
		})

		Convey("Hash sum should be consistent", func() {
			v := core.NewString("foobar")

			So(v.Hash(), ShouldEqual, v.Hash())
		})
	})

	Convey(".Length", t, func() {
		Convey("Should return unicode length", func() {
			str := core.NewString("Спутник")

			So(str.Length(), ShouldEqual, 7)
		})
	})

	Convey(".MarshalJSON", t, func() {
		Convey("It should correctly serialize Value", func() {
			value := "foobar"

			json1, err := json.Marshal(value)
			So(err, ShouldBeNil)

			json2, err := core.NewString(value).MarshalJSON()
			So(err, ShouldBeNil)

			So(json1, ShouldResemble, json2)
		})

		Convey("It should NOT escape HTML", func() {
			value := "<div><span>Foobar</span></div>"

			json1, err := json.Marshal(value)
			So(err, ShouldBeNil)

			json2, err := core.NewString(value).MarshalJSON()
			So(err, ShouldBeNil)

			So(json1, ShouldNotResemble, json2)
			So(string(json2), ShouldEqual, fmt.Sprintf(`"%s"`, value))
		})
	})
	Convey(".At", t, func() {
		Convey("It should return a character", func() {
			v := core.NewString("abcdefg")
			c := v.At(2)

			So(string(c), ShouldEqual, "c")
		})
	})
}
