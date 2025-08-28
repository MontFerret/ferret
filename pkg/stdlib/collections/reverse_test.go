package collections_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/stdlib/collections"
)

func TestReverse(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = collections.Reverse(context.Background(), runtime.None)

			So(err, ShouldBeError)
		})
	})

	Convey("When argument is invalid type", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = collections.Reverse(context.Background(), runtime.NewInt(123))

			So(err, ShouldBeError)
		})
	})

	Convey("Should reverse a text with right encoding", t, func() {
		out, _ := collections.Reverse(
			context.Background(),
			runtime.NewString("The quick brown 狐 jumped over the lazy 犬"),
		)

		So(string(out.(runtime.String)), ShouldEqual, "犬 yzal eht revo depmuj 狐 nworb kciuq ehT")
	})

	Convey("Should reverse an empty string", t, func() {
		out, err := collections.Reverse(
			context.Background(),
			runtime.NewString(""),
		)

		So(err, ShouldBeNil)
		So(string(out.(runtime.String)), ShouldEqual, "")
	})

	Convey("Should reverse a single character string", t, func() {
		out, err := collections.Reverse(
			context.Background(),
			runtime.NewString("a"),
		)

		So(err, ShouldBeNil)
		So(string(out.(runtime.String)), ShouldEqual, "a")
	})

	Convey("Should handle string with special characters", t, func() {
		out, err := collections.Reverse(
			context.Background(),
			runtime.NewString("a\nb\tc"),
		)

		So(err, ShouldBeNil)
		So(string(out.(runtime.String)), ShouldEqual, "c\tb\na")
	})

	Convey("Should return a copy of an array with reversed elements", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
			runtime.NewInt(6),
		)

		out, err := collections.Reverse(
			context.Background(),
			arr,
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[6,5,4,3,2,1]")
	})

	Convey("Should return an empty array when there no elements in a source one", t, func() {
		arr := runtime.NewArray(0)

		out, err := collections.Reverse(
			context.Background(),
			arr,
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[]")
	})

	Convey("Should reverse array with single element", t, func() {
		arr := runtime.NewArrayWith(runtime.NewInt(42))

		out, err := collections.Reverse(
			context.Background(),
			arr,
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[42]")
	})

	Convey("Should reverse array with different types", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewString("hello"),
			runtime.NewBoolean(true),
		)

		out, err := collections.Reverse(
			context.Background(),
			arr,
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[true,\"hello\",1]")
	})

	Convey("Should not modify original array", t, func() {
		original := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
		)
		originalStr := original.String()

		out, err := collections.Reverse(context.Background(), original)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[3,2,1]")
		So(original.String(), ShouldEqual, originalStr) // Original unchanged
	})
}
