package collections_test

import (
	"encoding/json"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGroup(t *testing.T) {
	makeObj := func(active bool, age int, city, gender string) *values.Object {
		obj := values.NewObject()

		obj.Set("active", values.NewBoolean(active))
		obj.Set("age", values.NewInt(age))
		obj.Set("city", values.NewString(city))
		obj.Set("gender", values.NewString(gender))

		return obj
	}

	Convey("Should group by a single key", t, func() {
		arr := []core.Value{
			makeObj(true, 31, "D.C.", "m"),
			makeObj(true, 29, "L.A.", "f"),
			makeObj(true, 36, "D.C.", "m"),
			makeObj(true, 34, "N.Y.C.", "f"),
			makeObj(true, 28, "L.A.", "f"),
			makeObj(true, 41, "Boston", "m"),
		}

		iter, err := collections.NewGroupIterator(
			collections.NewSliceIterator(arr),
			collections.NewGroupSelector(
				func(set collections.ResultSet) (core.Value, error) {
					val, _ := set[0].(*values.Object).Get("gender")

					return val, nil
				},
				func(set collections.ResultSet) (core.Value, error) {
					return set[0], nil
				},
			),
		)

		So(err, ShouldBeNil)

		res, err := collections.ToMap(iter)

		So(err, ShouldBeNil)

		j, _ := json.Marshal(res)

		So(string(j), ShouldEqual, `{"f":[{"active":true,"age":29,"city":"L.A.","gender":"f"},{"active":true,"age":34,"city":"N.Y.C.","gender":"f"},{"active":true,"age":28,"city":"L.A.","gender":"f"}],"m":[{"active":true,"age":31,"city":"D.C.","gender":"m"},{"active":true,"age":36,"city":"D.C.","gender":"m"},{"active":true,"age":41,"city":"Boston","gender":"m"}]}`)
	})
}
