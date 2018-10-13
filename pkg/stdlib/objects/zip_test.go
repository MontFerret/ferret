package objects_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/objects"

	. "github.com/smartystreets/goconvey/convey"
)

func TestZip(t *testing.T) {
	Convey("Invalid arguments", t, func() {
		Convey("When there are no arguments", func() {
			actual, err := objects.Zip(context.Background())

			So(err, ShouldBeError)
			So(actual.Compare(values.None), ShouldEqual, 0)
		})

		Convey("When single argument", func() {
			actual, err := objects.Zip(context.Background(), values.NewArray(0))

			So(err, ShouldBeError)
			So(actual.Compare(values.None), ShouldEqual, 0)

			actual, err = objects.Zip(context.Background(), values.NewInt(0))

			So(err, ShouldBeError)
			So(actual.Compare(values.None), ShouldEqual, 0)
		})

		Convey("When too many arguments", func() {
			actual, err := objects.Zip(context.Background(),
				values.NewArray(0), values.NewArray(0), values.NewArray(0))

			So(err, ShouldBeError)
			So(actual.Compare(values.None), ShouldEqual, 0)
		})

		Convey("When there is not array argument", func() {
			actual, err := objects.Zip(context.Background(), values.NewArray(0), values.NewInt(0))

			So(err, ShouldBeError)
			So(actual.Compare(values.None), ShouldEqual, 0)

			actual, err = objects.Zip(context.Background(), values.NewInt(0), values.NewArray(0))

			So(err, ShouldBeError)
			So(actual.Compare(values.None), ShouldEqual, 0)
		})

		Convey("When there is not string element into keys array", func() {
			keys := values.NewArrayWith(values.NewInt(0))
			vals := values.NewArrayWith(values.NewString("v1"))

			actual, err := objects.Zip(context.Background(), keys, vals)

			So(err, ShouldBeError)
			So(actual.Compare(values.None), ShouldEqual, 0)
		})

		Convey("When 1 key and 0 values", func() {
			keys := values.NewArrayWith(values.NewString("k1"))
			vals := values.NewArray(0)

			actual, err := objects.Zip(context.Background(), keys, vals)

			So(err, ShouldBeError)
			So(actual.Compare(values.None), ShouldEqual, 0)
		})

		Convey("When 0 keys and 1 values", func() {
			keys := values.NewArray(0)
			vals := values.NewArrayWith(values.NewString("v1"))

			actual, err := objects.Zip(context.Background(), keys, vals)

			So(err, ShouldBeError)
			So(actual.Compare(values.None), ShouldEqual, 0)
		})
	})
}
