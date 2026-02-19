package runtime_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"
)

func TestArrayIterator(t *testing.T) {
	Convey("No values", t, func() {
		ctx := context.Background()
		arr := runtime.NewArray(0)
		iter := runtime.NewArrayIterator(arr)

		_, _, err := iter.Next(ctx)
		So(errors.Is(err, io.EOF), ShouldBeTrue)
	})

	Convey("One value", t, func() {
		ctx := context.Background()
		arr := runtime.NewArray(1)
		arr.Append(ctx, runtime.NewInt(1))
		iter := runtime.NewArrayIterator(arr)

		val, key, err := iter.Next(ctx)

		So(err, ShouldBeNil)
		So(val, ShouldEqual, runtime.NewInt(1))
		So(key, ShouldEqual, runtime.NewInt(0))

		_, _, err = iter.Next(ctx)
		So(errors.Is(err, io.EOF), ShouldBeTrue)
	})

	Convey("Multiple values", t, func() {
		ctx := context.Background()
		arr := runtime.NewArray(5)
		arr.Append(ctx, runtime.NewInt(1))
		arr.Append(ctx, runtime.NewInt(2))
		arr.Append(ctx, runtime.NewInt(3))
		arr.Append(ctx, runtime.NewInt(4))
		arr.Append(ctx, runtime.NewInt(5))
		iter := runtime.NewArrayIterator(arr)

		actual := make([]runtime.Int, 0, 5)

		for {
			val, _, err := iter.Next(ctx)
			if errors.Is(err, io.EOF) {
				break
			}
			So(err, ShouldBeNil)
			actual = append(actual, val.(runtime.Int))
		}

		So(actual, ShouldResemble, []runtime.Int{1, 2, 3, 4, 5})
	})
}

func BenchmarkArrayIterator(b *testing.B) {
	size := 100
	arr := runtime.NewArray(size)
	ctx := context.Background()

	for i := 0; i < size; i++ {
		arr.Append(ctx, runtime.NewInt(i))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		iter := runtime.NewArrayIterator(arr)

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
