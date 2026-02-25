package sdk_test

import (
	"context"
	"io"
	"testing"
	"time"

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
	EmbeddedParams struct {
		URL         string `json:"url"`
		UserAgent   string `json:"userAgent"`
		KeepCookies bool   `json:"keepCookies"`
		Charset     string `json:"charset"`
	}
	EmbeddedPageLoadParams struct {
		EmbeddedParams
		Driver  string        `json:"driver"`
		Timeout time.Duration `json:"timeout"`
	}
	EmbeddedOuterPointer struct {
		*EmbeddedParams
		Driver string `json:"driver"`
	}
	EmbeddedNode struct {
		*EmbeddedNode
	}
	closableIterable struct {
		values []runtime.Value
		closed *bool
	}
	closableIterator struct {
		values []runtime.Value
		index  int
		closed *bool
	}
)

func (c closableIterable) String() string {
	return "closableIterable"
}

func (c closableIterable) Hash() uint64 {
	return 0
}

func (c closableIterable) Copy() runtime.Value {
	return c
}

func (c closableIterable) Iterate(_ context.Context) (runtime.Iterator, error) {
	return &closableIterator{
		values: c.values,
		closed: c.closed,
	}, nil
}

func (it *closableIterator) Next(_ context.Context) (runtime.Value, runtime.Value, error) {
	if it.index >= len(it.values) {
		return runtime.None, runtime.None, io.EOF
	}

	index := it.index
	value := it.values[index]
	it.index++

	return value, runtime.NewInt(index), nil
}

func (it *closableIterator) Close() error {
	if it.closed != nil {
		*it.closed = true
	}

	return nil
}

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

	Convey("Should bind anonymous embedded structs inline", t, func() {
		obj := runtime.NewObjectWith(map[string]runtime.Value{
			"url":         runtime.NewString("https://example.test"),
			"userAgent":   runtime.NewString("agent"),
			"keepCookies": runtime.NewBoolean(true),
			"charset":     runtime.NewString("utf-8"),
			"driver":      runtime.NewString("chrome"),
			"timeout":     runtime.NewInt(42),
		})

		var out EmbeddedPageLoadParams
		err := sdk.Decode(obj, &out)

		So(err, ShouldBeNil)
		So(out, ShouldResemble, EmbeddedPageLoadParams{
			EmbeddedParams: EmbeddedParams{
				URL:         "https://example.test",
				UserAgent:   "agent",
				KeepCookies: true,
				Charset:     "utf-8",
			},
			Driver:  "chrome",
			Timeout: 42,
		})
	})

	Convey("Should allocate embedded pointers only when matched", t, func() {
		Convey("no matching keys keeps pointer nil", func() {
			obj := runtime.NewObjectWith(map[string]runtime.Value{
				"driver": runtime.NewString("chrome"),
			})

			var out EmbeddedOuterPointer
			err := sdk.Decode(obj, &out)

			So(err, ShouldBeNil)
			So(out, ShouldResemble, EmbeddedOuterPointer{
				EmbeddedParams: nil,
				Driver:         "chrome",
			})
		})

		Convey("matching key allocates pointer", func() {
			obj := runtime.NewObjectWith(map[string]runtime.Value{
				"url":    runtime.NewString("https://example.test"),
				"driver": runtime.NewString("chrome"),
			})

			var out EmbeddedOuterPointer
			err := sdk.Decode(obj, &out)

			So(err, ShouldBeNil)
			So(out, ShouldResemble, EmbeddedOuterPointer{
				EmbeddedParams: &EmbeddedParams{
					URL: "https://example.test",
				},
				Driver: "chrome",
			})
		})
	})

	Convey("Should not overwrite existing embedded pointer fields on partial update", t, func() {
		obj := runtime.NewObjectWith(map[string]runtime.Value{
			"url": runtime.NewString("new-url"),
		})

		out := EmbeddedOuterPointer{
			EmbeddedParams: &EmbeddedParams{
				URL:         "old",
				UserAgent:   "ua",
				KeepCookies: true,
				Charset:     "utf-8",
			},
			Driver: "old-driver",
		}

		err := sdk.Decode(obj, &out)

		So(err, ShouldBeNil)
		So(out, ShouldResemble, EmbeddedOuterPointer{
			EmbeddedParams: &EmbeddedParams{
				URL:         "new-url",
				UserAgent:   "ua",
				KeepCookies: true,
				Charset:     "utf-8",
			},
			Driver: "old-driver",
		})
	})

	Convey("Should avoid infinite recursion on self-embedded pointers", t, func() {
		obj := runtime.NewObject()

		var out EmbeddedNode
		err := sdk.Decode(obj, &out)

		So(err, ShouldBeNil)
		So(out.EmbeddedNode, ShouldBeNil)
	})

	Convey("Should close iterators when decoding slices", t, func() {
		closed := false
		source := closableIterable{
			values: []runtime.Value{
				runtime.NewString("a"),
				runtime.NewString("b"),
			},
			closed: &closed,
		}

		var out []string
		err := sdk.Decode(source, &out)

		So(err, ShouldBeNil)
		So(out, ShouldResemble, []string{"a", "b"})
		So(closed, ShouldBeTrue)
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
