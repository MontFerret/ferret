package compiler_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/compiler"
)

func TestRangeOperator(t *testing.T) {
	Convey("Should compile RETURN 1..10", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
				RETURN 1..10
			`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[1,2,3,4,5,6,7,8,9,10]`)
	})

	Convey("Should compile FOR i IN 1..10 RETURN i * 2", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
				FOR i IN 1..10
					RETURN i * 2
			`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[2,4,6,8,10,12,14,16,18,20]`)
	})

	Convey("Should compile LET arr = 1..10 FOR i IN arr RETURN i * 2", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
				LET arr = 1..10
				FOR i IN arr
					RETURN i * 2
			`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[2,4,6,8,10,12,14,16,18,20]`)
	})

	Convey("Should use variables", t, func() {
		out := compiler.New().MustCompile(`
				LET max = 10
				
				FOR i IN 1..max
					RETURN i * 2
		`).MustRun(context.Background())

		So(string(out), ShouldEqual, `[2,4,6,8,10,12,14,16,18,20]`)

		out2 := compiler.New().MustCompile(`
				LET min = 1
				
				FOR i IN min..10
					RETURN i * 2
		`).MustRun(context.Background())

		So(string(out2), ShouldEqual, `[2,4,6,8,10,12,14,16,18,20]`)

		out3 := compiler.New().MustCompile(`
				LET min = 1
				LET max = 10
				
				FOR i IN min..max
					RETURN i * 2
		`).MustRun(context.Background())

		So(string(out3), ShouldEqual, `[2,4,6,8,10,12,14,16,18,20]`)
	})
}

func BenchmarkRangeOperator(b *testing.B) {
	p := compiler.New().MustCompile(`
				RETURN 1..10
			`)

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}
