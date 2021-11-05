package compiler_test

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestWaitforEventWithinTernaryExpression(t *testing.T) {
	SkipConvey("RETURN foo ? TRUE : (WAITFOR EVENT \"event\" IN obj)", t, func() {
		c := newCompilerWithObservable()

		out1, err := c.MustCompile(`
			LET foo = FALSE
			LET obj = X::CREATE()
			X::EMIT_WITH(obj, "event", "data", 100)

			RETURN foo ? TRUE : (WAITFOR EVENT "event" IN obj)
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `"data"`)

		out2, err := c.MustCompile(`
			LET foo = TRUE
			LET obj = X::CREATE()

			RETURN foo ? TRUE : (WAITFOR EVENT "event" IN obj)
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out2), ShouldEqual, `true`)
	})

	Convey("RETURN foo ? (WAITFOR EVENT \"event1\" IN obj) : (WAITFOR EVENT \"event2\" IN obj)", t, func() {
		c := newCompilerWithObservable()

		out1, err := c.MustCompile(`
			LET foo = FALSE
			LET obj = X::CREATE()
			X::EMIT_WITH(obj, "event2", "data2", 100)

			RETURN foo ? (WAITFOR EVENT "event1" IN obj) : (WAITFOR EVENT "event2" IN obj)
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `"data2"`)

		out2, err := c.MustCompile(`
			LET foo = TRUE
			LET obj = X::CREATE()
			X::EMIT_WITH(obj, "event1", "data1", 100)

			RETURN foo ? (WAITFOR EVENT "event1" IN obj) : (WAITFOR EVENT "event2" IN obj)
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out2), ShouldEqual, `"data1"`)
	})

	Convey("RETURN foo ? (FOR i IN 1..3 RETURN i*2) : (WAITFOR EVENT \"event2\" IN obj)", t, func() {
		c := newCompilerWithObservable()

		out1, err := c.MustCompile(`
			LET foo = FALSE
			LET obj = X::CREATE()
			X::EMIT_WITH(obj, "event", "data", 100)

			RETURN foo ? (FOR i IN 1..3 RETURN i*2) : (WAITFOR EVENT "event" IN obj)
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `"data"`)

		out2, err := c.MustCompile(`
			LET foo = TRUE
			LET obj = X::CREATE()

			RETURN foo ? (FOR i IN 1..3 RETURN i*2) : (WAITFOR EVENT "event" IN obj)
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out2), ShouldEqual, `[2,4,6]`)
	})

	Convey("RETURN foo ? (WAITFOR EVENT \"event\" IN obj) : (FOR i IN 1..3 RETURN i*2) ", t, func() {
		c := newCompilerWithObservable()

		out1, err := c.MustCompile(`
			LET foo = FALSE
			LET obj = X::CREATE()

			RETURN foo ? (WAITFOR EVENT "event" IN obj) : (FOR i IN 1..3 RETURN i*2)
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `[2,4,6]`)

		out2, err := c.MustCompile(`
			LET foo = TRUE
			LET obj = X::CREATE()
			X::EMIT_WITH(obj, "event", "data", 100)

			RETURN foo ? (WAITFOR EVENT "event" IN obj) : (FOR i IN 1..3 RETURN i*2)
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out2), ShouldEqual, `"data"`)
	})
}
