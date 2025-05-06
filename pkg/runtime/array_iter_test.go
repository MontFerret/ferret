package runtime_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"
)

func TestArrayIterator(t *testing.T) {
	Convey("No values", t, func() {
		ctx := context.Background()
		arr := runtime.NewArray(0)
		iter := runtime.NewArrayIterator(arr)

		hasNext, err := iter.HasNext(ctx)

		So(err, ShouldBeNil)
		So(hasNext, ShouldBeFalse)
	})

	Convey("One value", t, func() {
		ctx := context.Background()
		arr := runtime.NewArray(1)
		arr.Add(ctx, NewInt(1))
		iter := runtime.NewArrayIterator(arr)

		hasNext, err := iter.HasNext(ctx)

		So(err, ShouldBeNil)
		So(hasNext, ShouldBeTrue)

		val, key, err := iter.Next(ctx)

		So(err, ShouldBeNil)
		So(val, ShouldEqual, NewInt(1))
		So(key, ShouldEqual, NewInt(0))

		hasNext, err = iter.HasNext(ctx)

		So(err, ShouldBeNil)
		So(hasNext, ShouldBeFalse)
	})

	Convey("Multiple values", t, func() {
		ctx := context.Background()
		arr := runtime.NewArray(5)
		arr.Push(NewInt(1))
		arr.Push(NewInt(2))
		arr.Push(NewInt(3))
		arr.Push(NewInt(4))
		arr.Push(NewInt(5))
		iter := runtime.NewArrayIterator(arr)

		actual := make([]Int, 0, 5)

		for {
			hasNext, err := iter.HasNext(ctx)
			if !hasNext || err != nil {
				break
			}
			val, _, err := iter.Next(ctx)
			actual = append(actual, val.(Int))
		}

		So(actual, ShouldResemble, []Int{1, 2, 3, 4, 5})
	})
}

func BenchmarkArrayIterator(b *testing.B) {
	size := 100
	arr := runtime.NewArray(size)

	for i := 0; i < size; i++ {
		arr.Push(NewInt(i))
	}

	ctx := context.Background()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		iter := runtime.NewArrayIterator(arr)

		for {
			hasNext, err := iter.HasNext(ctx)
			if !hasNext || err != nil {
				break
			}

			iter.Next(ctx)
		}
	}
}
