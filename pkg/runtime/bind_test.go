package runtime_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	. "github.com/smartystreets/goconvey/convey"
)

type bindParams struct {
	Name  string `json:"name"`
	Age   int    `ferret:"age"`
	Count int64  `json:"count"`
	Alias string `json:"alias" ferret:"alt"`
	City  string
	Tags  []string `json:"tags"`
}

func TestBind(t *testing.T) {
	Convey("Should bind values into a struct", t, func() {
		obj := runtime.NewObject()
		So(obj.Set(context.Background(), runtime.NewString("name"), runtime.NewString("Alice")), ShouldBeNil)
		So(obj.Set(context.Background(), runtime.NewString("age"), runtime.NewInt(30)), ShouldBeNil)
		So(obj.Set(context.Background(), runtime.NewString("count"), runtime.NewInt64(42)), ShouldBeNil)
		So(obj.Set(context.Background(), runtime.NewString("alias"), runtime.NewString("primary")), ShouldBeNil)
		So(obj.Set(context.Background(), runtime.NewString("alt"), runtime.NewString("secondary")), ShouldBeNil)
		So(obj.Set(context.Background(), runtime.NewString("CITY"), runtime.NewString("Paris")), ShouldBeNil)
		So(obj.Set(context.Background(), runtime.NewString("tags"), runtime.NewArrayWith(
			runtime.NewString("a"),
			runtime.NewString("b"),
		)), ShouldBeNil)

		var out bindParams
		err := runtime.Bind(obj, &out)

		So(err, ShouldBeNil)
		So(out, ShouldResemble, bindParams{
			Name:  "Alice",
			Age:   30,
			Count: 42,
			Alias: "primary",
			City:  "Paris",
			Tags:  []string{"a", "b"},
		})
	})

	Convey("Should reject non-pointer targets", t, func() {
		obj := runtime.NewObject()
		var out bindParams
		err := runtime.Bind(obj, out)
		So(err, ShouldNotBeNil)
	})

	Convey("Should reject nil pointer targets", t, func() {
		obj := runtime.NewObject()
		var out *bindParams
		err := runtime.Bind(obj, out)
		So(err, ShouldNotBeNil)
	})

	Convey("Should reject non-string map keys", t, func() {
		obj := runtime.NewObject()
		var out map[int]string
		err := runtime.Bind(obj, &out)
		So(err, ShouldNotBeNil)
	})
}
