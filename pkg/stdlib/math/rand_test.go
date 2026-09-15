package math_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/rnd"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/math"
)

func TestRand(t *testing.T) {
	Convey("Should return pseudo-random value", t, func() {
		out, err := math.Rand(rnd.WithContext(context.Background(), rnd.NewSeed(0)))

		So(err, ShouldBeNil)
		So(out, ShouldBeLessThan, 1)
	})
}
