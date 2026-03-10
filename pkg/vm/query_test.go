package vm

import (
	"context"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type queryStub struct {
	queries []runtime.Query
	result  runtime.List
	err     error
}

func (q *queryStub) Query(_ context.Context, query runtime.Query) (runtime.List, error) {
	q.queries = append(q.queries, query)

	if q.err != nil {
		return nil, q.err
	}

	if q.result != nil {
		return q.result, nil
	}

	return runtime.NewArray(0), nil
}

func (q *queryStub) String() string {
	return "query-stub"
}

func (q *queryStub) Hash() uint64 {
	return runtime.NewString("query-stub").Hash()
}

func (q *queryStub) Copy() runtime.Value {
	return q
}

func validDescriptor() runtime.Value {
	return runtime.NewArrayWith(
		runtime.NewString("css"),
		runtime.NewString(".items"),
		runtime.None,
	)
}

func assertStringArray(t *testing.T, out runtime.Value, expected ...runtime.String) {
	t.Helper()

	arr, err := runtime.CastArray(out)
	if err != nil {
		t.Fatalf("expected array output, got %T: %v", out, err)
	}

	length, err := arr.Length(context.Background())
	if err != nil {
		t.Fatalf("failed to read result length: %v", err)
	}

	if int(length) != len(expected) {
		t.Fatalf("unexpected result length: got %d, want %d", length, len(expected))
	}

	for i, want := range expected {
		item, err := arr.At(context.Background(), runtime.NewInt(i))
		if err != nil {
			t.Fatalf("failed to read result item %d: %v", i, err)
		}

		if item != want {
			t.Fatalf("unexpected item %d: got %v, want %v", i, item, want)
		}
	}
}

func TestApplyQuery_ObjectDescriptorIgnoresModifier(t *testing.T) {
	src := &queryStub{
		result: runtime.NewArrayWith(runtime.NewString("ok")),
	}

	obj := runtime.NewObject()
	ctx := context.Background()
	_ = obj.Set(ctx, runtime.NewString("kind"), runtime.NewString("css"))
	_ = obj.Set(ctx, runtime.NewString("payload"), runtime.NewString(".items"))
	_ = obj.Set(ctx, runtime.NewString("options"), runtime.None)
	_ = obj.Set(ctx, runtime.NewString("modifier"), runtime.NewString("ONE"))

	out, err := applyQuery(ctx, src, obj)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertStringArray(t, out, runtime.NewString("ok"))

	if len(src.queries) != 1 {
		t.Fatalf("unexpected query count: got %d, want 1", len(src.queries))
	}
}

func TestApplyQuery_ArrayDescriptorRequiresExactTupleSize(t *testing.T) {
	src := &queryStub{
		result: runtime.NewArrayWith(runtime.NewString("ok")),
	}

	descriptor := runtime.NewArrayWith(
		runtime.NewString("css"),
		runtime.NewString(".items"),
		runtime.None,
		runtime.NewString("count"),
	)

	_, err := applyQuery(context.Background(), src, descriptor)
	if err == nil {
		t.Fatal("expected runtime error")
	}

	if !strings.Contains(strings.ToLower(err.Error()), "unexpected query format") {
		t.Fatalf("expected unexpected query format error, got %v", err)
	}
}

func TestApplyQuery_ArrayDescriptorPayloadTypeValidation(t *testing.T) {
	src := &queryStub{
		result: runtime.NewArrayWith(runtime.NewString("ok")),
	}

	descriptor := runtime.NewArrayWith(
		runtime.NewString("css"),
		runtime.NewInt(1),
		runtime.None,
	)

	_, err := applyQuery(context.Background(), src, descriptor)
	if err == nil {
		t.Fatal("expected runtime error")
	}

	if !strings.Contains(strings.ToLower(err.Error()), "invalid type") {
		t.Fatalf("expected invalid type error, got %v", err)
	}
}

func TestApplyQuery_QueryableNilResultNormalizedToEmptyArray(t *testing.T) {
	src := &queryStub{}

	out, err := applyQuery(context.Background(), src, validDescriptor())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	arr, err := runtime.CastArray(out)
	if err != nil {
		t.Fatalf("expected array output, got %T: %v", out, err)
	}

	length, err := arr.Length(context.Background())
	if err != nil {
		t.Fatalf("failed to read result length: %v", err)
	}

	if length != 0 {
		t.Fatalf("expected empty array result, got len=%d", length)
	}
}

func TestApplyQuery_ListSourceFlattensResults(t *testing.T) {
	a := &queryStub{
		result: runtime.NewArrayWith(runtime.NewString("a"), runtime.NewString("b")),
	}
	b := &queryStub{
		result: runtime.NewArrayWith(runtime.NewString("c")),
	}

	src := runtime.NewArrayWith(a, b)

	out, err := applyQuery(context.Background(), src, validDescriptor())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertStringArray(t, out, runtime.NewString("a"), runtime.NewString("b"), runtime.NewString("c"))

	if len(a.queries) != 1 || len(b.queries) != 1 {
		t.Fatalf("expected both queryables to be queried once, got a=%d b=%d", len(a.queries), len(b.queries))
	}
}

func TestApplyQuery_NonQueryableSourceTypeError(t *testing.T) {
	_, err := applyQuery(context.Background(), runtime.NewInt(1), validDescriptor())
	if err == nil {
		t.Fatal("expected type error")
	}

	if !strings.Contains(strings.ToLower(err.Error()), "invalid type") {
		t.Fatalf("expected invalid type error, got %v", err)
	}
}
