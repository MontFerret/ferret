package runtime_test

import (
	c "context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"

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

	Convey(".At", t, func() {
		Convey("It should return a character", func() {
			v := runtime.NewString("abcdefg")
			c := v.At(2)

			So(string(c), ShouldEqual, "c")
		})
	})
}
