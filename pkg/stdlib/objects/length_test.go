package objects_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/objects"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLength(t *testing.T) {
	Convey("Should return number 0", t, func() {
		obj := values.NewObject()

		objLen, err := objects.Length(context.Background(), obj)

		So(err, ShouldEqual, nil)

		len := objLen.(values.Int)

		So(len, ShouldEqual, values.NewInt(0))
		So(int(len), ShouldEqual, int(0))
	})

	Convey("Should return number 2", t, func() {
		obj := values.NewObjectWith(
			values.NewObjectProperty("Int", values.NewInt(1)),
			values.NewObjectProperty("String", values.NewString("one")),
		)

		objLen, err := objects.Length(context.Background(), obj)

		So(err, ShouldEqual, nil)

		len := objLen.(values.Int)

		So(len, ShouldEqual, values.NewInt(2))
		So(int(len), ShouldEqual, int(2))
	})

	Convey("When argument isn't object", t, func() {
		notObj := values.NewInt(1)

		len, err := objects.Length(context.Background(), notObj)

		So(err, ShouldBeError)
		So(len, ShouldEqual, values.None)
	})

	Convey("Not enought arguments", t, func() {
		len, err := objects.Length(context.Background())

		So(err, ShouldBeError)
		So(len, ShouldEqual, values.None)
	})

	Convey("Too many arguments", t, func() {
		obj := values.NewObject()
		arg1 := values.NewInt(1)
		arg2 := values.NewString("str")

		len, err := objects.Length(context.Background(), obj, arg1, arg2)

		So(err, ShouldBeError)
		So(len, ShouldEqual, values.None)
	})
}
