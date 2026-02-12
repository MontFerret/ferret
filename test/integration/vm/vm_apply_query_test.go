package vm_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type testQueryable struct {
	queries []runtime.Query
	result  runtime.Value
	err     error
}

func (t *testQueryable) ApplyQuery(_ context.Context, q runtime.Query) (runtime.Value, error) {
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

func TestApplyQuery(t *testing.T) {
	queryable := &testQueryable{
		result: runtime.NewString("ok"),
	}

	RunUseCases(t, []UseCase{
		Case("RETURN @doc[~ css`.items`]", "ok", "Should apply query literal"),
		Case("RETURN @doc[~ sql`SELECT * FROM products`({ c: \"laptops\" })]", "ok", "Should apply query literal with params"),
		RuntimeErrorCase("RETURN @val[~ css`x`]", ExpectedRuntimeError{
			Message: "Invalid type",
		}),
	}, vm.WithParams(map[string]runtime.Value{
		"doc": queryable,
		"val": runtime.NewInt(1),
	}))

	var hasCSS bool
	var hasSQLParams bool

	for _, q := range queryable.queries {
		if q.Kind == runtime.NewString("css") && q.Payload == runtime.NewString(".items") {
			hasCSS = true
		}

		if q.Kind == runtime.NewString("sql") {
			params := runtime.ToMap(context.Background(), q.Params)
			value, err := params.Get(context.Background(), runtime.NewString("c"))
			if err == nil && value == runtime.NewString("laptops") {
				hasSQLParams = true
			}
		}
	}

	if !hasCSS {
		t.Fatalf("expected query kind %q with payload %q", "css", ".items")
	}
	if !hasSQLParams {
		t.Fatalf("expected query params to contain %q=%q", "c", "laptops")
	}
}
