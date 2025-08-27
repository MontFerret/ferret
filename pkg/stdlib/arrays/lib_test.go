package arrays_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestToUniqueList(t *testing.T) {
	Convey("Should return unique items from a list", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(3),
			runtime.NewInt(5),
			runtime.NewInt(6),
			runtime.NewInt(5),
			runtime.NewInt(6),
		)

		result, err := arrays.ToUniqueList(context.Background(), arr)

		So(err, ShouldBeNil)
		So(result.String(), ShouldEqual, `[1,2,3,4,5,6]`)
	})

	Convey("Should return empty list when input is empty", t, func() {
		arr := runtime.NewArrayWith()

		result, err := arrays.ToUniqueList(context.Background(), arr)

		So(err, ShouldBeNil)
		So(result.String(), ShouldEqual, `[]`)
	})

	Convey("Should preserve first occurrence of duplicate items", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewString("a"),
			runtime.NewString("b"),
			runtime.NewString("a"),
			runtime.NewString("c"),
			runtime.NewString("b"),
		)

		result, err := arrays.ToUniqueList(context.Background(), arr)

		So(err, ShouldBeNil)
		So(result.String(), ShouldEqual, `["a","b","c"]`)
	})

	Convey("Should handle mixed data types", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewString("test"),
			runtime.NewBoolean(true),
			runtime.NewInt(1),
			runtime.NewString("test"),
			runtime.NewBoolean(false),
		)

		result, err := arrays.ToUniqueList(context.Background(), arr)

		So(err, ShouldBeNil)
		So(result.String(), ShouldEqual, `[1,"test",true,false]`)
	})
}