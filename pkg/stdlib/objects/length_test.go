package objects_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/objects"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLength(t *testing.T) {
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

	Convey("Should return number 0", t, func() {
		obj := values.NewObject()

		objLen, err := objects.Length(context.Background(), obj)

		So(err, ShouldEqual, nil)

		len := objLen.(values.Int)

		So(len, ShouldEqual, values.NewInt(0))
		So(int(len), ShouldEqual, int(0))
	})

}
