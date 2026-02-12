package vm_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type testQueryable struct {
	last   runtime.Query
	result runtime.Value
	err    error
}

func (t *testQueryable) ApplyQuery(_ context.Context, q runtime.Query) (runtime.Value, error) {
	t.last = q
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
		RuntimeErrorCase("RETURN @val[~ css`x`]", ExpectedRuntimeError{
			Message: "Invalid type",
		}),
	}, vm.WithParams(map[string]runtime.Value{
		"doc": queryable,
		"val": runtime.NewInt(1),
	}))

	if queryable.last.Kind != runtime.NewString("css") {
		t.Fatalf("expected query kind %q, got %q", "css", queryable.last.Kind.String())
	}

	if queryable.last.Payload != runtime.NewString(".items") {
		t.Fatalf("expected query payload %q, got %q", ".items", queryable.last.Payload.String())
	}
}
