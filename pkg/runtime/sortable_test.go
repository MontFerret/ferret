package runtime_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime"
)

func TestSortList(t *testing.T) {
	Convey("Should sort Array DESC", t, func() {
		arr1 := runtime.NewArrayWith(runtime.Int(3), runtime.Int(1), runtime.Int(2), runtime.Int(6), runtime.Int(4), runtime.Int(5))

		err := runtime.SortList(context.Background(), arr1, false)
		So(err, ShouldBeNil)

		j, err := arr1.MarshalJSON()

		So(err, ShouldBeNil)
		So(string(j), ShouldEqual, `[6,5,4,3,2,1]`)
	})

	Convey("Should sort Array ASC", t, func() {
		arr1 := runtime.NewArrayWith(runtime.Int(3), runtime.Int(1), runtime.Int(2), runtime.Int(6), runtime.Int(4), runtime.Int(5))

		err := runtime.SortList(context.Background(), arr1, true)
		So(err, ShouldBeNil)

		j, err := arr1.MarshalJSON()

		So(err, ShouldBeNil)
		So(string(j), ShouldEqual, `[1,2,3,4,5,6]`)
	})
}
