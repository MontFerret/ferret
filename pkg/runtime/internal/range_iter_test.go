package internal_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRangeIterator(t *testing.T) {
	Convey("Zero value", t, func() {
		ctx := context.Background()
		r := NewRange(0, 0)
		iter := NewRangeIterator(r)

		hasNext, err := iter.HasNext(ctx)
		So(err, ShouldBeNil)
		So(hasNext, ShouldBeTrue)

		val, key, err := iter.Next(ctx)

		So(err, ShouldBeNil)
		So(val, ShouldEqual, core.NewInt(0))
		So(key, ShouldEqual, core.NewInt(0))

		hasNext, err = iter.HasNext(ctx)
		So(err, ShouldBeNil)
		So(hasNext, ShouldBeFalse)
	})

	Convey("Two values", t, func() {
		ctx := context.Background()
		r := NewRange(0, 1)
		iter := NewRangeIterator(r)

		hasNext, err := iter.HasNext(ctx)
		So(err, ShouldBeNil)
		So(hasNext, ShouldBeTrue)

		val, key, err := iter.Next(ctx)

		So(err, ShouldBeNil)
		So(val, ShouldEqual, core.NewInt(0))
		So(key, ShouldEqual, core.NewInt(0))

		val, key, err = iter.Next(ctx)

		So(err, ShouldBeNil)
		So(val, ShouldEqual, core.NewInt(1))
		So(key, ShouldEqual, core.NewInt(1))

		hasNext, err = iter.HasNext(ctx)
		So(err, ShouldBeNil)
		So(hasNext, ShouldBeFalse)
	})

	Convey("Two values (2)", t, func() {
		ctx := context.Background()
		r := NewRange(1, 2)
		iter := NewRangeIterator(r)

		hasNext, err := iter.HasNext(ctx)
		So(err, ShouldBeNil)
		So(hasNext, ShouldBeTrue)

		val, key, err := iter.Next(ctx)

		So(err, ShouldBeNil)
		So(val, ShouldEqual, core.NewInt(1))
		So(key, ShouldEqual, core.NewInt(0))

		val, key, err = iter.Next(ctx)

		So(err, ShouldBeNil)
		So(val, ShouldEqual, core.NewInt(2))
		So(key, ShouldEqual, core.NewInt(1))

		hasNext, err = iter.HasNext(ctx)
		So(err, ShouldBeNil)
		So(hasNext, ShouldBeFalse)
	})

	Convey("Multiple ascending values", t, func() {
		ctx := context.Background()
		r := NewRange(0, 10)
		iter := NewRangeIterator(r)

		actual := make([]core.Int, 0, 10)

		for {
			hasNext, err := iter.HasNext(ctx)
			if !hasNext || err != nil {
				break
			}

			val, _, err := iter.Next(ctx)
			actual = append(actual, val.(core.Int))
		}

		So(actual, ShouldResemble, []core.Int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	})

	Convey("Multiple descending values", t, func() {
		ctx := context.Background()
		r := NewRange(10, 0)
		iter := NewRangeIterator(r)

		actual := make([]core.Int, 0, 10)

		for {
			hasNext, err := iter.HasNext(ctx)
			if !hasNext || err != nil {
				break
			}

			val, _, err := iter.Next(ctx)
			actual = append(actual, val.(core.Int))
		}

		So(actual, ShouldResemble, []core.Int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0})
	})
}

func BenchmarkRangeIterator(b *testing.B) {
	size := 100
	ctx := context.Background()
	r := NewRange(0, int64(size))
	iter := NewRangeIterator(r)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for {
			hasNext, err := iter.HasNext(ctx)
			if !hasNext || err != nil {
				break
			}
			iter.Next(ctx)
		}
	}
}
