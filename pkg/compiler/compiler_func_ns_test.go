package compiler_test

import (
	"context"
	"fmt"
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

	Convey("Should use keywords", t, func() {
		c := compiler.New()

		keywords := []string{
			"And",
			"Or",
			"For",
			"Return",
			"Distinct",
			"Filter",
			"Sort",
			"Limit",
			"Let",
			"Collect",
			"Desc",
			"Asc",
			"None",
			"Null",
			"True",
			"False",
			"Use",
			"Into",
			"Keep",
			"With",
			"Count",
			"All",
			"Any",
			"Aggregate",
			"Like",
			"Not",
			"In",
		}

		for _, kw := range keywords {
			segment := kw
			err := c.Namespace("T").Namespace(segment).RegisterFunction("TEST", func(ctx context.Context, args ...core.Value) (core.Value, error) {
				return values.True, nil
			})

			So(err, ShouldBeNil)

			err = c.Namespace("T").Namespace(segment).RegisterFunction(segment, func(ctx context.Context, args ...core.Value) (core.Value, error) {
				return values.True, nil
			})

			So(err, ShouldBeNil)

			p, err := c.Compile(fmt.Sprintf(`
			RETURN T::%s::TEST()
		`, segment))

			So(err, ShouldBeNil)

			out := p.MustRun(context.Background())

			So(string(out), ShouldEqual, "true")

			p, err = c.Compile(fmt.Sprintf(`
			RETURN T::%s::%s()
		`, segment, segment))

			So(err, ShouldBeNil)

			out = p.MustRun(context.Background())

			So(string(out), ShouldEqual, "true")
		}
	})

}
