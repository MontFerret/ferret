package vm

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

const (
	testQueryDescriptorKeyKind    = "kind"
	testQueryDescriptorKeyPayload = "payload"
	testQueryDescriptorKeyOptions = "options"
	testQueryDescriptorKeyMod     = "modifier"
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

func programWithApplyQueryConstSource(src runtime.Value) *bytecode.Program {
	return programWithApplyQueryDescriptor(
		src,
		runtime.NewArrayWith(
			runtime.NewString("css"),
			runtime.NewString(".items"),
			runtime.None,
		),
	)
}

func programWithApplyQueryDescriptor(src runtime.Value, descriptor runtime.Value) *bytecode.Program {
	return &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  1,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpApplyQuery, bytecode.NewRegister(0), bytecode.NewConstant(0), bytecode.NewConstant(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		},
		Constants: []runtime.Value{
			src,
			descriptor,
		},
	}
}

func assertSingleStringArrayResult(t *testing.T, result runtime.Value, expected runtime.String) {
	t.Helper()

	arr, err := runtime.CastArray(result)
	if err != nil {
		t.Fatalf("expected array result, got %T: %v", result, err)
	}

	length, err := arr.Length(context.Background())
	if err != nil {
		t.Fatalf("failed to read result length: %v", err)
	}

	if length != 1 {
		t.Fatalf("unexpected result length: got %d, want 1", length)
	}

	item, err := arr.At(context.Background(), runtime.NewInt(0))
	if err != nil {
		t.Fatalf("failed to read result item: %v", err)
	}

	if item != expected {
		t.Fatalf("unexpected result item: got %v, want %v", item, expected)
	}
}

func TestApplyQuery_ConstantSourceStrict(t *testing.T) {
	stub := &queryStub{
		result: runtime.NewArrayWith(runtime.NewString("ok")),
	}

	instance := New(programWithApplyQueryConstSource(stub))
	result, err := instance.Run(context.Background(), nil)
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}

	assertSingleStringArrayResult(t, result, runtime.NewString("ok"))

	if len(stub.queries) != 1 {
		t.Fatalf("unexpected query count: got %d, want 1", len(stub.queries))
	}

	if got, want := stub.queries[0].Kind, runtime.NewString("css"); got != want {
		t.Fatalf("unexpected query kind: got %q, want %q", got, want)
	}

	if got, want := stub.queries[0].Payload, runtime.NewString(".items"); got != want {
		t.Fatalf("unexpected query payload: got %q, want %q", got, want)
	}
}

func TestApplyQuery_ConstantSourceFast(t *testing.T) {
	stub := &queryStub{
		result: runtime.NewArrayWith(runtime.NewString("ok")),
	}

	instance := NewWithOptions(programWithApplyQueryConstSource(stub), WithRunSafetyMode(RunSafetyFast))

	defer func() {
		if recovered := recover(); recovered != nil {
			t.Fatalf("unexpected panic in fast mode: %v", recovered)
		}
	}()

	result, err := instance.Run(context.Background(), nil)
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}

	assertSingleStringArrayResult(t, result, runtime.NewString("ok"))
}

func TestApplyQuery_ConstantSourceNonQueryableReturnsTypeError(t *testing.T) {
	program := programWithApplyQueryConstSource(runtime.NewInt(1))

	cases := []struct {
		name     string
		instance *VM
	}{
		{name: "strict", instance: New(program)},
		{name: "fast", instance: NewWithOptions(program, WithRunSafetyMode(RunSafetyFast))},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if recovered := recover(); recovered != nil {
					t.Fatalf("unexpected panic: %v", recovered)
				}
			}()

			_, err := tc.instance.Run(context.Background(), nil)
			if err == nil {
				t.Fatal("expected type error")
			}

			var rtErr *RuntimeError
			if !errors.As(err, &rtErr) {
				t.Fatalf("expected runtime error, got %T", err)
			}

			if !strings.Contains(strings.ToLower(rtErr.Message), "invalid type") {
				t.Fatalf("expected invalid type error, got %q", rtErr.Message)
			}
		})
	}
}

func TestApplyQuery_ArrayDescriptorWithModifier(t *testing.T) {
	stub := &queryStub{
		result: runtime.NewArrayWith(runtime.NewString("a"), runtime.NewString("b")),
	}

	program := programWithApplyQueryDescriptor(
		stub,
		runtime.NewArrayWith(
			runtime.NewString("css"),
			runtime.NewString(".items"),
			runtime.None,
			runtime.NewString("count"),
		),
	)

	_, err := New(program).Run(context.Background(), nil)
	if err == nil {
		t.Fatal("expected runtime error")
	}

	var rtErr *RuntimeError
	if !errors.As(err, &rtErr) {
		t.Fatalf("expected runtime error, got %T", err)
	}

	if !strings.Contains(strings.ToLower(rtErr.Format()), "unexpected query format") {
		t.Fatalf("expected unexpected query format error, got:\n%s", rtErr.Format())
	}
}

func TestApplyQuery_ObjectDescriptorWithModifier(t *testing.T) {
	stub := &queryStub{
		result: runtime.NewArrayWith(runtime.NewString("ok")),
	}

	obj := runtime.NewObject()
	ctx := context.Background()
	_ = obj.Set(ctx, runtime.NewString(testQueryDescriptorKeyKind), runtime.NewString("css"))
	_ = obj.Set(ctx, runtime.NewString(testQueryDescriptorKeyPayload), runtime.NewString(".items"))
	_ = obj.Set(ctx, runtime.NewString(testQueryDescriptorKeyOptions), runtime.None)
	_ = obj.Set(ctx, runtime.NewString(testQueryDescriptorKeyMod), runtime.NewString("EXISTS"))

	program := programWithApplyQueryDescriptor(stub, obj)
	result, err := New(program).Run(ctx, nil)
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}

	assertSingleStringArrayResult(t, result, runtime.NewString("ok"))

	if len(stub.queries) != 1 {
		t.Fatalf("unexpected query count: got %d, want 1", len(stub.queries))
	}

	if got, want := stub.queries[0].Kind, runtime.NewString("css"); got != want {
		t.Fatalf("unexpected query kind: got %q, want %q", got, want)
	}

	if got, want := stub.queries[0].Payload, runtime.NewString(".items"); got != want {
		t.Fatalf("unexpected query payload: got %q, want %q", got, want)
	}

	if got := stub.queries[0].Options; got != runtime.None {
		t.Fatalf("unexpected query options: got %v, want %v", got, runtime.None)
	}
}

func TestApplyQuery_ArrayDescriptorRequiresExactTupleSize(t *testing.T) {
	stub := &queryStub{
		result: runtime.NewArrayWith(runtime.NewString("a")),
	}

	program := programWithApplyQueryDescriptor(
		stub,
		runtime.NewArrayWith(
			runtime.NewString("css"),
			runtime.NewString(".items"),
		),
	)

	_, err := New(program).Run(context.Background(), nil)
	if err == nil {
		t.Fatal("expected runtime error")
	}

	var rtErr *RuntimeError
	if !errors.As(err, &rtErr) {
		t.Fatalf("expected runtime error, got %T", err)
	}

	if !strings.Contains(strings.ToLower(rtErr.Format()), "unexpected query format") {
		t.Fatalf("expected unexpected query format error, got:\n%s", rtErr.Format())
	}
}
