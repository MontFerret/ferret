package values_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func TestRangeIterator(t *testing.T) {
	Convey("Zero value", t, func() {
		ctx := context.Background()
		r := values.NewRange(0, 0)
		iter := values.NewRangeIterator(r)

		hasNext, err := iter.HasNext(ctx)
		So(err, ShouldBeNil)
		So(hasNext, ShouldBeTrue)

		So(iter.Next(ctx), ShouldBeNil)
		So(iter.Value(), ShouldEqual, values.NewInt(0))
		So(iter.Key(), ShouldEqual, values.NewInt(0))

		hasNext, err = iter.HasNext(ctx)
		So(err, ShouldBeNil)
		So(hasNext, ShouldBeFalse)
	})

	Convey("Two values", t, func() {
		ctx := context.Background()
		r := values.NewRange(0, 1)
		iter := values.NewRangeIterator(r)

		hasNext, err := iter.HasNext(ctx)
		So(err, ShouldBeNil)
		So(hasNext, ShouldBeTrue)

		So(iter.Next(ctx), ShouldBeNil)
		So(iter.Value(), ShouldEqual, values.NewInt(0))
		So(iter.Key(), ShouldEqual, values.NewInt(0))

		So(iter.Next(ctx), ShouldBeNil)
		So(iter.Value(), ShouldEqual, values.NewInt(1))
		So(iter.Key(), ShouldEqual, values.NewInt(1))

		hasNext, err = iter.HasNext(ctx)
		So(err, ShouldBeNil)
		So(hasNext, ShouldBeFalse)
	})

	Convey("Two values (2)", t, func() {
		ctx := context.Background()
		r := values.NewRange(1, 2)
		iter := values.NewRangeIterator(r)

		hasNext, err := iter.HasNext(ctx)
		So(err, ShouldBeNil)
		So(hasNext, ShouldBeTrue)

		So(iter.Next(ctx), ShouldBeNil)
		So(iter.Value(), ShouldEqual, values.NewInt(1))
		So(iter.Key(), ShouldEqual, values.NewInt(0))

		So(iter.Next(ctx), ShouldBeNil)
		So(iter.Value(), ShouldEqual, values.NewInt(2))
		So(iter.Key(), ShouldEqual, values.NewInt(1))

		hasNext, err = iter.HasNext(ctx)
		So(err, ShouldBeNil)
		So(hasNext, ShouldBeFalse)
	})

	Convey("Multiple ascending values", t, func() {
		ctx := context.Background()
		r := values.NewRange(0, 10)
		iter := values.NewRangeIterator(r)

		actual := make([]values.Int, 0, 10)

		for {
			hasNext, err := iter.HasNext(ctx)
			if !hasNext || err != nil {
				break
			}
			err = iter.Next(ctx)
			actual = append(actual, iter.Value().(values.Int))
		}

		So(actual, ShouldResemble, []values.Int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	})

	Convey("Multiple descending values", t, func() {
		ctx := context.Background()
		r := values.NewRange(10, 0)
		iter := values.NewRangeIterator(r)

		actual := make([]values.Int, 0, 10)

		for {
			hasNext, err := iter.HasNext(ctx)
			if !hasNext || err != nil {
				break
			}
			err = iter.Next(ctx)
			actual = append(actual, iter.Value().(values.Int))
		}

		So(actual, ShouldResemble, []values.Int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0})
	})
}

func BenchmarkRangeIterator(b *testing.B) {
	size := 100
	ctx := context.Background()
	r := values.NewRange(0, int64(size))
	iter := values.NewRangeIterator(r)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for {
			hasNext, err := iter.HasNext(ctx)
			if !hasNext || err != nil {
				break
			}
			err = iter.Next(ctx)
			iter.Value()
		}
	}
}
