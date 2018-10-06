package objects

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLength(t *testing.T) {
	Convey("Should return number 2", t, func() {
		obj := values.NewObjectWith(
			values.NewObjectProperty("Int", values.NewInt(1)),
			values.NewObjectProperty("String", values.NewString("one")),
		)
		objLen := obj.Length()

		So(objLen, ShouldEqual, values.NewInt(2))
		So(int(objLen), ShouldEqual, 2)
	})

	Convey("Should return number 0", t, func() {
		obj := values.NewObject()
		objLen := obj.Length()

		So(objLen, ShouldEqual, values.NewInt(0))
		So(int(objLen), ShouldEqual, 0)
	})

}
