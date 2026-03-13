package ferret

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type testIterableValue struct {
	values       []runtime.Value
	closeCalls   int
	iterateCalls int
	iterClosed   int
}

func (v *testIterableValue) String() string {
	return "testIterableValue"
}

func (v *testIterableValue) Hash() uint64 {
	return 0
}

func (v *testIterableValue) Copy() runtime.Value {
	return v
}

func (v *testIterableValue) Iterate(context.Context) (runtime.Iterator, error) {
	v.iterateCalls++

	return &testIterator{
		values:     v.values,
		closeCalls: &v.iterClosed,
	}, nil
}

func (v *testIterableValue) Close() error {
	v.closeCalls++
	return nil
}

type testIterator struct {
	closeCalls *int
	values     []runtime.Value
	index      int
}

func (it *testIterator) Next(context.Context) (runtime.Value, runtime.Value, error) {
	if it.index >= len(it.values) {
		return runtime.None, runtime.None, io.EOF
	}

	index := it.index
	value := it.values[index]
	it.index++

	return value, runtime.NewInt(index), nil
}

func (it *testIterator) Close() error {
	if it.closeCalls != nil {
		(*it.closeCalls)++
	}

	return nil
}

func assertIntResult(t *testing.T, got runtime.Value, want runtime.Int) {
	t.Helper()

	value, ok := got.(runtime.Int)
	if !ok {
		t.Fatalf("expected runtime.Int, got %T", got)
	}

	if value != want {
		t.Fatalf("unexpected value: got %v, want %v", value, want)
	}
}

func TestResultScalarNextAndValue(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	result := newResult(runtime.NewInt(42))
	if !result.IsScalar() {
		t.Fatal("expected scalar result")
	}

	val, err := result.Value()
	if err != nil {
		t.Fatalf("value failed: %v", err)
	}

	assertIntResult(t, val, runtime.NewInt(42))

	val, err = result.Next(ctx)
	if err != nil {
		t.Fatalf("next failed: %v", err)
	}

	assertIntResult(t, val, runtime.NewInt(42))

	_, err = result.Next(ctx)
	if !errors.Is(err, io.EOF) {
		t.Fatalf("expected EOF after scalar consumption, got: %v", err)
	}
}

func TestResultCollectScalar(t *testing.T) {
	t.Parallel()

	res := newResult(runtime.NewInt(7))
	values, err := res.Collect(context.Background())
	if err != nil {
		t.Fatalf("collect failed: %v", err)
	}

	if len(values) != 1 {
		t.Fatalf("expected one collected value, got %d", len(values))
	}

	assertIntResult(t, values[0], runtime.NewInt(7))
}

func TestResultIterableIterationOrder(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	result := newResult(runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2)))

	if result.IsScalar() {
		t.Fatal("expected iterable result")
	}

	first, err := result.First(ctx)
	if err != nil {
		t.Fatalf("first failed: %v", err)
	}

	assertIntResult(t, first, runtime.NewInt(1))

	second, err := result.Next(ctx)
	if err != nil {
		t.Fatalf("next failed: %v", err)
	}

	assertIntResult(t, second, runtime.NewInt(2))

	_, err = result.Next(ctx)
	if !errors.Is(err, io.EOF) {
		t.Fatalf("expected EOF after iterable consumption, got: %v", err)
	}
}

func TestResultValueReturnsErrNotScalar(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	result := newResult(runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2)))

	val, err := result.Value()
	if !errors.Is(err, ErrNotScalar) {
		t.Fatalf("expected ErrNotScalar, got: %v", err)
	}

	if val != runtime.None {
		t.Fatalf("expected runtime.None, got %T", val)
	}

	first, err := result.First(ctx)
	if err != nil {
		t.Fatalf("first failed after Value on iterable: %v", err)
	}

	assertIntResult(t, first, runtime.NewInt(1))
}

func TestResultMapIsScalar(t *testing.T) {
	t.Parallel()

	obj := runtime.NewObjectWith(map[string]runtime.Value{
		"a": runtime.NewInt(1),
	})

	result := newResult(obj)
	if !result.IsScalar() {
		t.Fatal("expected map result to be scalar")
	}

	val, err := result.Value()
	if err != nil {
		t.Fatalf("value failed: %v", err)
	}

	if val != obj {
		t.Fatalf("expected original object value, got %T", val)
	}
}

func TestResultFirstReturnsErrNoResult(t *testing.T) {
	t.Parallel()

	val, err := newResult(runtime.EmptyArray()).First(context.Background())
	if !errors.Is(err, ErrNoResult) {
		t.Fatalf("expected ErrNoResult, got: %v", err)
	}

	if val != runtime.None {
		t.Fatalf("expected runtime.None, got %T", val)
	}
}

func TestResultCollectIterable(t *testing.T) {
	t.Parallel()

	values, err := newResult(runtime.NewArrayWith(
		runtime.NewInt(1),
		runtime.NewInt(2),
		runtime.NewInt(3),
	)).Collect(context.Background())
	if err != nil {
		t.Fatalf("collect failed: %v", err)
	}

	if len(values) != 3 {
		t.Fatalf("expected three collected values, got %d", len(values))
	}

	assertIntResult(t, values[0], runtime.NewInt(1))
	assertIntResult(t, values[1], runtime.NewInt(2))
	assertIntResult(t, values[2], runtime.NewInt(3))
}

func TestResultCloseClosesIteratorAfterIterationStarts(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	value := &testIterableValue{
		values: []runtime.Value{
			runtime.NewInt(1),
			runtime.NewInt(2),
		},
	}

	result := newResult(value)
	if _, err := result.Next(ctx); err != nil {
		t.Fatalf("next failed: %v", err)
	}

	if err := result.Close(); err != nil {
		t.Fatalf("close failed: %v", err)
	}

	if value.iterClosed != 1 {
		t.Fatalf("expected iterator close once, got %d", value.iterClosed)
	}

	if value.closeCalls != 1 {
		t.Fatalf("expected underlying value close once, got %d", value.closeCalls)
	}

	if _, err := result.Next(ctx); !errors.Is(err, io.EOF) {
		t.Fatalf("expected EOF after close, got: %v", err)
	}
}

func TestResultCloseClosesUnderlyingValueBeforeIterationStarts(t *testing.T) {
	t.Parallel()

	value := &testIterableValue{
		values: []runtime.Value{
			runtime.NewInt(1),
		},
	}

	result := newResult(value)
	if err := result.Close(); err != nil {
		t.Fatalf("close failed: %v", err)
	}

	if value.closeCalls != 1 {
		t.Fatalf("expected underlying value close once, got %d", value.closeCalls)
	}

	if value.iterateCalls != 0 {
		t.Fatalf("expected close not to initialize iteration, got %d iterate calls", value.iterateCalls)
	}

	if value.iterClosed != 0 {
		t.Fatalf("expected iterator to remain unopened, got %d closes", value.iterClosed)
	}

	if _, err := result.Next(context.Background()); !errors.Is(err, io.EOF) {
		t.Fatalf("expected EOF after close, got: %v", err)
	}
}
