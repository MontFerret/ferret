package compiler_test

import (
	"context"
	"errors"
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestLogicalOperators(t *testing.T) {
	Convey("Should compile RETURN 2 > 1 AND 1 > 0", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			RETURN 2 > 1 AND 1 > 0
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "true")
	})

	Convey("Should compile RETURN 2 > 1 OR 1 < 0", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			RETURN 2 > 1 OR 1 < 0
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "true")
	})

	Convey("1 || 7  should return 1", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			RETURN 1 || 7
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "1")
	})

	Convey("NONE || 'foo'  should return 'foo'", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			RETURN NONE || 'foo'
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `"foo"`)
	})

	Convey("ERROR()? || 'boo'  should return 'boo'", t, func() {
		c := compiler.New()
		c.RegisterFunction("ERROR", func(ctx context.Context, args ...core.Value) (core.Value, error) {
			return nil, errors.New("test")
		})

		p, err := c.Compile(`
			RETURN ERROR()? || 'boo'
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `"boo"`)
	})

	Convey("!ERROR()? && TRUE should return false", t, func() {
		c := compiler.New()
		c.RegisterFunction("ERROR", func(ctx context.Context, args ...core.Value) (core.Value, error) {
			return nil, errors.New("test")
		})

		p, err := c.Compile(`
			RETURN !ERROR()? && TRUE
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `true`)
	})

	Convey("NONE && true should return null", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			RETURN NONE && true
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `null`)
	})

	Convey("'' && true  should return ''", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			RETURN '' && true
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `""`)
	})

	Convey("true && 23  should return '23", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			RETURN true && 23 
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `23`)
	})

	Convey("NOT TRUE should return false", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			RETURN NOT TRUE
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `false`)
	})

	Convey("NOT u.valid should return true", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			LET u = { valid: false }

			RETURN NOT u.valid
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `true`)
	})
}

func BenchmarkLogicalOperatorsAnd(b *testing.B) {
	p := compiler.New().MustCompile(`
			RETURN 2 > 1 AND 1 > 0
		`)

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}

func BenchmarkLogicalOperatorsOr(b *testing.B) {
	p := compiler.New().MustCompile(`
			RETURN 2 > 1 OR 1 < 0
		`)

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}
