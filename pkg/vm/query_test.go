package vm

import (
	"context"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type queryStub struct {
	result           runtime.List
	err              error
	oneResult        runtime.Value
	oneErr           error
	countResult      runtime.Int
	countResultSet   bool
	countErr         error
	existsResult     runtime.Boolean
	existsResultSet  bool
	existsErr        error
	queries          []runtime.Query
	queryCalls       int
	queryOneCalls    int
	queryCountCalls  int
	queryExistsCalls int
}

func (q *queryStub) Query(_ context.Context, query runtime.Query) (runtime.List, error) {
	q.queryCalls++
	q.queries = append(q.queries, query)

	if q.err != nil {
		return nil, q.err
	}

	if q.result != nil {
		return q.result, nil
	}

	return runtime.NewArray(0), nil
}

func (q *queryStub) QueryOne(ctx context.Context, query runtime.Query) (runtime.Value, error) {
	q.queryOneCalls++
	if q.oneErr != nil {
		return runtime.None, q.oneErr
	}
	if q.oneResult != nil {
		return q.oneResult, nil
	}

	return runtime.DefaultQueryOne(ctx, query, q.Query)
}

func (q *queryStub) QueryCount(ctx context.Context, query runtime.Query) (runtime.Int, error) {
	q.queryCountCalls++
	if q.countErr != nil {
		return runtime.ZeroInt, q.countErr
	}
	if q.countResultSet {
		return q.countResult, nil
	}

	return runtime.DefaultQueryCount(ctx, query, q.Query)
}

func (q *queryStub) QueryExists(ctx context.Context, query runtime.Query) (runtime.Boolean, error) {
	q.queryExistsCalls++
	if q.existsErr != nil {
		return runtime.False, q.existsErr
	}
	if q.existsResultSet {
		return q.existsResult, nil
	}

	return runtime.DefaultQueryExists(ctx, query, q.Query)
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

func TestApplyQueryExists_UsesQueryableModifier(t *testing.T) {
	src := &queryStub{
		existsResult:    runtime.True,
		existsResultSet: true,
	}

	out, err := applyQueryExists(context.Background(), src, validDescriptor())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out != runtime.True {
		t.Fatalf("expected true, got %v", out)
	}
	if src.queryExistsCalls != 1 || src.queryCalls != 0 {
		t.Fatalf("expected QueryExists only, got exists=%d query=%d", src.queryExistsCalls, src.queryCalls)
	}
}

func TestApplyQueryCount_UsesQueryableModifier(t *testing.T) {
	src := &queryStub{
		countResult:    runtime.NewInt(3),
		countResultSet: true,
	}

	out, err := applyQueryCount(context.Background(), src, validDescriptor())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out != runtime.NewInt(3) {
		t.Fatalf("expected count 3, got %v", out)
	}
	if src.queryCountCalls != 1 || src.queryCalls != 0 {
		t.Fatalf("expected QueryCount only, got count=%d query=%d", src.queryCountCalls, src.queryCalls)
	}
}

func TestApplyQueryOne_UsesQueryableModifier(t *testing.T) {
	src := &queryStub{
		oneResult: runtime.NewString("only"),
	}

	out, err := applyQueryOne(context.Background(), src, validDescriptor())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out != runtime.NewString("only") {
		t.Fatalf("expected only result, got %v", out)
	}
	if src.queryOneCalls != 1 || src.queryCalls != 0 {
		t.Fatalf("expected QueryOne only, got one=%d query=%d", src.queryOneCalls, src.queryCalls)
	}
}

func TestApplyQueryCount_ListSourceSumsCounts(t *testing.T) {
	a := &queryStub{countResult: runtime.NewInt(2), countResultSet: true}
	b := &queryStub{countResult: runtime.NewInt(3), countResultSet: true}

	out, err := applyQueryCount(context.Background(), runtime.NewArrayWith(a, b), validDescriptor())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out != runtime.NewInt(5) {
		t.Fatalf("expected count 5, got %v", out)
	}
	if a.queryCountCalls != 1 || b.queryCountCalls != 1 {
		t.Fatalf("expected both queryables to receive QueryCount, got a=%d b=%d", a.queryCountCalls, b.queryCountCalls)
	}
}

func TestApplyQueryExists_ListSourceShortCircuits(t *testing.T) {
	a := &queryStub{existsResult: runtime.False, existsResultSet: true}
	b := &queryStub{existsResult: runtime.True, existsResultSet: true}
	c := &queryStub{existsResult: runtime.True, existsResultSet: true}

	out, err := applyQueryExists(context.Background(), runtime.NewArrayWith(a, b, c), validDescriptor())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out != runtime.True {
		t.Fatalf("expected true, got %v", out)
	}
	if a.queryExistsCalls != 1 || b.queryExistsCalls != 1 || c.queryExistsCalls != 0 {
		t.Fatalf("expected short-circuit after second queryable, got a=%d b=%d c=%d", a.queryExistsCalls, b.queryExistsCalls, c.queryExistsCalls)
	}
}

func TestApplyQueryOne_ListSourceUsesCountThenOne(t *testing.T) {
	a := &queryStub{countResult: runtime.ZeroInt, countResultSet: true}
	b := &queryStub{countResult: runtime.NewInt(1), countResultSet: true, oneResult: runtime.NewString("only")}
	c := &queryStub{countResult: runtime.ZeroInt, countResultSet: true}

	out, err := applyQueryOne(context.Background(), runtime.NewArrayWith(a, b, c), validDescriptor())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out != runtime.NewString("only") {
		t.Fatalf("expected only result, got %v", out)
	}
	if a.queryCountCalls != 1 || b.queryCountCalls != 1 || c.queryCountCalls != 1 || b.queryOneCalls != 1 {
		t.Fatalf("expected count across all and one on matching queryable, got aCount=%d bCount=%d cCount=%d bOne=%d", a.queryCountCalls, b.queryCountCalls, c.queryCountCalls, b.queryOneCalls)
	}
}

func TestApplyQueryOne_ListSourceFailsForMultipleCombinedMatches(t *testing.T) {
	a := &queryStub{countResult: runtime.NewInt(1), countResultSet: true, oneResult: runtime.NewString("a")}
	b := &queryStub{countResult: runtime.NewInt(1), countResultSet: true, oneResult: runtime.NewString("b")}

	_, err := applyQueryOne(context.Background(), runtime.NewArrayWith(a, b), validDescriptor())
	if err == nil {
		t.Fatal("expected runtime error")
	}

	if !strings.Contains(err.Error(), runtime.QueryOneErrorMessage) {
		t.Fatalf("expected QUERY ONE cardinality error, got %v", err)
	}
	if a.queryOneCalls != 0 || b.queryOneCalls != 0 {
		t.Fatalf("did not expect QueryOne after combined count failure, got a=%d b=%d", a.queryOneCalls, b.queryOneCalls)
	}
}
