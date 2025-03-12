package objects_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"testing"

	"github.com/MontFerret/ferret/pkg/stdlib/objects"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMerge(t *testing.T) {
	Convey("When not enough arguments", t, func() {
		obj, err := objects.Merge(context.Background())

		So(err, ShouldBeError)
		So(obj.Compare(core.None), ShouldEqual, 0)
	})

	Convey("When wrong type of arguments", t, func() {
		obj, err := objects.Merge(context.Background(), core.NewInt(0))

		So(err, ShouldBeError)
		So(obj.Compare(core.None), ShouldEqual, 0)

		obj, err = objects.Merge(context.Background(), internal.NewObject(), core.NewInt(0))

		So(err, ShouldBeError)
		So(obj.Compare(core.None), ShouldEqual, 0)
	})

	Convey("When too many arrays", t, func() {
		obj, err := objects.Merge(context.Background(), internal.NewArray(0), internal.NewArray(0))

		So(err, ShouldBeError)
		So(obj.Compare(core.None), ShouldEqual, 0)
	})

	Convey("Merged object should be independent of source objects", t, func() {
		obj1 := internal.NewObjectWith(
			internal.NewObjectProperty("prop1", core.NewInt(1)),
			internal.NewObjectProperty("prop2", core.NewString("str")),
		)
		obj2 := internal.NewObjectWith(
			internal.NewObjectProperty("prop3", core.NewInt(3)),
		)

		result := internal.NewObjectWith(
			internal.NewObjectProperty("prop1", core.NewInt(1)),
			internal.NewObjectProperty("prop2", core.NewString("str")),
			internal.NewObjectProperty("prop3", core.NewInt(3)),
		)

		merged, err := objects.Merge(context.Background(), obj1, obj2)

		So(err, ShouldBeNil)

		obj1.Remove(core.NewString("prop1"))

		So(merged.Compare(result), ShouldEqual, 0)
	})
}

func TestMergeObjects(t *testing.T) {
	Convey("Merge single object", t, func() {
		obj1 := internal.NewObjectWith(
			internal.NewObjectProperty("prop1", core.NewInt(1)),
			internal.NewObjectProperty("prop2", core.NewString("str")),
		)
		result := internal.NewObjectWith(
			internal.NewObjectProperty("prop1", core.NewInt(1)),
			internal.NewObjectProperty("prop2", core.NewString("str")),
		)

		merged, err := objects.Merge(context.Background(), obj1)

		So(err, ShouldBeNil)
		So(merged.Compare(result), ShouldEqual, 0)
	})

	Convey("Merge two objects", t, func() {
		obj1 := internal.NewObjectWith(
			internal.NewObjectProperty("prop1", core.NewInt(1)),
			internal.NewObjectProperty("prop2", core.NewString("str")),
		)
		obj2 := internal.NewObjectWith(
			internal.NewObjectProperty("prop3", core.NewInt(3)),
		)

		result := internal.NewObjectWith(
			internal.NewObjectProperty("prop1", core.NewInt(1)),
			internal.NewObjectProperty("prop2", core.NewString("str")),
			internal.NewObjectProperty("prop3", core.NewInt(3)),
		)

		merged, err := objects.Merge(context.Background(), obj1, obj2)

		So(err, ShouldBeNil)
		So(merged.Compare(result), ShouldEqual, 0)
	})

	Convey("When keys are repeated", t, func() {
		obj1 := internal.NewObjectWith(
			internal.NewObjectProperty("prop1", core.NewInt(1)),
			internal.NewObjectProperty("prop2", core.NewString("str")),
		)
		obj2 := internal.NewObjectWith(
			internal.NewObjectProperty("prop1", core.NewInt(3)),
		)
		result := internal.NewObjectWith(
			internal.NewObjectProperty("prop1", core.NewInt(3)),
			internal.NewObjectProperty("prop2", core.NewString("str")),
		)

		merged, err := objects.Merge(context.Background(), obj1, obj2)

		So(err, ShouldBeNil)
		So(merged.Compare(result), ShouldEqual, 0)
	})
}

func TestMergeArray(t *testing.T) {
	Convey("Merge array", t, func() {
		objArr := internal.NewArrayWith(
			internal.NewObjectWith(
				internal.NewObjectProperty("prop1", core.NewInt(1)),
			),
			internal.NewObjectWith(
				internal.NewObjectProperty("prop2", core.NewInt(2)),
			),
		)
		result := internal.NewObjectWith(
			internal.NewObjectProperty("prop1", core.NewInt(1)),
			internal.NewObjectProperty("prop2", core.NewInt(2)),
		)

		merged, err := objects.Merge(context.Background(), objArr)

		So(err, ShouldBeNil)
		So(merged.Compare(result), ShouldEqual, 0)
	})

	Convey("Merge empty array", t, func() {
		objArr := internal.NewArray(0)
		result := internal.NewObject()

		merged, err := objects.Merge(context.Background(), objArr)

		So(err, ShouldBeNil)
		So(merged.Compare(result), ShouldEqual, 0)
	})

	Convey("When there is not object element inside the array", t, func() {
		objArr := internal.NewArrayWith(
			internal.NewObject(),
			core.NewInt(0),
		)

		obj, err := objects.Merge(context.Background(), objArr)

		So(err, ShouldBeError)
		So(obj.Compare(core.None), ShouldEqual, 0)
	})
}
