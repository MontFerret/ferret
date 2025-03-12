package internal_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"slices"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

func TestObjectIterator(t *testing.T) {
	Convey("No values", t, func() {
		ctx := context.Background()
		obj := internal.NewObject()
		iter := internal.NewObjectIterator(obj)

		hasNext, err := iter.HasNext(ctx)

		So(err, ShouldBeNil)
		So(hasNext, ShouldBeFalse)
	})

	Convey("One value", t, func() {
		ctx := context.Background()
		obj := internal.NewObject()
		obj.Set(core.NewString("key"), core.NewInt(1))
		iter := internal.NewObjectIterator(obj)

		hasNext, err := iter.HasNext(ctx)

		So(err, ShouldBeNil)
		So(hasNext, ShouldBeTrue)

		val, key, err := iter.Next(ctx)

		So(err, ShouldBeNil)
		So(val, ShouldEqual, core.NewInt(1))
		So(key, ShouldEqual, core.NewString("key"))

		hasNext, err = iter.HasNext(ctx)

		So(err, ShouldBeNil)
		So(hasNext, ShouldBeFalse)
	})

	Convey("Multiple values", t, func() {
		ctx := context.Background()
		obj := internal.NewObject()
		obj.Set(core.NewString("key1"), core.NewInt(1))
		obj.Set(core.NewString("key2"), core.NewInt(2))
		obj.Set(core.NewString("key3"), core.NewInt(3))
		obj.Set(core.NewString("key4"), core.NewInt(4))
		obj.Set(core.NewString("key5"), core.NewInt(5))
		iter := internal.NewObjectIterator(obj)

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
			return int(core.CompareValues(a[1], b[1]))
		})

		So(actual, ShouldResemble, [][2]core.Value{
			{core.NewString("key1"), core.NewInt(1)},
			{core.NewString("key2"), core.NewInt(2)},
			{core.NewString("key3"), core.NewInt(3)},
			{core.NewString("key4"), core.NewInt(4)},
			{core.NewString("key5"), core.NewInt(5)},
		})
	})
}

func BenchmarkObjectIterator(b *testing.B) {
	size := 100
	obj := internal.NewObject()

	for i := 0; i < size; i++ {
		obj.Set(core.NewString("key"+ToString(core.NewInt(i)).String()), core.NewInt(i))
	}

	ctx := context.Background()
	iter := internal.NewObjectIterator(obj)

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
