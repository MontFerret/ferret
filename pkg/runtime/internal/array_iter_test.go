package internal_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestArrayIterator(t *testing.T) {
	Convey("No values", t, func() {
		ctx := context.Background()
		arr := internal.NewArray(0)
		iter := internal.NewArrayIterator(arr)

		hasNext, err := iter.HasNext(ctx)

		So(err, ShouldBeNil)
		So(hasNext, ShouldBeFalse)
	})

	Convey("One value", t, func() {
		ctx := context.Background()
		arr := internal.NewArray(1)
		arr.Add(ctx, core.NewInt(1))
		iter := internal.NewArrayIterator(arr)

		hasNext, err := iter.HasNext(ctx)

		So(err, ShouldBeNil)
		So(hasNext, ShouldBeTrue)

		val, key, err := iter.Next(ctx)

		So(err, ShouldBeNil)
		So(val, ShouldEqual, core.NewInt(1))
		So(key, ShouldEqual, core.NewInt(0))

		hasNext, err = iter.HasNext(ctx)

		So(err, ShouldBeNil)
		So(hasNext, ShouldBeFalse)
	})

	Convey("Multiple values", t, func() {
		ctx := context.Background()
		arr := internal.NewArray(5)
		arr.Push(core.NewInt(1))
		arr.Push(core.NewInt(2))
		arr.Push(core.NewInt(3))
		arr.Push(core.NewInt(4))
		arr.Push(core.NewInt(5))
		iter := internal.NewArrayIterator(arr)

		actual := make([]core.Int, 0, 5)

		for {
			hasNext, err := iter.HasNext(ctx)
			if !hasNext || err != nil {
				break
			}
			val, _, err := iter.Next(ctx)
			actual = append(actual, val.(core.Int))
		}

		So(actual, ShouldResemble, []core.Int{1, 2, 3, 4, 5})
	})
}

func BenchmarkArrayIterator(b *testing.B) {
	size := 100
	arr := internal.NewArray(size)

	for i := 0; i < size; i++ {
		arr.Push(core.NewInt(i))
	}

	ctx := context.Background()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		iter := internal.NewArrayIterator(arr)

		for {
			hasNext, err := iter.HasNext(ctx)
			if !hasNext || err != nil {
				break
			}

			iter.Next(ctx)
		}
	}
}
