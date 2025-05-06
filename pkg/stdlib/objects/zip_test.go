package objects_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"

	"github.com/MontFerret/ferret/pkg/stdlib/objects"

	. "github.com/smartystreets/goconvey/convey"
)

func TestZip(t *testing.T) {
	Convey("Invalid arguments", t, func() {
		Convey("When there are no arguments", func() {
			actual, err := objects.Zip(context.Background())
			expected := core.None

			So(err, ShouldBeError)
			So(actual.Compare(expected), ShouldEqual, 0)
		})

		Convey("When single argument", func() {
			actual, err := objects.Zip(context.Background(), internal.NewArray(0))
			expected := core.None

			So(err, ShouldBeError)
			So(actual.Compare(expected), ShouldEqual, 0)

			actual, err = objects.Zip(context.Background(), core.NewInt(0))

			So(err, ShouldBeError)
			So(actual.Compare(expected), ShouldEqual, 0)
		})

		Convey("When too many arguments", func() {
			actual, err := objects.Zip(context.Background(),
				internal.NewArray(0), internal.NewArray(0), internal.NewArray(0))
			expected := core.None

			So(err, ShouldBeError)
			So(actual.Compare(expected), ShouldEqual, 0)
		})

		Convey("When there is not array argument", func() {
			actual, err := objects.Zip(context.Background(), internal.NewArray(0), core.NewInt(0))
			expected := core.None

			So(err, ShouldBeError)
			So(actual.Compare(expected), ShouldEqual, 0)

			actual, err = objects.Zip(context.Background(), core.NewInt(0), internal.NewArray(0))

			So(err, ShouldBeError)
			So(actual.Compare(expected), ShouldEqual, 0)
		})

		Convey("When there is not string element into keys array", func() {
			keys := internal.NewArrayWith(core.NewInt(0))
			vals := internal.NewArrayWith(core.NewString("v1"))
			expected := core.None

			actual, err := objects.Zip(context.Background(), keys, vals)

			So(err, ShouldBeError)
			So(actual.Compare(expected), ShouldEqual, 0)
		})

		Convey("When 1 key and 0 values", func() {
			keys := internal.NewArrayWith(core.NewString("k1"))
			vals := internal.NewArray(0)
			expected := core.None

			actual, err := objects.Zip(context.Background(), keys, vals)

			So(err, ShouldBeError)
			So(actual.Compare(expected), ShouldEqual, 0)
		})

		Convey("When 0 keys and 1 values", func() {
			keys := internal.NewArray(0)
			vals := internal.NewArrayWith(core.NewString("v1"))
			expected := core.None

			actual, err := objects.Zip(context.Background(), keys, vals)

			So(err, ShouldBeError)
			So(actual.Compare(expected), ShouldEqual, 0)
		})
	})

	Convey("Zip 2 keys and 2 values", t, func() {
		keys := internal.NewArrayWith(
			core.NewString("k1"),
			core.NewString("k2"),
		)
		vals := internal.NewArrayWith(
			core.NewString("v1"),
			core.NewInt(2),
		)
		expected := internal.NewObjectWith(
			internal.NewObjectProperty("k1", core.NewString("v1")),
			internal.NewObjectProperty("k2", core.NewInt(2)),
		)

		actual, err := objects.Zip(context.Background(), keys, vals)

		So(err, ShouldBeNil)
		So(actual.Compare(expected), ShouldEqual, 0)
	})

	Convey("Zip 3 keys and 3 values. 1 key repeats", t, func() {
		keys := internal.NewArrayWith(
			core.NewString("k1"),
			core.NewString("k2"),
			core.NewString("k1"),
		)
		vals := internal.NewArrayWith(
			core.NewInt(1),
			core.NewInt(2),
			core.NewInt(3),
		)
		expected := internal.NewObjectWith(
			internal.NewObjectProperty("k1", core.NewInt(1)),
			internal.NewObjectProperty("k2", core.NewInt(2)),
		)

		actual, err := objects.Zip(context.Background(), keys, vals)

		So(err, ShouldBeNil)
		So(actual.Compare(expected), ShouldEqual, 0)
	})
}
