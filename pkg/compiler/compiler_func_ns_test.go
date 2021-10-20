package compiler_test

import (
	"context"
	"fmt"
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/parser"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
	"regexp"
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

	Convey("T::FAIL()? should return NONE", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			RETURN T::FAIL()?
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `null`)
	})

	Convey("Should use keywords", t, func() {
		p := parser.New("RETURN TRUE")
		c := compiler.New()

		r := regexp.MustCompile("\\w+")

		for _, l := range p.GetLiteralNames() {
			if r.MatchString(l) {
				kw := l[1 : len(l)-1]

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
		}
	})

}
