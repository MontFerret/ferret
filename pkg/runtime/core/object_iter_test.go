package core_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"slices"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestObjectIterator(t *testing.T) {
	Convey("No values", t, func() {
		ctx := context.Background()
		obj := core.NewObject()
		iter := core.NewObjectIterator(obj)

		hasNext, err := iter.HasNext(ctx)

		So(err, ShouldBeNil)
		So(hasNext, ShouldBeFalse)
	})

	Convey("One value", t, func() {
		ctx := context.Background()
		obj := core.NewObject()
		obj.Set(NewString("key"), NewInt(1))
		iter := core.NewObjectIterator(obj)

		hasNext, err := iter.HasNext(ctx)

		So(err, ShouldBeNil)
		So(hasNext, ShouldBeTrue)

		val, key, err := iter.Next(ctx)

		So(err, ShouldBeNil)
		So(val, ShouldEqual, NewInt(1))
		So(key, ShouldEqual, NewString("key"))

		hasNext, err = iter.HasNext(ctx)

		So(err, ShouldBeNil)
		So(hasNext, ShouldBeFalse)
	})

	Convey("Multiple values", t, func() {
		ctx := context.Background()
		obj := core.NewObject()
		obj.Set(NewString("key1"), NewInt(1))
		obj.Set(NewString("key2"), NewInt(2))
		obj.Set(NewString("key3"), NewInt(3))
		obj.Set(NewString("key4"), NewInt(4))
		obj.Set(NewString("key5"), NewInt(5))
		iter := core.NewObjectIterator(obj)

		actual := make([][2]Value, 0, 5)

		for {
			hasNext, err := iter.HasNext(ctx)
			if !hasNext || err != nil {
				break
			}

			val, key, err := iter.Next(ctx)

			actual = append(actual, [2]Value{key, val})
		}

		slices.SortStableFunc(actual, func(a, b [2]Value) int {
			return int(CompareValues(a[1], b[1]))
		})

		So(actual, ShouldResemble, [][2]Value{
			{NewString("key1"), NewInt(1)},
			{NewString("key2"), NewInt(2)},
			{NewString("key3"), NewInt(3)},
			{NewString("key4"), NewInt(4)},
			{NewString("key5"), NewInt(5)},
		})
	})
}

func BenchmarkObjectIterator(b *testing.B) {
	size := 100
	obj := core.NewObject()

	for i := 0; i < size; i++ {
		obj.Set(NewString("key"+ToString(NewInt(i)).String()), NewInt(i))
	}

	ctx := context.Background()
	iter := core.NewObjectIterator(obj)

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
