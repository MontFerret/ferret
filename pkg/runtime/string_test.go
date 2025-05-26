package runtime_test

import (
	c "context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"
)

func TestString(t *testing.T) {
	Convey(".Hash", t, func() {
		Convey("It should calculate hash", func() {
			v := runtime.NewString("a")

			h := v.Hash()

			So(h, ShouldBeGreaterThan, 0)

			v2 := runtime.NewString("b")

			So(h, ShouldNotEqual, v2.Hash())
		})

		Convey("Hash sum should be consistent", func() {
			v := runtime.NewString("foobar")

			So(v.Hash(), ShouldEqual, v.Hash())
		})
	})

	Convey(".Length", t, func() {
		Convey("Should return unicode length", func() {
			str := runtime.NewString("Спутник")
			size, err := str.Length(c.Background())

			So(err, ShouldBeNil)
			So(size, ShouldEqual, 7)
		})
	})

	Convey(".MarshalJSON", t, func() {
		Convey("It should correctly serialize Value", func() {
			value := "foobar"

			json1, err := json.Marshal(value)
			So(err, ShouldBeNil)

			json2, err := runtime.NewString(value).MarshalJSON()
			So(err, ShouldBeNil)

			So(json1, ShouldResemble, json2)
		})

		Convey("It should NOT escape HTML", func() {
			value := "<div><span>Foobar</span></div>"

			json1, err := json.Marshal(value)
			So(err, ShouldBeNil)

			json2, err := runtime.NewString(value).MarshalJSON()
			So(err, ShouldBeNil)

			So(json1, ShouldNotResemble, json2)
			So(string(json2), ShouldEqual, fmt.Sprintf(`"%s"`, value))
		})
	})
	Convey(".At", t, func() {
		Convey("It should return a character", func() {
			v := runtime.NewString("abcdefg")
			c := v.At(2)

			So(string(c), ShouldEqual, "c")
		})
	})
}
