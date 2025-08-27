package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRange(t *testing.T) {
	Convey("Should return range of numbers", t, func() {
		out, err := math.Range(context.Background(), runtime.NewInt(1), runtime.NewInt(4))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,3,4]")

		out, err = math.Range(context.Background(),
			runtime.NewInt(1),
			runtime.NewInt(4),
			runtime.NewInt(2))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,3]")

		out, err = math.Range(context.Background(),
			runtime.NewInt(1),
			runtime.NewInt(4),
			runtime.NewInt(3),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,4]")

		out, err = math.Range(context.Background(),
			runtime.NewFloat(1.5),
			runtime.NewFloat(2.5),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1.5,2.5]")

		out, err = math.Range(context.Background(),
			runtime.NewFloat(1.5),
			runtime.NewFloat(2.5),
			runtime.NewFloat(0.5),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1.5,2,2.5]")

		out, err = math.Range(context.Background(),
			runtime.NewFloat(-0.75),
			runtime.NewFloat(1.1),
			runtime.NewFloat(0.5),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[-0.75,-0.25,0.25,0.75]")
	})
}
