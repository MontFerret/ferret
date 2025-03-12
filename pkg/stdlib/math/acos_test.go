package math_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"testing"

	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAcos(t *testing.T) {
	Convey("Should return arccosine", t, func() {
		out, err := math.Acos(context.Background(), core.NewInt(-1))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 3.141592653589793)

		out, err = math.Acos(context.Background(), core.NewInt(0))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 1.5707963267948966)
	})
}
