package strings_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

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
				runtime.NewString("foo"),
				runtime.NewString("bar"),
				runtime.NewString("qaz"),
			)

			So(out.String(), ShouldEqual, "foobarqaz")
		})
	})

	Convey("When args are not strings", t, func() {
		Convey("Concat('foo', None, 'bar') should return 'foobar'", func() {
			out, _ := strings.Concat(
				context.Background(),
				runtime.NewString("foo"),
				runtime.None,
				runtime.NewString("bar"),
			)

			So(out.String(), ShouldEqual, "foobar")
		})
		Convey("Concat('foo', 1, false) should return 'foo1false'", func() {
			out, _ := strings.Concat(
				context.Background(),
				runtime.NewString("foo"),
				runtime.NewInt(1),
				runtime.False,
			)

			So(out.String(), ShouldEqual, "foo1false")
		})

		Convey("Concat(['foo', 'bar']) should return 'foobar'", func() {
			out, _ := strings.Concat(
				context.Background(),
				runtime.NewArrayWith(runtime.NewString("foo"), runtime.NewString("bar")),
			)

			So(out.String(), ShouldEqual, "foobar")
		})

		Convey("Concat([1,2,3]) should return '123'", func() {
			out, _ := strings.Concat(
				context.Background(),
				runtime.NewArrayWith(
					runtime.NewInt(1),
					runtime.NewInt(2),
					runtime.NewInt(3),
				),
			)

			So(out.String(), ShouldEqual, "123")
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
				runtime.NewString(","),
				runtime.NewString("foo"),
				runtime.NewString("bar"),
				runtime.NewString("qaz"),
			)

			So(out.String(), ShouldEqual, "foo,bar,qaz")
		})
	})

	Convey("When args are not strings", t, func() {
		Convey("ConcatWithSeparator(',' ['foo', 'bar', 'qaz']) should return 'foo,bar,qaz'", func() {
			out, _ := strings.ConcatWithSeparator(
				context.Background(),
				runtime.NewString(","),
				runtime.NewArrayWith(
					runtime.NewString("foo"),
					runtime.NewString("bar"),
					runtime.NewString("qaz"),
				),
			)

			So(out.String(), ShouldEqual, "foo,bar,qaz")
		})

		Convey("ConcatWithSeparator(',' ['foo', None, 'qaz']) should return 'foo,qaz'", func() {
			out, _ := strings.ConcatWithSeparator(
				context.Background(),
				runtime.NewString(","),
				runtime.NewArrayWith(
					runtime.NewString("foo"),
					runtime.None,
					runtime.NewString("qaz"),
				),
			)

			So(out.String(), ShouldEqual, "foo,qaz")
		})

		Convey("ConcatWithSeparator(',' 'foo', None, 'qaz') should return 'foo,qaz'", func() {
			out, _ := strings.ConcatWithSeparator(
				context.Background(),
				runtime.NewString(","),
				runtime.NewArrayWith(
					runtime.NewString("foo"),
					runtime.None,
					runtime.NewString("qaz"),
				),
			)

			So(out.String(), ShouldEqual, "foo,qaz")
		})
	})
}
