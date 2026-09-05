package encoding_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/encoding"
)

func TestJSONParse(t *testing.T) {
	Convey("When invalid JSON", t, func() {
		Convey("It should return an error", func() {
			_, err := encoding.JSONParse(
				context.Background(),
				runtime.NewString("invalid json"),
			)

			So(err, ShouldBeError)
		})
	})

	Convey("Should parse null", t, func() {
		out, err := encoding.JSONParse(
			context.Background(),
			runtime.NewString("null"),
		)

		So(err, ShouldBeNil)
		So(out, ShouldEqual, runtime.None)
	})

	Convey("Should parse boolean", t, func() {
		out, err := encoding.JSONParse(
			context.Background(),
			runtime.NewString("true"),
		)

		So(err, ShouldBeNil)
		So(out, ShouldEqual, runtime.True)
	})

	Convey("Should parse string", t, func() {
		out, err := encoding.JSONParse(
			context.Background(),
			runtime.NewString(`"hello"`),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "hello")
	})

}

func TestJSONStringify(t *testing.T) {
	Convey("It should serialize none", t, func() {
		out, err := encoding.JSONStringify(
			context.Background(),
			runtime.None,
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "null")
	})

	Convey("It should serialize boolean", t, func() {
		out, err := encoding.JSONStringify(
			context.Background(),
			runtime.False,
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "false")
	})

	Convey("It should serialize string", t, func() {
		out, err := encoding.JSONStringify(
			context.Background(),
			runtime.NewString("foobar"),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, `"foobar"`)
	})

	Convey("It should serialize int", t, func() {
		out, err := encoding.JSONStringify(
			context.Background(),
			runtime.NewInt(1),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, `1`)
	})

	Convey("It should serialize float", t, func() {
		out, err := encoding.JSONStringify(
			context.Background(),
			runtime.NewFloat(1.1),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, `1.1`)
	})

	Convey("It should serialize array", t, func() {
		out, err := encoding.JSONStringify(
			context.Background(),
			runtime.NewArrayWith(
				runtime.NewString("foo"),
				runtime.NewString("bar"),
			),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, `["foo","bar"]`)
	})

	Convey("It should serialize datetime", t, func() {
		obj, err := runtime.ParseDateTime("2006-01-02T15:04:05Z")

		So(err, ShouldBeNil)

		out, err := encoding.JSONStringify(
			context.Background(),
			obj,
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, `"2006-01-02T15:04:05Z"`)
	})
}
