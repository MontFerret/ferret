package runtime_test

import (
	"testing"

	. "github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBoolean(t *testing.T) {
	Convey(".MarshalJSON", t, func() {
		Convey("Should serialize a boolean items", func() {
			b := True
			marshaled, err := b.MarshalJSON()

			So(err, ShouldBeNil)

			So(string(marshaled), ShouldEqual, "true")
		})
	})

	Convey(".Unwrap", t, func() {
		Convey("Should return an unwrapped items", func() {
			So(True.Unwrap(), ShouldHaveSameTypeAs, true)
		})
	})

	Convey(".String", t, func() {
		Convey("Should return a string representation ", func() {
			So(True.String(), ShouldEqual, "true")
		})
	})

	Convey(".CompareValues", t, func() {
		Convey("It should return 1 when compared to None", func() {
			So(True.Compare(None), ShouldEqual, 1)
		})

		Convey("It should return -1 for all non-boolean values", func() {
			vals := []Value{
				NewString("foo"),
				NewInt(1),
				NewFloat(1.1),
				NewArray(10),
				NewObject(),
			}

			for _, v := range vals {
				So(True.Compare(v), ShouldEqual, -1)
				So(False.Compare(v), ShouldEqual, -1)
			}
		})

		Convey("It should return 0 when both are True or False", func() {
			So(True.Compare(True), ShouldEqual, 0)
			So(False.Compare(False), ShouldEqual, 0)
		})

		Convey("It should return 1 when other is false", func() {
			So(True.Compare(False), ShouldEqual, 1)
		})

		Convey("It should return -1 when other are true", func() {
			So(False.Compare(True), ShouldEqual, -1)
		})
	})

	Convey(".Hash", t, func() {
		Convey("It should calculate hash", func() {
			So(True.Hash(), ShouldBeGreaterThan, 0)
			So(False.Hash(), ShouldBeGreaterThan, 0)
		})

		Convey("Hash sum should be consistent", func() {
			So(True.Hash(), ShouldEqual, True.Hash())
			So(False.Hash(), ShouldEqual, False.Hash())
		})
	})
}
