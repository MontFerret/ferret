package compiler_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestFunctionNSCall(t *testing.T) {
	Convey("Should compile RETURN T::SPY", t, func() {
		c := compiler.New()

		var counter int
		err := c.Namespace("T").RegisterFunction("SPY", func(_ context.Context, _ ...core.Value) (core.Value, error) {
			counter++

			return values.None, nil
		})

		So(err, ShouldBeNil)

		p, err := c.Compile(`
			RETURN T::SPY()
		`)

		So(err, ShouldBeNil)

		_, err = p.Run(context.Background())

		So(err, ShouldBeNil)

		So(counter, ShouldEqual, 1)
	})

	Convey("Should compile RETURN T::UTILS::SPY", t, func() {
		c := compiler.New()

		var counter int
		err := c.Namespace("T").Namespace("UTILS").RegisterFunction("SPY", func(_ context.Context, _ ...core.Value) (core.Value, error) {
			counter++

			return values.None, nil
		})

		So(err, ShouldBeNil)

		p, err := c.Compile(`
			RETURN T::UTILS::SPY()
		`)

		So(err, ShouldBeNil)

		_, err = p.Run(context.Background())

		So(err, ShouldBeNil)

		So(counter, ShouldEqual, 1)
	})

	Convey("Should NOT compile RETURN T:UTILS::SPY", t, func() {
		c := compiler.New()

		var counter int
		err := c.Namespace("T").Namespace("UTILS").RegisterFunction("SPY", func(_ context.Context, _ ...core.Value) (core.Value, error) {
			counter++

			return values.None, nil
		})

		So(err, ShouldBeNil)

		_, err = c.Compile(`
			RETURN T:UTILS::SPY()
		`)

		So(err, ShouldNotBeNil)
	})
}
