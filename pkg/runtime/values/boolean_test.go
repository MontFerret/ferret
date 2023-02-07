package values_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

func TestBoolean(t *testing.T) {
	Convey(".MarshalJSON", t, func() {
		Convey("Should serialize a boolean items", func() {
			b := values.True
			marshaled, err := b.MarshalJSON()

			So(err, ShouldBeNil)

			So(string(marshaled), ShouldEqual, "true")
		})
	})

	Convey(".Type", t, func() {
		Convey("Should return a type", func() {
			So(values.True.Type().Equals(types.Boolean), ShouldBeTrue)
		})
	})

	Convey(".Unwrap", t, func() {
		Convey("Should return an unwrapped items", func() {
			So(values.True.Unwrap(), ShouldHaveSameTypeAs, true)
		})
	})

	Convey(".String", t, func() {
		Convey("Should return a string representation ", func() {
			So(values.True.String(), ShouldEqual, "true")
		})
	})

	Convey(".Compare", t, func() {
		Convey("It should return 1 when compared to None", func() {
			So(values.True.Compare(values.None), ShouldEqual, 1)
		})

		Convey("It should return -1 for all non-boolean values", func() {
			vals := []core.Value{
				values.NewString("foo"),
				values.NewInt(1),
				values.NewFloat(1.1),
				values.NewArray(10),
				values.NewObject(),
			}

			for _, v := range vals {
				So(values.True.Compare(v), ShouldEqual, -1)
				So(values.False.Compare(v), ShouldEqual, -1)
			}
		})

		Convey("It should return 0 when both are True or False", func() {
			So(values.True.Compare(values.True), ShouldEqual, 0)
			So(values.False.Compare(values.False), ShouldEqual, 0)
		})

		Convey("It should return 1 when other is false", func() {
			So(values.True.Compare(values.False), ShouldEqual, 1)
		})

		Convey("It should return -1 when other are true", func() {
			So(values.False.Compare(values.True), ShouldEqual, -1)
		})
	})

	Convey(".Hash", t, func() {
		Convey("It should calculate hash", func() {
			So(values.True.Hash(), ShouldBeGreaterThan, 0)
			So(values.False.Hash(), ShouldBeGreaterThan, 0)
		})

		Convey("Hash sum should be consistent", func() {
			So(values.True.Hash(), ShouldEqual, values.True.Hash())
			So(values.False.Hash(), ShouldEqual, values.False.Hash())
		})
	})
}
