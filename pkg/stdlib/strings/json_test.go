package strings_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values/types"
	"github.com/MontFerret/ferret/pkg/stdlib/strings"
)

func TestJSONParse(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.JSONParse(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("It should parse none", t, func() {
		val := core.None

		b, err := val.MarshalJSON()

		So(err, ShouldBeNil)

		out, err := strings.JSONParse(
			context.Background(),
			core.NewString(string(b)),
		)

		So(err, ShouldBeNil)
		So(out.Type().Equals(types.None), ShouldBeTrue)
	})

	Convey("It should parse a string", t, func() {
		val := core.NewString("foobar")

		b, err := val.MarshalJSON()

		So(err, ShouldBeNil)

		out, err := strings.JSONParse(
			context.Background(),
			core.NewString(string(b)),
		)

		So(err, ShouldBeNil)
		So(out.Type().Equals(types.String), ShouldBeTrue)
	})

	Convey("It should parse an int", t, func() {
		val := core.NewInt(1)

		b, err := val.MarshalJSON()

		So(err, ShouldBeNil)

		out, err := strings.JSONParse(
			context.Background(),
			core.NewString(string(b)),
		)

		So(err, ShouldBeNil)
		So(out.Type().Equals(types.Float), ShouldBeTrue)
	})

	Convey("It should parse a float", t, func() {
		val := core.NewFloat(1.1)

		b, err := val.MarshalJSON()

		So(err, ShouldBeNil)

		out, err := strings.JSONParse(
			context.Background(),
			core.NewString(string(b)),
		)

		So(err, ShouldBeNil)
		So(out.Type().Equals(types.Float), ShouldBeTrue)
	})

	Convey("It should parse a boolean", t, func() {
		val := core.True

		b, err := val.MarshalJSON()

		So(err, ShouldBeNil)

		out, err := strings.JSONParse(
			context.Background(),
			core.NewString(string(b)),
		)

		So(err, ShouldBeNil)
		So(out.Type().Equals(types.Boolean), ShouldBeTrue)
	})

	Convey("It should parse an array", t, func() {
		val := internal.NewArrayWith(
			core.Int(1),
			core.Int(2),
			core.Int(3),
		)

		b, err := val.MarshalJSON()

		So(err, ShouldBeNil)

		out, err := strings.JSONParse(
			context.Background(),
			core.NewString(string(b)),
		)

		So(err, ShouldBeNil)
		So(out.Type().Equals(types.Array), ShouldBeTrue)
		So(out.String(), ShouldEqual, "[1,2,3]")
	})

	Convey("It should parse an object", t, func() {
		val := internal.NewObject()
		val.Set(core.NewString("foo"), core.NewString("bar"))

		b, err := val.MarshalJSON()

		So(err, ShouldBeNil)

		out, err := strings.JSONParse(
			context.Background(),
			core.NewString(string(b)),
		)

		So(err, ShouldBeNil)
		So(out.Type().Equals(types.Object), ShouldBeTrue)
		So(out.String(), ShouldEqual, `{"foo":"bar"}`)
	})
}

func TestJSONStringify(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.JSONStringify(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("It should serialize none", t, func() {
		out, err := strings.JSONStringify(
			context.Background(),
			core.None,
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "null")
	})

	Convey("It should serialize boolean", t, func() {
		out, err := strings.JSONStringify(
			context.Background(),
			core.False,
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "false")
	})

	Convey("It should serialize string", t, func() {
		out, err := strings.JSONStringify(
			context.Background(),
			core.NewString("foobar"),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, `"foobar"`)
	})

	Convey("It should serialize int", t, func() {
		out, err := strings.JSONStringify(
			context.Background(),
			core.NewInt(1),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, `1`)
	})

	Convey("It should serialize float", t, func() {
		out, err := strings.JSONStringify(
			context.Background(),
			core.NewFloat(1.1),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, `1.1`)
	})

	Convey("It should serialize array", t, func() {
		out, err := strings.JSONStringify(
			context.Background(),
			internal.NewArrayWith(
				core.NewString("foo"),
				core.NewString("bar"),
			),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, `["foo","bar"]`)
	})

	Convey("It should serialize object", t, func() {
		obj := internal.NewObject()
		obj.Set(core.NewString("foo"), core.NewString("bar"))

		out, err := strings.JSONStringify(
			context.Background(),
			obj,
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, `{"foo":"bar"}`)
	})

	Convey("It should serialize datetime", t, func() {
		obj, err := core.ParseDateTime("2006-01-02T15:04:05Z")

		So(err, ShouldBeNil)

		out, err := strings.JSONStringify(
			context.Background(),
			obj,
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, `"2006-01-02T15:04:05Z"`)
	})
}
