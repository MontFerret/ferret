package math_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/stdlib/math"
	. "github.com/smartystreets/goconvey/convey"
	m "math"
	"testing"
)

func TestPi(t *testing.T) {
	Convey("Should return Pi value", t, func() {
		out, err := math.Pi(context.Background())

		So(err, ShouldBeNil)
		So(out, ShouldEqual, m.Pi)
	})
}
