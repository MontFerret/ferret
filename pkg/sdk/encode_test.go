package sdk_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/sdk"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type encodeParams struct {
	Name    string `ferret:"name"`
	City    string
	Ignored string `ferret:"-"`
	private string `ferret:"private"`
	Age     int    `json:"age"`
}

type encodeInner struct {
	Value  string `ferret:"value"`
	Hidden string
}

type encodeOuter struct {
	Inner *encodeInner `ferret:"inner"`
	Count int          `ferret:"count"`
}

type encodeNestedAddress struct {
	City string `ferret:"city"`
	Zip  int    `ferret:"zip"`
}

type encodeNestedProfile struct {
	Name    string              `ferret:"name"`
	Address encodeNestedAddress `ferret:"address"`
}

type encodeNestedFriendMeta struct {
	Active bool `ferret:"active"`
}

type encodeNestedFriend struct {
	Meta *encodeNestedFriendMeta `ferret:"meta"`
	Tags []string                `ferret:"tags"`
	ID   int                     `ferret:"id"`
}

type encodeNestedPayload struct {
	Profile encodeNestedProfile  `ferret:"profile"`
	Matrix  [][]int              `ferret:"matrix"`
	Friends []encodeNestedFriend `ferret:"friends"`
}

type EncodeEmbeddedParams struct {
	URL       string `json:"url"`
	UserAgent string `json:"userAgent"`
}

type EncodeEmbeddedPageLoadParams struct {
	EncodeEmbeddedParams
	Driver string `json:"driver"`
	URL    string `json:"url"`
}

type EncodeEmbeddedNode struct {
	*EncodeEmbeddedNode
}

type encodeTaggedNode struct {
	Next  *encodeTaggedNode `json:"next"`
	Value string            `json:"value"`
}

func TestEncode(t *testing.T) {
	Convey("Should encode tagged fields only", t, func() {
		input := encodeParams{
			Name:    "Alice",
			Age:     30,
			City:    "Paris",
			Ignored: "skip",
			private: "secret",
		}

		out := sdk.Encode(input)

		expected := runtime.NewObjectWith(
			map[string]runtime.Value{
				"name": runtime.NewString("Alice"),
				"age":  runtime.NewInt(30),
			},
		)

		So(out, ShouldResemble, expected)
	})

	Convey("Should encode deeply nested structs and slices", t, func() {
		input := encodeNestedPayload{
			Profile: encodeNestedProfile{
				Name: "Alice",
				Address: encodeNestedAddress{
					City: "Paris",
					Zip:  75001,
				},
			},
			Matrix: [][]int{{1, 2}, {3, 4}},
			Friends: []encodeNestedFriend{
				{
					ID:   1,
					Tags: []string{"a", "b"},
					Meta: &encodeNestedFriendMeta{Active: true},
				},
				{
					ID:   2,
					Tags: []string{"c"},
					Meta: &encodeNestedFriendMeta{Active: false},
				},
			},
		}

		out := sdk.Encode(input)

		expected := runtime.NewObjectWith(
			map[string]runtime.Value{
				"profile": runtime.NewObjectWith(
					map[string]runtime.Value{
						"name": runtime.NewString("Alice"),
						"address": runtime.NewObjectWith(
							map[string]runtime.Value{
								"city": runtime.NewString("Paris"),
								"zip":  runtime.NewInt(75001),
							},
						),
					},
				),
				"matrix": runtime.NewArrayWith(
					runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2)),
					runtime.NewArrayWith(runtime.NewInt(3), runtime.NewInt(4)),
				),
				"friends": runtime.NewArrayWith(
					runtime.NewObjectWith(
						map[string]runtime.Value{
							"id":   runtime.NewInt(1),
							"tags": runtime.NewArrayWith(runtime.NewString("a"), runtime.NewString("b")),
							"meta": runtime.NewObjectWith(
								map[string]runtime.Value{
									"active": runtime.NewBoolean(true),
								},
							),
						},
					),
					runtime.NewObjectWith(
						map[string]runtime.Value{
							"id":   runtime.NewInt(2),
							"tags": runtime.NewArrayWith(runtime.NewString("c")),
							"meta": runtime.NewObjectWith(
								map[string]runtime.Value{
									"active": runtime.NewBoolean(false),
								},
							),
						},
					),
				),
			},
		)

		So(out, ShouldResemble, expected)
	})

	Convey("Should encode anonymous embedded structs inline", t, func() {
		input := EncodeEmbeddedPageLoadParams{
			EncodeEmbeddedParams: EncodeEmbeddedParams{
				URL:       "https://example.test",
				UserAgent: "agent",
			},
			Driver: "chrome",
			URL:    "parent-url",
		}

		out := sdk.Encode(input)

		expected := runtime.NewObjectWith(
			map[string]runtime.Value{
				"driver":    runtime.NewString("chrome"),
				"url":       runtime.NewString("parent-url"),
				"userAgent": runtime.NewString("agent"),
			},
		)

		So(out, ShouldResemble, expected)
	})

	Convey("Should avoid infinite recursion on self-embedded pointers", t, func() {
		node := EncodeEmbeddedNode{}
		node.EncodeEmbeddedNode = &node

		out := sdk.Encode(node)

		So(out, ShouldResemble, runtime.NewObject())
	})

	Convey("Should encode tagged pointer cycles as none", t, func() {
		node := &encodeTaggedNode{
			Value: "root",
		}
		node.Next = node

		out := sdk.Encode(node)

		expected := runtime.NewObjectWith(map[string]runtime.Value{
			"value": runtime.NewString("root"),
			"next":  runtime.None,
		})

		So(out, ShouldResemble, expected)
	})

	Convey("Should encode tagged pointer chains without cycles", t, func() {
		tail := &encodeTaggedNode{
			Value: "tail",
		}
		head := &encodeTaggedNode{
			Value: "head",
			Next:  tail,
		}

		out := sdk.Encode(head)

		expected := runtime.NewObjectWith(map[string]runtime.Value{
			"value": runtime.NewString("head"),
			"next": runtime.NewObjectWith(map[string]runtime.Value{
				"value": runtime.NewString("tail"),
				"next":  runtime.None,
			}),
		})

		So(out, ShouldResemble, expected)
	})

	Convey("Should encode nested tagged structs", t, func() {
		input := encodeOuter{
			Inner: &encodeInner{
				Value:  "ok",
				Hidden: "skip",
			},
			Count: 2,
		}

		out := sdk.Encode(input)

		expected := runtime.NewObjectWith(
			map[string]runtime.Value{
				"inner": runtime.NewObjectWith(
					map[string]runtime.Value{
						"value": runtime.NewString("ok"),
					},
				),
				"count": runtime.NewInt(2),
			},
		)

		So(out, ShouldResemble, expected)
	})
}
