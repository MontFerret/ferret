package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPow(t *testing.T) {
	Convey("Should return a value", t, func() {
		out, err := math.Pow(context.Background(), runtime.NewInt(2), runtime.NewInt(4))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 16)

		out, err = math.Pow(context.Background(), runtime.NewInt(5), runtime.NewInt(-1))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 0.2)

		out, err = math.Pow(context.Background(), runtime.NewInt(5), runtime.NewInt(0))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 1)

		out, err = math.Pow(context.Background(), runtime.NewFloat(2.5), runtime.NewFloat(2))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 6.25)
	})

	Convey("Should return error for invalid arguments", t, func() {
		// Non-numeric first argument
		out, err := math.Pow(context.Background(), runtime.NewString("invalid"), runtime.NewInt(2))

		So(err, ShouldNotBeNil)
		So(out, ShouldEqual, runtime.None)

		// Non-numeric second argument
		out, err = math.Pow(context.Background(), runtime.NewInt(2), runtime.NewString("invalid"))

		So(err, ShouldNotBeNil)
		So(out, ShouldEqual, runtime.None)

		// None values
		out, err = math.Pow(context.Background(), runtime.None, runtime.NewInt(2))

		So(err, ShouldNotBeNil)
		So(out, ShouldEqual, runtime.None)
	})
}
