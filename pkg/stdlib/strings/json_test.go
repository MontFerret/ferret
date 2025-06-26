package strings_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/strings"
)

// TODO: Fix the tests
func TestJSONParse(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.JSONParse(context.Background())

			So(err, ShouldBeError)
		})
	})

	//Convey("It should parse none", t, func() {
	//	val := runtime.None
	//
	//	b, err := val.MarshalJSON()
	//
	//	So(err, ShouldBeNil)
	//
	//	out, err := strings.JSONParse(
	//		context.Background(),
	//		runtime.NewString(string(b)),
	//	)
	//
	//	So(err, ShouldBeNil)
	//	So(out.Type().Equals(types.None), ShouldBeTrue)
	//})

	//Convey("It should parse a string", t, func() {
	//	val := runtime.NewString("foobar")
	//
	//	b, err := val.MarshalJSON()
	//
	//	So(err, ShouldBeNil)
	//
	//	out, err := strings.JSONParse(
	//		context.Background(),
	//		runtime.NewString(string(b)),
	//	)
	//
	//	So(err, ShouldBeNil)
	//	So(out.Type().Equals(types.String), ShouldBeTrue)
	//})
	//
	//Convey("It should parse an int", t, func() {
	//	val := runtime.NewInt(1)
	//
	//	b, err := val.MarshalJSON()
	//
	//	So(err, ShouldBeNil)
	//
	//	out, err := strings.JSONParse(
	//		context.Background(),
	//		runtime.NewString(string(b)),
	//	)
	//
	//	So(err, ShouldBeNil)
	//	So(out.Type().Equals(types.Float), ShouldBeTrue)
	//})
	//
	//Convey("It should parse a float", t, func() {
	//	val := runtime.NewFloat(1.1)
	//
	//	b, err := val.MarshalJSON()
	//
	//	So(err, ShouldBeNil)
	//
	//	out, err := strings.JSONParse(
	//		context.Background(),
	//		runtime.NewString(string(b)),
	//	)
	//
	//	So(err, ShouldBeNil)
	//	So(out.Type().Equals(types.Float), ShouldBeTrue)
	//})
	//
	//Convey("It should parse a boolean", t, func() {
	//	val := runtime.True
	//
	//	b, err := val.MarshalJSON()
	//
	//	So(err, ShouldBeNil)
	//
	//	out, err := strings.JSONParse(
	//		context.Background(),
	//		runtime.NewString(string(b)),
	//	)
	//
	//	So(err, ShouldBeNil)
	//	So(out.Type().Equals(types.Boolean), ShouldBeTrue)
	//})
	//
	//Convey("It should parse an array", t, func() {
	//	val := runtime.NewArrayWith(
	//		runtime.Int(1),
	//		runtime.Int(2),
	//		runtime.Int(3),
	//	)
	//
	//	b, err := val.MarshalJSON()
	//
	//	So(err, ShouldBeNil)
	//
	//	out, err := strings.JSONParse(
	//		context.Background(),
	//		runtime.NewString(string(b)),
	//	)
	//
	//	So(err, ShouldBeNil)
	//	So(out.Type().Equals(types.Array), ShouldBeTrue)
	//	So(out.String(), ShouldEqual, "[1,2,3]")
	//})
	//
	//Convey("It should parse an object", t, func() {
	//	val := runtime.NewObject()
	//	val.Set(runtime.NewString("foo"), runtime.NewString("bar"))
	//
	//	b, err := val.MarshalJSON()
	//
	//	So(err, ShouldBeNil)
	//
	//	out, err := strings.JSONParse(
	//		context.Background(),
	//		runtime.NewString(string(b)),
	//	)
	//
	//	So(err, ShouldBeNil)
	//	So(out.Type().Equals(types.Object), ShouldBeTrue)
	//	So(out.String(), ShouldEqual, `{"foo":"bar"}`)
	//})
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
			runtime.None,
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "null")
	})

	Convey("It should serialize boolean", t, func() {
		out, err := strings.JSONStringify(
			context.Background(),
			runtime.False,
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "false")
	})

	Convey("It should serialize string", t, func() {
		out, err := strings.JSONStringify(
			context.Background(),
			runtime.NewString("foobar"),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, `"foobar"`)
	})

	Convey("It should serialize int", t, func() {
		out, err := strings.JSONStringify(
			context.Background(),
			runtime.NewInt(1),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, `1`)
	})

	Convey("It should serialize float", t, func() {
		out, err := strings.JSONStringify(
			context.Background(),
			runtime.NewFloat(1.1),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, `1.1`)
	})

	Convey("It should serialize array", t, func() {
		out, err := strings.JSONStringify(
			context.Background(),
			runtime.NewArrayWith(
				runtime.NewString("foo"),
				runtime.NewString("bar"),
			),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, `["foo","bar"]`)
	})

	//Convey("It should serialize object", t, func() {
	//	obj := runtime.NewObject()
	//	obj.Set(runtime.NewString("foo"), runtime.NewString("bar"))
	//
	//	out, err := strings.JSONStringify(
	//		context.Background(),
	//		obj,
	//	)
	//
	//	So(err, ShouldBeNil)
	//	So(out.String(), ShouldEqual, `{"foo":"bar"}`)
	//})

	Convey("It should serialize datetime", t, func() {
		obj, err := runtime.ParseDateTime("2006-01-02T15:04:05Z")

		So(err, ShouldBeNil)

		out, err := strings.JSONStringify(
			context.Background(),
			obj,
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, `"2006-01-02T15:04:05Z"`)
	})
}
