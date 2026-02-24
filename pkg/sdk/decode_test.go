package sdk_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/sdk"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	someOthers struct {
		Other string `json:"other"`
	}
	bindParams struct {
		Name     string   `ferret:"name"`
		Age      int      `ferret:"age"`
		Count    int64    `json:"count"`
		Alias    string   `ferret:"alt"`
		City     string   `ferret:"city"`
		Tags     []string `ferret:"tags"`
		Untagged string
		Pointer  *someOthers `ferret:"pointer"`
	}
	nestedAddress struct {
		City string `ferret:"city"`
		Zip  int    `ferret:"zip"`
	}
	nestedProfile struct {
		Name    string        `ferret:"name"`
		Address nestedAddress `ferret:"address"`
	}
	nestedFriendMeta struct {
		Active bool `ferret:"active"`
	}
	nestedFriend struct {
		ID   int               `ferret:"id"`
		Tags []string          `ferret:"tags"`
		Meta *nestedFriendMeta `ferret:"meta"`
	}
	nestedPayload struct {
		Profile nestedProfile  `ferret:"profile"`
		Matrix  [][]int        `ferret:"matrix"`
		Friends []nestedFriend `ferret:"friends"`
	}
)

func TestDecode(t *testing.T) {
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
		So(obj.Set(context.Background(), runtime.NewString("untagged"), runtime.NewString("ignored")), ShouldBeNil)
		So(obj.Set(context.Background(), runtime.NewString("pointer"), runtime.NewObjectWith(map[string]runtime.Value{
			"other": runtime.NewString("value"),
		})), ShouldBeNil)

		var out bindParams
		err := sdk.Decode(obj, &out)

		So(err, ShouldBeNil)
		So(out, ShouldResemble, bindParams{
			Name:  "Alice",
			Age:   30,
			Count: 42,
			Alias: "secondary",
			City:  "Paris",
			Tags:  []string{"a", "b"},
			Pointer: &someOthers{
				Other: "value",
			},
		})
	})

	Convey("Should bind deeply nested structs and slices", t, func() {
		obj := runtime.NewObjectWith(map[string]runtime.Value{
			"profile": runtime.NewObjectWith(map[string]runtime.Value{
				"name": runtime.NewString("Alice"),
				"address": runtime.NewObjectWith(map[string]runtime.Value{
					"city": runtime.NewString("Paris"),
					"zip":  runtime.NewInt(75001),
				}),
			}),
			"matrix": runtime.NewArrayWith(
				runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2)),
				runtime.NewArrayWith(runtime.NewInt(3), runtime.NewInt(4)),
			),
			"friends": runtime.NewArrayWith(
				runtime.NewObjectWith(map[string]runtime.Value{
					"id":   runtime.NewInt(1),
					"tags": runtime.NewArrayWith(runtime.NewString("a"), runtime.NewString("b")),
					"meta": runtime.NewObjectWith(map[string]runtime.Value{
						"active": runtime.NewBoolean(true),
					}),
				}),
				runtime.NewObjectWith(map[string]runtime.Value{
					"id":   runtime.NewInt(2),
					"tags": runtime.NewArrayWith(runtime.NewString("c")),
					"meta": runtime.NewObjectWith(map[string]runtime.Value{
						"active": runtime.NewBoolean(false),
					}),
				}),
			),
		})

		var out nestedPayload
		err := sdk.Decode(obj, &out)

		So(err, ShouldBeNil)
		So(out, ShouldResemble, nestedPayload{
			Profile: nestedProfile{
				Name: "Alice",
				Address: nestedAddress{
					City: "Paris",
					Zip:  75001,
				},
			},
			Matrix: [][]int{{1, 2}, {3, 4}},
			Friends: []nestedFriend{
				{
					ID:   1,
					Tags: []string{"a", "b"},
					Meta: &nestedFriendMeta{Active: true},
				},
				{
					ID:   2,
					Tags: []string{"c"},
					Meta: &nestedFriendMeta{Active: false},
				},
			},
		})
	})

	Convey("Should reject non-pointer targets", t, func() {
		obj := runtime.NewObject()
		var out bindParams
		err := sdk.Decode(obj, out)
		So(err, ShouldNotBeNil)
	})

	Convey("Should reject nil pointer targets", t, func() {
		obj := runtime.NewObject()
		var out *bindParams
		err := sdk.Decode(obj, out)
		So(err, ShouldNotBeNil)
	})

	Convey("Should reject non-string map keys", t, func() {
		obj := runtime.NewObject()
		var out map[int]string
		err := sdk.Decode(obj, &out)
		So(err, ShouldNotBeNil)
	})
}
