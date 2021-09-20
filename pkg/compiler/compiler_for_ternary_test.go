package compiler_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/compiler"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestForTernaryExpression(t *testing.T) {
	Convey("RETURN foo ? TRUE : (FOR i IN 1..5 RETURN i*2)", t, func() {
		c := compiler.New()

		out1, err := c.MustCompile(`
			LET foo = FALSE
			RETURN foo ? TRUE : (FOR i IN 1..5 RETURN i*2)
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `[2,4,6,8,10]`)

		out2, err := c.MustCompile(`
			LET foo = TRUE
			RETURN foo ? TRUE : (FOR i IN 1..5 RETURN i*2)
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out2), ShouldEqual, `true`)
	})

	Convey("RETURN foo ? TRUE : (FOR i IN 1..5 T::FAIL() RETURN i*2)?", t, func() {
		c := compiler.New()

		out1, err := c.MustCompile(`
			LET foo = FALSE
			RETURN foo ? TRUE : (FOR i IN 1..5 T::FAIL() RETURN i*2)?
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `null`)
	})

	Convey("RETURN foo ? (FOR i IN 1..5 RETURN i) : (FOR i IN 1..5 RETURN i*2)", t, func() {
		c := compiler.New()

		out1, err := c.MustCompile(`
			LET foo = FALSE
			RETURN foo ? (FOR i IN 1..5 RETURN i) : (FOR i IN 1..5 RETURN i*2)
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `[2,4,6,8,10]`)

		out2, err := c.MustCompile(`
			LET foo = TRUE
			RETURN foo ? (FOR i IN 1..5 RETURN i) : (FOR i IN 1..5 RETURN i*2)
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out2), ShouldEqual, `[1,2,3,4,5]`)
	})

	Convey("RETURN foo ? (FOR i IN 1..5 RETURN T::FAIL())? : (FOR i IN 1..5 RETURN T::FAIL())?", t, func() {
		c := compiler.New()

		out1, err := c.MustCompile(`
			LET foo = FALSE
			RETURN foo ? (FOR i IN 1..5 RETURN T::FAIL()) : (FOR i IN 1..5 RETURN T::FAIL())?
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `null`)

		out2, err := c.MustCompile(`
			LET foo = TRUE
			RETURN foo ? (FOR i IN 1..5 RETURN T::FAIL())? : (FOR i IN 1..5 RETURN T::FAIL())
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out2), ShouldEqual, `null`)
	})

	Convey("LET res =  foo ? TRUE : (FOR i IN 1..5 RETURN i*2)", t, func() {
		c := compiler.New()

		out1, err := c.MustCompile(`
			LET foo = FALSE
			LET res = foo ? TRUE : (FOR i IN 1..5 RETURN i*2) 
			RETURN res
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `[2,4,6,8,10]`)

		out2, err := c.MustCompile(`
			LET foo = TRUE
			LET res = foo ? TRUE : (FOR i IN 1..5 RETURN i*2)
			RETURN res
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out2), ShouldEqual, `true`)
	})

	Convey("LET res = foo ? (FOR i IN 1..5 RETURN i) : (FOR i IN 1..5 RETURN i*2)", t, func() {
		c := compiler.New()

		out1, err := c.MustCompile(`
			LET foo = FALSE
			LET res = foo ? (FOR i IN 1..5 RETURN i) : (FOR i IN 1..5 RETURN i*2)
			RETURN res
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `[2,4,6,8,10]`)

		out2, err := c.MustCompile(`
			LET foo = TRUE
			LET res = foo ? (FOR i IN 1..5 RETURN i) : (FOR i IN 1..5 RETURN i*2)
			RETURN res
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out2), ShouldEqual, `[1,2,3,4,5]`)
	})

	Convey("LET res = (FOR i IN 1..5 RETURN T::FAIL())? ? TRUE : FALSE", t, func() {
		c := compiler.New()

		out1, err := c.MustCompile(`
			LET res = (FOR i IN 1..5 RETURN T::FAIL())? ? TRUE : FALSE
			RETURN res
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `false`)

		out2, err := c.MustCompile(`
			LET res = (FOR i IN 1..5 RETURN i)? ? TRUE : FALSE
			RETURN res
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out2), ShouldEqual, `true`)
	})
}

func BenchmarkForTernary(b *testing.B) {
	p := compiler.New().MustCompile(`
			LET foo = FALSE
			RETURN foo ? TRUE : (FOR i IN 1..5 RETURN i*2)
		`)

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}
