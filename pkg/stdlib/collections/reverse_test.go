package collections_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	collections2 "github.com/MontFerret/ferret/pkg/runtime/internal"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/collections"
)

func TestReverse(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = collections.Reverse(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("Should reverse a text with right encoding", t, func() {
		out, _ := collections.Reverse(
			context.Background(),
			core.NewString("The quick brown 狐 jumped over the lazy 犬"),
		)

		So(out, ShouldEqual, "犬 yzal eht revo depmuj 狐 nworb kciuq ehT")
	})

	Convey("Should return a copy of an array with reversed elements", t, func() {
		arr := internal.NewArrayWith(
			core.NewInt(1),
			core.NewInt(2),
			core.NewInt(3),
			core.NewInt(4),
			core.NewInt(5),
			core.NewInt(6),
		)

		out, err := collections.Reverse(
			context.Background(),
			arr,
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[6,5,4,3,2,1]")
	})

	Convey("Should return an empty array when there no elements in a source one", t, func() {
		arr := internal.NewArray(0)

		out, err := collections.Reverse(
			context.Background(),
			arr,
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[]")
	})
}
