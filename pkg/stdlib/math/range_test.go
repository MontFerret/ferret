package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"

	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRange(t *testing.T) {
	Convey("Should return range of numbers", t, func() {
		out, err := math.Range(context.Background(), core.NewInt(1), core.NewInt(4))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,3,4]")

		out, err = math.Range(context.Background(),
			core.NewInt(1),
			core.NewInt(4),
			core.NewInt(2))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,3]")

		out, err = math.Range(context.Background(),
			core.NewInt(1),
			core.NewInt(4),
			core.NewInt(3),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,4]")

		out, err = math.Range(context.Background(),
			core.NewFloat(1.5),
			core.NewFloat(2.5),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1.5,2.5]")

		out, err = math.Range(context.Background(),
			core.NewFloat(1.5),
			core.NewFloat(2.5),
			core.NewFloat(0.5),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1.5,2,2.5]")

		out, err = math.Range(context.Background(),
			core.NewFloat(-0.75),
			core.NewFloat(1.1),
			core.NewFloat(0.5),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[-0.75,-0.25,0.25,0.75]")
	})
}
