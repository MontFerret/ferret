package values_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func TestArrayIterator(t *testing.T) {
	Convey("No values", t, func() {
		ctx := context.Background()
		arr := values.NewArray(0)
		iter := values.NewArrayIterator(arr)

		hasNext, err := iter.HasNext(ctx)

		So(err, ShouldBeNil)
		So(hasNext, ShouldBeFalse)
	})

	Convey("One value", t, func() {
		ctx := context.Background()
		arr := values.NewArray(1)
		arr.Push(values.NewInt(1))
		iter := values.NewArrayIterator(arr)

		hasNext, err := iter.HasNext(ctx)

		So(err, ShouldBeNil)
		So(hasNext, ShouldBeTrue)

		val, key, err := iter.Next(ctx)

		So(err, ShouldBeNil)
		So(val, ShouldEqual, values.NewInt(1))
		So(key, ShouldEqual, values.NewInt(0))

		hasNext, err = iter.HasNext(ctx)

		So(err, ShouldBeNil)
		So(hasNext, ShouldBeFalse)
	})

	Convey("Multiple values", t, func() {
		ctx := context.Background()
		arr := values.NewArray(5)
		arr.Push(values.NewInt(1))
		arr.Push(values.NewInt(2))
		arr.Push(values.NewInt(3))
		arr.Push(values.NewInt(4))
		arr.Push(values.NewInt(5))
		iter := values.NewArrayIterator(arr)

		actual := make([]values.Int, 0, 5)

		for {
			hasNext, err := iter.HasNext(ctx)
			if !hasNext || err != nil {
				break
			}
			val, _, err := iter.Next(ctx)
			actual = append(actual, val.(values.Int))
		}

		So(actual, ShouldResemble, []values.Int{1, 2, 3, 4, 5})
	})
}

func BenchmarkArrayIterator(b *testing.B) {
	size := 100
	arr := values.NewArray(size)

	for i := 0; i < size; i++ {
		arr.Push(values.NewInt(i))
	}

	ctx := context.Background()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		iter := values.NewArrayIterator(arr)

		for {
			hasNext, err := iter.HasNext(ctx)
			if !hasNext || err != nil {
				break
			}

			iter.Next(ctx)
		}
	}
}
