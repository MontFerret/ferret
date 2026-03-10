package runtime_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRangeIterator(t *testing.T) {
	Convey("Zero value", t, func() {
		ctx := context.Background()
		r := runtime.NewRange(0, 0)
		iter := runtime.NewRangeIterator(r)

		val, key, err := iter.Next(ctx)

		So(err, ShouldBeNil)
		So(val, ShouldEqual, runtime.NewInt(0))
		So(key, ShouldEqual, runtime.NewInt(0))

		_, _, err = iter.Next(ctx)
		So(errors.Is(err, io.EOF), ShouldBeTrue)
	})

	Convey("Two values", t, func() {
		ctx := context.Background()
		r := runtime.NewRange(0, 1)
		iter := runtime.NewRangeIterator(r)

		val, key, err := iter.Next(ctx)

		So(err, ShouldBeNil)
		So(val, ShouldEqual, runtime.NewInt(0))
		So(key, ShouldEqual, runtime.NewInt(0))

		val, key, err = iter.Next(ctx)

		So(err, ShouldBeNil)
		So(val, ShouldEqual, runtime.NewInt(1))
		So(key, ShouldEqual, runtime.NewInt(1))

		_, _, err = iter.Next(ctx)
		So(errors.Is(err, io.EOF), ShouldBeTrue)
	})

	Convey("Two values (2)", t, func() {
		ctx := context.Background()
		r := runtime.NewRange(1, 2)
		iter := runtime.NewRangeIterator(r)

		val, key, err := iter.Next(ctx)

		So(err, ShouldBeNil)
		So(val, ShouldEqual, runtime.NewInt(1))
		So(key, ShouldEqual, runtime.NewInt(0))

		val, key, err = iter.Next(ctx)

		So(err, ShouldBeNil)
		So(val, ShouldEqual, runtime.NewInt(2))
		So(key, ShouldEqual, runtime.NewInt(1))

		_, _, err = iter.Next(ctx)
		So(errors.Is(err, io.EOF), ShouldBeTrue)
	})

	Convey("Multiple ascending values", t, func() {
		ctx := context.Background()
		r := runtime.NewRange(0, 10)
		iter := runtime.NewRangeIterator(r)

		actual := make([]runtime.Int, 0, 10)

		for {
			val, _, err := iter.Next(ctx)
			if errors.Is(err, io.EOF) {
				break
			}
			So(err, ShouldBeNil)
			actual = append(actual, val.(runtime.Int))
		}

		So(actual, ShouldResemble, []runtime.Int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	})

	Convey("Multiple descending values", t, func() {
		ctx := context.Background()
		r := runtime.NewRange(10, 0)
		iter := runtime.NewRangeIterator(r)

		actual := make([]runtime.Int, 0, 10)

		for {
			val, _, err := iter.Next(ctx)
			if errors.Is(err, io.EOF) {
				break
			}
			So(err, ShouldBeNil)
			actual = append(actual, val.(runtime.Int))
		}

		So(actual, ShouldResemble, []runtime.Int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0})
	})
}

func BenchmarkRangeIterator(b *testing.B) {
	var size runtime.Int
	size = 100
	ctx := context.Background()
	r := runtime.NewRange(0, size)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		iter := runtime.NewRangeIterator(r)
		for {
			_, _, err := iter.Next(ctx)
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				break
			}
		}
	}
}
