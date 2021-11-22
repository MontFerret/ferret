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
			LET obj = X::VAL("event", ["data"])

			RETURN foo ? TRUE : (WAITFOR EVENT "event" IN obj)
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `"data"`)

		out2, err := c.MustCompile(`
			LET foo = TRUE
			LET obj = X::VAL("event", ["data"])

			RETURN foo ? TRUE : (WAITFOR EVENT "event" IN obj)
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out2), ShouldEqual, `true`)
	})

	SkipConvey("RETURN foo ? (WAITFOR EVENT \"event1\" IN obj) : (WAITFOR EVENT \"event2\" IN obj)", t, func() {
		c := newCompilerWithObservable()

		out1, err := c.MustCompile(`
			LET foo = FALSE
			LET obj = X::VAL("event2", ["data2"])

			RETURN foo ? (WAITFOR EVENT "event1" IN obj) : (WAITFOR EVENT "event2" IN obj)
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `"data2"`)

		c = newCompilerWithObservable()
		out2, err := c.MustCompile(`
			LET foo = TRUE
			LET obj = X::VAL("event1", ["data1"])

			RETURN foo ? (WAITFOR EVENT "event1" IN obj) : (WAITFOR EVENT "event2" IN obj)
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out2), ShouldEqual, `"data1"`)
	})

	SkipConvey("RETURN foo ? (FOR i IN 1..3 RETURN i*2) : (WAITFOR EVENT \"event2\" IN obj)", t, func() {
		c := newCompilerWithObservable()

		out1, err := c.MustCompile(`
			LET foo = FALSE
			LET obj = X::VAL("event", ["data"])

			RETURN foo ? (FOR i IN 1..3 RETURN i*2) : (WAITFOR EVENT "event" IN obj)
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `"data"`)

		out2, err := c.MustCompile(`
			LET foo = TRUE
			LET obj = X::VAL("event", ["data"])

			RETURN foo ? (FOR i IN 1..3 RETURN i*2) : (WAITFOR EVENT "event" IN obj)
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out2), ShouldEqual, `[2,4,6]`)
	})

	SkipConvey("RETURN foo ? (WAITFOR EVENT \"event\" IN obj) : (FOR i IN 1..3 RETURN i*2) ", t, func() {
		c := newCompilerWithObservable()

		out1, err := c.MustCompile(`
			LET foo = FALSE
			LET obj = X::VAL("event", ["data"], 1000)

			RETURN foo ? (WAITFOR EVENT "event" IN obj) : (FOR i IN 1..3 RETURN i*2)
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `[2,4,6]`)

		out2, err := c.MustCompile(`
			LET foo = TRUE
			LET obj = X::VAL("event", ["data"])

			RETURN foo ? (WAITFOR EVENT "event" IN obj) : (FOR i IN 1..3 RETURN i*2)
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out2), ShouldEqual, `"data"`)
	})
}
