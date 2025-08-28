package collections_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/collections"
)

func TestIncludes(t *testing.T) {
	Convey("When searching in string", t, func() {
		Convey("Should find existing substring", func() {
			str := runtime.NewString("Hello World")
			needle := runtime.NewString("World")

			result, err := collections.Includes(context.Background(), str, needle)

			So(err, ShouldBeNil)
			So(result, ShouldEqual, runtime.True)
		})

		Convey("Should not find non-existing substring", func() {
			str := runtime.NewString("Hello World")
			needle := runtime.NewString("Mars")

			result, err := collections.Includes(context.Background(), str, needle)

			So(err, ShouldBeNil)
			So(result, ShouldEqual, runtime.False)
		})

		Convey("Should handle empty string needle", func() {
			str := runtime.NewString("Hello World")
			needle := runtime.NewString("")

			result, err := collections.Includes(context.Background(), str, needle)

			So(err, ShouldBeNil)
			So(result, ShouldEqual, runtime.True) // Empty string is contained in any string
		})

		Convey("Should handle empty string haystack", func() {
			str := runtime.NewString("")
			needle := runtime.NewString("test")

			result, err := collections.Includes(context.Background(), str, needle)

			So(err, ShouldBeNil)
			So(result, ShouldEqual, runtime.False)
		})

		Convey("Should convert needle to string", func() {
			str := runtime.NewString("123")
			needle := runtime.NewInt(123)

			result, err := collections.Includes(context.Background(), str, needle)

			So(err, ShouldBeNil)
			So(result, ShouldEqual, runtime.True)
		})
	})

	Convey("When searching in array", t, func() {
		Convey("Should find existing element", func() {
			arr := runtime.NewArrayWith(
				runtime.NewInt(1),
				runtime.NewInt(2),
				runtime.NewInt(3),
			)
			needle := runtime.NewInt(2)

			result, err := collections.Includes(context.Background(), arr, needle)

			So(err, ShouldBeNil)
			So(result, ShouldEqual, runtime.True)
		})

		Convey("Should not find non-existing element", func() {
			arr := runtime.NewArrayWith(
				runtime.NewInt(1),
				runtime.NewInt(2),
				runtime.NewInt(3),
			)
			needle := runtime.NewInt(4)

			result, err := collections.Includes(context.Background(), arr, needle)

			So(err, ShouldBeNil)
			So(result, ShouldEqual, runtime.False)
		})

		Convey("Should handle empty array", func() {
			arr := runtime.NewArray(0)
			needle := runtime.NewInt(1)

			result, err := collections.Includes(context.Background(), arr, needle)

			So(err, ShouldBeNil)
			So(result, ShouldEqual, runtime.False)
		})

		Convey("Should handle different types", func() {
			arr := runtime.NewArrayWith(
				runtime.NewInt(1),
				runtime.NewString("hello"),
				runtime.NewBoolean(true),
			)
			needle := runtime.NewString("hello")

			result, err := collections.Includes(context.Background(), arr, needle)

			So(err, ShouldBeNil)
			So(result, ShouldEqual, runtime.True)
		})
	})

	Convey("When searching in object", t, func() {
		Convey("Should find existing value", func() {
			obj := runtime.NewObjectWith(
				runtime.NewObjectProperty("a", runtime.NewInt(1)),
				runtime.NewObjectProperty("b", runtime.NewString("hello")),
				runtime.NewObjectProperty("c", runtime.NewBoolean(true)),
			)
			needle := runtime.NewString("hello")

			result, err := collections.Includes(context.Background(), obj, needle)

			So(err, ShouldBeNil)
			So(result, ShouldEqual, runtime.True)
		})

		Convey("Should not find non-existing value", func() {
			obj := runtime.NewObjectWith(
				runtime.NewObjectProperty("a", runtime.NewInt(1)),
				runtime.NewObjectProperty("b", runtime.NewString("hello")),
			)
			needle := runtime.NewString("world")

			result, err := collections.Includes(context.Background(), obj, needle)

			So(err, ShouldBeNil)
			So(result, ShouldEqual, runtime.False)
		})

		Convey("Should not find key as value", func() {
			obj := runtime.NewObjectWith(
				runtime.NewObjectProperty("hello", runtime.NewString("world")),
			)
			needle := runtime.NewString("hello") // This is a key, not a value

			result, err := collections.Includes(context.Background(), obj, needle)

			So(err, ShouldBeNil)
			So(result, ShouldEqual, runtime.False)
		})

		Convey("Should handle empty object", func() {
			obj := runtime.NewObject()
			needle := runtime.NewString("test")

			result, err := collections.Includes(context.Background(), obj, needle)

			So(err, ShouldBeNil)
			So(result, ShouldEqual, runtime.False)
		})
	})

	Convey("When haystack is invalid type", t, func() {
		needle := runtime.NewString("test")
		haystack := runtime.NewInt(123)

		result, err := collections.Includes(context.Background(), haystack, needle)

		So(err, ShouldBeError)
		So(result, ShouldEqual, runtime.None)
	})

	Convey("When searching with complex objects", t, func() {
		obj1 := runtime.NewObjectWith(
			runtime.NewObjectProperty("x", runtime.NewInt(1)),
		)
		obj2 := runtime.NewObjectWith(
			runtime.NewObjectProperty("x", runtime.NewInt(1)),
		)
		arr := runtime.NewArrayWith(obj1, runtime.NewString("test"))

		result, err := collections.Includes(context.Background(), arr, obj2)

		So(err, ShouldBeNil)
		// Objects with same content should be considered equal
		So(result, ShouldEqual, runtime.True)
	})

	Convey("Edge cases for string searching", t, func() {
		Convey("Should handle Unicode properly", func() {
			str := runtime.NewString("🚀 Hello 🌍")
			needle := runtime.NewString("🌍")

			result, err := collections.Includes(context.Background(), str, needle)

			So(err, ShouldBeNil)
			So(result, ShouldEqual, runtime.True)
		})

		Convey("Should handle case sensitivity", func() {
			str := runtime.NewString("Hello World")
			needle := runtime.NewString("hello")

			result, err := collections.Includes(context.Background(), str, needle)

			So(err, ShouldBeNil)
			So(result, ShouldEqual, runtime.False)
		})

		Convey("Should handle numeric conversion", func() {
			str := runtime.NewString("123.45")
			needle := runtime.NewFloat(123.45)

			result, err := collections.Includes(context.Background(), str, needle)

			So(err, ShouldBeNil)
			So(result, ShouldEqual, runtime.True)
		})
	})

	Convey("Edge cases for array/object searching", t, func() {
		Convey("Should handle nested objects", func() {
			innerObj := runtime.NewObjectWith(
				runtime.NewObjectProperty("nested", runtime.NewString("value")),
			)
			arr := runtime.NewArrayWith(runtime.NewInt(1), innerObj)

			result, err := collections.Includes(context.Background(), arr, innerObj)

			So(err, ShouldBeNil)
			So(result, ShouldEqual, runtime.True)
		})

		Convey("Should handle boolean values", func() {
			arr := runtime.NewArrayWith(
				runtime.NewBoolean(true),
				runtime.NewBoolean(false),
				runtime.NewString("true"),
			)

			result1, err1 := collections.Includes(context.Background(), arr, runtime.NewBoolean(true))
			result2, err2 := collections.Includes(context.Background(), arr, runtime.NewString("true"))

			So(err1, ShouldBeNil)
			So(result1, ShouldEqual, runtime.True)
			So(err2, ShouldBeNil)
			So(result2, ShouldEqual, runtime.True)
		})
	})
}
