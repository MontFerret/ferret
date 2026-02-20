package vm_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type testQueryable struct {
	queries []runtime.Query
	result  runtime.Value
	err     error
}

func (t *testQueryable) Query(_ context.Context, q runtime.Query) (runtime.Value, error) {
	t.queries = append(t.queries, q)
	if t.err != nil {
		return runtime.None, t.err
	}
	if t.result != nil {
		return t.result, nil
	}
	return runtime.None, nil
}

func (t *testQueryable) MarshalJSON() ([]byte, error) {
	return json.Marshal("queryable")
}

func (t *testQueryable) String() string {
	return "queryable"
}

func (t *testQueryable) Unwrap() interface{} {
	return "queryable"
}

func (t *testQueryable) Hash() uint64 {
	return 0
}

func (t *testQueryable) Copy() runtime.Value {
	return t
}

func newObjectWithMap(props map[string]runtime.Value) runtime.Value {
	obj := runtime.NewObject()

	for key, value := range props {
		_ = obj.Set(context.Background(), runtime.NewString(key), value)
	}

	return obj
}

type mockText struct {
	value runtime.String
}

func newMockText(value string) *mockText {
	return &mockText{value: runtime.NewString(value)}
}

func (t *mockText) MarshalJSON() ([]byte, error) {
	return t.value.MarshalJSON()
}

func (t *mockText) String() string {
	return t.value.String()
}

func (t *mockText) Hash() uint64 {
	return t.value.Hash()
}

func (t *mockText) Copy() runtime.Value {
	return t
}

func (t *mockText) Iterate(ctx context.Context) (runtime.Iterator, error) {
	return runtime.NewArrayWith(t.value).Iterate(ctx)
}

type mockNode struct {
	kind string
}

func newMockNode(kind string) *mockNode {
	return &mockNode{kind: kind}
}

func (n *mockNode) MarshalJSON() ([]byte, error) {
	return runtime.NewString(n.kind).MarshalJSON()
}

func (n *mockNode) String() string {
	return n.kind
}

func (n *mockNode) Unwrap() interface{} {
	return n.kind
}

func (n *mockNode) Hash() uint64 {
	return runtime.NewString(n.kind).Hash()
}

func (n *mockNode) Copy() runtime.Value {
	return n
}

func (n *mockNode) Iterate(ctx context.Context) (runtime.Iterator, error) {
	return runtime.NewArrayWith(n).Iterate(ctx)
}

func (n *mockNode) Query(_ context.Context, q runtime.Query) (runtime.Value, error) {
	switch q.Kind.String() {
	case "css":
		switch q.Payload.String() {
		case ".product":
			return newMockNode("product"), nil
		case ".title":
			return newMockNode("title"), nil
		case ".price":
			return newMockNode("price"), nil
		default:
			return newMockNode("node"), nil
		}
	case "text":
		return newMockText(n.kind), nil
	default:
		return runtime.None, nil
	}
}

type mockDBQueryable struct {
	testQueryable
}

func (m *mockDBQueryable) Query(ctx context.Context, q runtime.Query) (runtime.Value, error) {
	m.queries = append(m.queries, q)

	if q.Kind.String() == "nil" {
		return runtime.None, nil
	}

	if q.Kind.String() != "sql" {
		return runtime.NewArray(0), nil
	}

	params, err := runtime.ToMap(ctx, q.Params)
	convey.So(err, convey.ShouldBeNil)
	category, _ := params.Get(ctx, runtime.NewString("c"))
	if category == runtime.NewString("laptops") {
		return runtime.NewArrayWith(
			newObjectWithMap(map[string]runtime.Value{
				"name":  runtime.NewString("Laptop Pro"),
				"price": runtime.NewInt(200),
			}),
		), nil
	}

	return runtime.NewArrayWith(
		newObjectWithMap(map[string]runtime.Value{
			"name":  runtime.NewString("Laptop Pro"),
			"price": runtime.NewInt(200),
		}),
		newObjectWithMap(map[string]runtime.Value{
			"name":  runtime.NewString("Mouse"),
			"price": runtime.NewInt(50),
		}),
	), nil
}

type mockJSONQueryable struct {
	testQueryable
}

func (m *mockJSONQueryable) Query(ctx context.Context, q runtime.Query) (runtime.Value, error) {
	m.queries = append(m.queries, q)

	if q.Kind.String() != "jp" {
		return runtime.NewArray(0), nil
	}

	orders := runtime.NewArrayWith(
		newObjectWithMap(map[string]runtime.Value{
			"id":    runtime.NewInt(1),
			"total": runtime.NewInt(150),
			"items": runtime.NewArrayWith(
				newObjectWithMap(map[string]runtime.Value{"name": runtime.NewString("Item A")}),
				newObjectWithMap(map[string]runtime.Value{"name": runtime.NewString("Item B")}),
			),
		}),
		newObjectWithMap(map[string]runtime.Value{
			"id":    runtime.NewInt(2),
			"total": runtime.NewInt(80),
			"items": runtime.NewArrayWith(
				newObjectWithMap(map[string]runtime.Value{"name": runtime.NewString("Item C")}),
			),
		}),
	)

	return orders, nil
}

func TestQueryable(t *testing.T) {
	queryable := &testQueryable{result: runtime.NewString("ok")}

	RunUseCases(t, []UseCase{
		Case("RETURN @doc[~ css`.items`]", "ok", "Should apply query literal"),
		Case("RETURN @doc[~ sql`SELECT * FROM products`({ c: \"laptops\" })]", "ok", "Should apply query literal with params"),
		Case("RETURN @doc[~ text]", "ok", "Should apply query literal with no string payload"),
		RuntimeErrorCase("RETURN @val[~ css`x`]", ExpectedRuntimeError{Message: "Invalid type"}),
	}, vm.WithParams(map[string]runtime.Value{
		"doc": queryable,
		"val": runtime.NewInt(1),
	}))

	t.Run("Should receive correct queries", func(t *testing.T) {
		convey.Convey("Should be ok", t, func() {
			var hasCSS bool
			var hasSQLParams bool
			var hasText bool

			for _, q := range queryable.queries {
				switch q.Kind {
				case runtime.NewString("css"):
					if q.Payload == runtime.NewString(".items") {
						hasCSS = true
					}
				case runtime.NewString("text"):
					if q.Payload == runtime.EmptyString {
						hasText = true
					}
				case runtime.NewString("sql"):
					params, err := runtime.ToMap(context.Background(), q.Params)
					convey.So(err, convey.ShouldBeNil)

					value, err := params.Get(context.Background(), runtime.NewString("c"))
					if err == nil && value == runtime.NewString("laptops") {
						hasSQLParams = true
					}
				}
			}

			convey.SoMsg(fmt.Sprintf("Expected to receive a query with kind %q and payload %q", "css", ".items"), hasCSS, convey.ShouldBeTrue)
			convey.SoMsg(fmt.Sprintf("Expected to receive a query with kind %q and empty payload", "text"), hasText, convey.ShouldBeTrue)
			convey.SoMsg(fmt.Sprintf("Expected to receive a query with kind %q and params containing %q=%q", "sql", "c", "laptops"), hasSQLParams, convey.ShouldBeTrue)
		})
	})
}

func TestComplexQueries(t *testing.T) {
	queryableDoc := newMockNode("doc")
	queryableDB := &mockDBQueryable{}
	queryableJSON := &mockJSONQueryable{}

	RunUseCases(t, []UseCase{
		Case(
			"RETURN @doc\n    [~ css`.product`]\n    [~ css`.title`]\n    [~ text]",
			"title",
			"Should chain apply operators",
		),
		CaseArray(
			"RETURN\n  @db[~ sql`\n    SELECT name, price\n    FROM products\n    WHERE category = $c\n  `({ c: \"laptops\" })]",
			[]any{
				map[string]any{"name": "Laptop Pro", "price": 200},
			},
			"Should pass params to query",
		),
		CaseArray(
			"RETURN @doc\n    [~ css`.product`]\n    [~ css`.title`]\n    [~ text]\n    [* FILTER CURRENT != \"\" RETURN UPPER(CURRENT)]",
			[]any{"TITLE"},
			"Should apply array operators after query chain",
		),
		CaseArray(
			"RETURN @db[~ sql`SELECT name, price FROM products`]\n    [* FILTER CURRENT.price > 100 RETURN CURRENT.name]",
			[]any{"Laptop Pro"},
			"Should filter results from query",
		),
		CaseArray(
			"RETURN @doc\n    [~ css`.product`]\n    [* FILTER CURRENT[~ css`.price`][~ text] != \"\"]\n    [~ css`.title`]\n    [~ text]",
			[]any{"title"},
			"Should combine query apply inside array filter",
		),
		CaseArray(
			"RETURN @json\n    [~ jp`$.orders[*]`]\n    [* FILTER CURRENT.total > 100]\n    [* RETURN {\n         id: CURRENT.id,\n         items: CURRENT.items[* RETURN CURRENT.name]\n       }]",
			[]any{
				map[string]any{
					"id":    1,
					"items": []any{"Item A", "Item B"},
				},
			},
			"Should project nested array operators",
		),
		CaseArray(
			"RETURN @doc\n    [~ css`.product`]\n    [* RETURN {\n         title: CURRENT[~ css`.title`][~ text],\n         price: CURRENT[~ css`.price`][~ text]\n       }]",
			[]any{
				map[string]any{"title": "title", "price": "price"},
			},
			"Should support nested apply inside projections",
		),
		CaseNil("RETURN @doc[~ nil`foo`]?.foo", "Should return null for queryable that returns None"),
		SkipCaseNil("RETURN @doc[~ nil`foo`]?.[*].name", "Should return null for queryable that returns None"),
	}, vm.WithParams(map[string]runtime.Value{
		"doc":  queryableDoc,
		"db":   queryableDB,
		"json": queryableJSON,
	}))
}
