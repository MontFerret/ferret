package values_test

import (
	"context"
	"slices"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func TestObjectIterator(t *testing.T) {
	Convey("No values", t, func() {
		ctx := context.Background()
		obj := values.NewObject()
		iter := values.NewObjectIterator(obj)

		hasNext, err := iter.HasNext(ctx)

		So(err, ShouldBeNil)
		So(hasNext, ShouldBeFalse)
	})

	Convey("One value", t, func() {
		ctx := context.Background()
		obj := values.NewObject()
		obj.Set(values.NewString("key"), values.NewInt(1))
		iter := values.NewObjectIterator(obj)

		hasNext, err := iter.HasNext(ctx)

		So(err, ShouldBeNil)
		So(hasNext, ShouldBeTrue)

		val, key, err := iter.Next(ctx)

		So(err, ShouldBeNil)
		So(val, ShouldEqual, values.NewInt(1))
		So(key, ShouldEqual, values.NewString("key"))

		hasNext, err = iter.HasNext(ctx)

		So(err, ShouldBeNil)
		So(hasNext, ShouldBeFalse)
	})

	Convey("Multiple values", t, func() {
		ctx := context.Background()
		obj := values.NewObject()
		obj.Set(values.NewString("key1"), values.NewInt(1))
		obj.Set(values.NewString("key2"), values.NewInt(2))
		obj.Set(values.NewString("key3"), values.NewInt(3))
		obj.Set(values.NewString("key4"), values.NewInt(4))
		obj.Set(values.NewString("key5"), values.NewInt(5))
		iter := values.NewObjectIterator(obj)

		actual := make([][2]core.Value, 0, 5)

		for {
			hasNext, err := iter.HasNext(ctx)
			if !hasNext || err != nil {
				break
			}

			val, key, err := iter.Next(ctx)

			actual = append(actual, [2]core.Value{key, val})
		}

		slices.SortStableFunc(actual, func(a, b [2]core.Value) int {
			return int(values.Compare(a[1], b[1]))
		})

		So(actual, ShouldResemble, [][2]core.Value{
			{values.NewString("key1"), values.NewInt(1)},
			{values.NewString("key2"), values.NewInt(2)},
			{values.NewString("key3"), values.NewInt(3)},
			{values.NewString("key4"), values.NewInt(4)},
			{values.NewString("key5"), values.NewInt(5)},
		})
	})
}

func BenchmarkObjectIterator(b *testing.B) {
	size := 100
	obj := values.NewObject()

	for i := 0; i < size; i++ {
		obj.Set(values.NewString("key"+values.ToString(values.NewInt(i)).String()), values.NewInt(i))
	}

	ctx := context.Background()
	iter := values.NewObjectIterator(obj)

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
