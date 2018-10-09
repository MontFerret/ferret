package objects_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"

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
		So(isEqual(mergedObj, resultObj), ShouldEqual, true)
	})

	Convey("When keys are repeated", t, func() {
		obj1 := values.NewObjectWith(
			values.NewObjectProperty("prop1", values.NewInt(1)),
			values.NewObjectProperty("prop2", values.NewString("str")),
		)
		obj2 := values.NewObjectWith(
			values.NewObjectProperty("prop1", values.NewInt(3)),
		)
		resultObj := values.NewObjectWith(
			values.NewObjectProperty("prop1", values.NewInt(3)),
			values.NewObjectProperty("prop2", values.NewString("str")),
		)

		merged, err := objects.Merge(context.Background(), obj1, obj2)
		mergedObj := merged.(*values.Object)

		So(err, ShouldEqual, nil)
		So(isEqual(mergedObj, resultObj), ShouldEqual, true)
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
		So(isEqual(mergedObj, resultObj), ShouldEqual, true)
	})

	Convey("Merge empty array", t, func() {
		objArr := values.NewArray(0)
		resultObj := values.NewObject()

		merged, err := objects.Merge(context.Background(), objArr)
		mergedObj := merged.(*values.Object)

		So(err, ShouldEqual, nil)
		So(isEqual(mergedObj, resultObj), ShouldEqual, true)
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

func isEqual(obj1 *values.Object, obj2 *values.Object) bool {
	var val1 core.Value
	var val2 core.Value

	for _, key := range obj1.Keys() {
		val1, _ = obj1.Get(values.NewString(key))
		val2, _ = obj2.Get(values.NewString(key))
		if val1.Compare(val2) != 0 {
			return false
		}
	}
	for _, key := range obj2.Keys() {
		val1, _ = obj1.Get(values.NewString(key))
		val2, _ = obj2.Get(values.NewString(key))
		if val2.Compare(val1) != 0 {
			return false
		}
	}
	return true
}
