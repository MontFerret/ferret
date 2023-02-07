package strings_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/strings"
)

func TestConcat(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := strings.Concat(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("When args are strings", t, func() {
		Convey("Concat('foo', 'bar', 'qaz') should return 'foobarqaz'", func() {
			out, _ := strings.Concat(
				context.Background(),
				values.NewString("foo"),
				values.NewString("bar"),
				values.NewString("qaz"),
			)

			So(out, ShouldEqual, "foobarqaz")
		})
	})

	Convey("When args are not strings", t, func() {
		Convey("Concat('foo', None, 'bar') should return 'foobar'", func() {
			out, _ := strings.Concat(
				context.Background(),
				values.NewString("foo"),
				values.None,
				values.NewString("bar"),
			)

			So(out, ShouldEqual, "foobar")
		})
		Convey("Concat('foo', 1, false) should return 'foo1false'", func() {
			out, _ := strings.Concat(
				context.Background(),
				values.NewString("foo"),
				values.NewInt(1),
				values.False,
			)

			So(out, ShouldEqual, "foo1false")
		})

		Convey("Concat(['foo', 'bar']) should return 'foobar'", func() {
			out, _ := strings.Concat(
				context.Background(),
				values.NewArrayWith(values.NewString("foo"), values.NewString("bar")),
			)

			So(out, ShouldEqual, "foobar")
		})

		Convey("Concat([1,2,3]) should return '123'", func() {
			out, _ := strings.Concat(
				context.Background(),
				values.NewArrayWith(
					values.NewInt(1),
					values.NewInt(2),
					values.NewInt(3),
				),
			)

			So(out, ShouldEqual, "123")
		})
	})
}

func TestConcatWithSeparator(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := strings.ConcatWithSeparator(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("When args are strings", t, func() {
		Convey("ConcatWithSeparator(',' 'foo', 'bar', 'qaz') should return 'foo,bar,qaz'", func() {
			out, _ := strings.ConcatWithSeparator(
				context.Background(),
				values.NewString(","),
				values.NewString("foo"),
				values.NewString("bar"),
				values.NewString("qaz"),
			)

			So(out, ShouldEqual, "foo,bar,qaz")
		})
	})

	Convey("When args are not strings", t, func() {
		Convey("ConcatWithSeparator(',' ['foo', 'bar', 'qaz']) should return 'foo,bar,qaz'", func() {
			out, _ := strings.ConcatWithSeparator(
				context.Background(),
				values.NewString(","),
				values.NewArrayWith(
					values.NewString("foo"),
					values.NewString("bar"),
					values.NewString("qaz"),
				),
			)

			So(out, ShouldEqual, "foo,bar,qaz")
		})

		Convey("ConcatWithSeparator(',' ['foo', None, 'qaz']) should return 'foo,qaz'", func() {
			out, _ := strings.ConcatWithSeparator(
				context.Background(),
				values.NewString(","),
				values.NewArrayWith(
					values.NewString("foo"),
					values.None,
					values.NewString("qaz"),
				),
			)

			So(out, ShouldEqual, "foo,qaz")
		})

		Convey("ConcatWithSeparator(',' 'foo', None, 'qaz') should return 'foo,qaz'", func() {
			out, _ := strings.ConcatWithSeparator(
				context.Background(),
				values.NewString(","),
				values.NewArrayWith(
					values.NewString("foo"),
					values.None,
					values.NewString("qaz"),
				),
			)

			So(out, ShouldEqual, "foo,qaz")
		})
	})
}
