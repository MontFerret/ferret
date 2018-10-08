package objects_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/objects"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMergeObjects(t *testing.T) {
	Convey("Merge 2 objects", t, func() {
		obj1 := values.NewObjectWith(
			values.NewObjectProperty("prop1", values.NewInt(1)),
			values.NewObjectProperty("prop2", values.NewString("str")),
		)
		obj2 := values.NewObjectWith(
			values.NewObjectProperty("prop3", values.NewInt(3)),
		)

		resultObj := values.NewObjectWith(
			values.NewObjectProperty("prop1", values.NewInt(1)),
			values.NewObjectProperty("prop2", values.NewString("str")),
			values.NewObjectProperty("prop3", values.NewInt(3)),
		)

		merged, err := objects.Merge(context.Background(), obj1, obj2)
		mergedObj := merged.(*values.Object)

		So(err, ShouldEqual, nil)
		So(mergedObj.Compare(resultObj), ShouldEqual, 0)
	})

	Convey("When not enought arguments", t, func() {
		obj, err := objects.Merge(context.Background(), values.NewObject())

		So(err, ShouldBeError)
		So(obj, ShouldEqual, values.None)

		obj, err = objects.Merge(context.Background())

		So(err, ShouldBeError)
		So(obj, ShouldEqual, values.None)
	})

	Convey("When wrong argument", t, func() {
		obj, err := objects.Merge(context.Background(), values.NewInt(0))

		So(err, ShouldBeError)
		So(obj, ShouldEqual, values.None)

		obj, err = objects.Merge(context.Background(), values.NewObject(), values.NewInt(0))

		So(err, ShouldBeError)
		So(obj, ShouldEqual, values.None)
	})
}

func TestMergeArray(t *testing.T) {
	Convey("Merge array", t, func() {
		objArr := values.NewArrayWith(
			values.NewObjectWith(
				values.NewObjectProperty("prop1", values.NewInt(1)),
			),
			values.NewObjectWith(
				values.NewObjectProperty("prop2", values.NewInt(2)),
			),
		)
		resultObj := values.NewObjectWith(
			values.NewObjectProperty("prop1", values.NewInt(1)),
			values.NewObjectProperty("prop2", values.NewInt(2)),
		)

		merged, err := objects.Merge(context.Background(), objArr)
		mergedObj := merged.(*values.Object)

		So(err, ShouldEqual, nil)
		So(mergedObj.Compare(resultObj), ShouldEqual, 0)
	})

	Convey("Merge empty array", t, func() {
		objArr := values.NewArray(0)
		resultObj := values.NewObject()

		merged, err := objects.Merge(context.Background(), objArr)
		mergedObj := merged.(*values.Object)

		So(err, ShouldEqual, nil)
		So(mergedObj.Compare(resultObj), ShouldEqual, 0)
	})

	Convey("When there is not object element inside the array", t, func() {
		objArr := values.NewArrayWith(
			values.NewObject(),
			values.NewInt(0),
		)

		obj, err := objects.Merge(context.Background(), objArr)

		So(err, ShouldBeError)
		So(obj, ShouldEqual, values.None)
	})

	Convey("When not enought arguments", t, func() {
		obj, err := objects.Merge(context.Background())

		So(err, ShouldBeError)
		So(obj, ShouldEqual, values.None)
	})

	Convey("When too many arguments", t, func() {
		obj, err := objects.Merge(context.Background(), values.NewArray(0), values.NewArray(0))

		So(err, ShouldBeError)
		So(obj, ShouldEqual, values.None)
	})

	Convey("When argument isn't array", t, func() {
		obj, err := objects.Merge(context.Background(), values.NewInt(0))

		So(err, ShouldBeError)
		So(obj, ShouldEqual, values.None)
	})
}
