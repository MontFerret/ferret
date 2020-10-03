package values_test

import (
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestString(t *testing.T) {
	Convey(".Hash", t, func() {
		Convey("It should calculate hash", func() {
			v := values.NewString("a")

			h := v.Hash()

			So(h, ShouldBeGreaterThan, 0)

			v2 := values.NewString("b")

			So(h, ShouldNotEqual, v2.Hash())
		})

		Convey("Hash sum should be consistent", func() {
			v := values.NewString("foobar")

			So(v.Hash(), ShouldEqual, v.Hash())
		})
	})

	Convey(".Length", t, func() {
		Convey("Should return unicode length", func() {
			str := values.NewString("Спутник")

			So(str.Length(), ShouldEqual, 7)
		})
	})

	Convey(".MarshalJSON", t, func() {
		Convey("Should not HTML escape", func() {
			str := values.NewString("name=Jane&age=38")
			expectedValue := []byte(`"name=Jane&age=38"`)

			marshalled, err := str.MarshalJSON()
			So(err, ShouldBeNil)
			So(string(marshalled), ShouldEqual, string(expectedValue))
		})
	})
}
