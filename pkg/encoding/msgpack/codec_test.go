package msgpack_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"math"
	"reflect"
	"strconv"
	"testing"
	"time"

	vmmsgpack "github.com/vmihailenco/msgpack/v5"

	ferretmsgpack "github.com/MontFerret/ferret/v2/pkg/encoding/msgpack"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type iterOnly struct {
	items []runtime.Value
}

type iterOnlyIter struct {
	items []runtime.Value
	idx   int
}

func (i *iterOnly) String() string {
	return "iterable"
}

func (i *iterOnly) Hash() uint64 {
	return 0
}

func (i *iterOnly) Copy() runtime.Value {
	return i
}

func (i *iterOnly) Iterate(_ context.Context) (runtime.Iterator, error) {
	return &iterOnlyIter{items: i.items}, nil
}

func (it *iterOnlyIter) Next(_ context.Context) (runtime.Value, runtime.Value, error) {
	if it.idx >= len(it.items) {
		return runtime.None, runtime.None, io.EOF
	}

	value := it.items[it.idx]
	key := runtime.NewInt(it.idx)
	it.idx++

	return value, key, nil
}

type unwrapValue struct {
	value any
}

func (u *unwrapValue) String() string {
	return "unwrap"
}

func (u *unwrapValue) Hash() uint64 {
	return 0
}

func (u *unwrapValue) Copy() runtime.Value {
	return u
}

func (u *unwrapValue) Unwrap() any {
	return u.value
}

type marshalerValue struct{}

func (m *marshalerValue) String() string {
	return "marshal"
}

func (m *marshalerValue) Hash() uint64 {
	return 0
}

func (m *marshalerValue) Copy() runtime.Value {
	return m
}

func (m *marshalerValue) MarshalMsgpack() ([]byte, error) {
	return vmmsgpack.Marshal(map[string]any{"custom": true})
}

type badValue struct {
	Fn func()
}

func (b *badValue) String() string {
	return "bad"
}

func (b *badValue) Hash() uint64 {
	return 0
}

func (b *badValue) Copy() runtime.Value {
	return b
}

func nestedArrayWithLeaf(depth int, leaf runtime.Value) runtime.Value {
	value := leaf

	for i := 0; i < depth; i++ {
		value = runtime.NewArrayWith(value)
	}

	return value
}

func nestedArray(depth int) runtime.Value {
	return nestedArrayWithLeaf(depth, runtime.NewInt(1))
}

func nestedObject(depth int) runtime.Value {
	value := runtime.Value(runtime.NewInt(1))

	for i := 0; i < depth; i++ {
		value = runtime.NewObjectWith(map[string]runtime.Value{"x": value})
	}

	return value
}

func hookValueLabel(value runtime.Value) string {
	switch v := value.(type) {
	case nil:
		return "<nil>"
	case *badValue:
		return "bad"
	case runtime.Map:
		return "map"
	case runtime.List:
		return "list"
	case runtime.String:
		return "string:" + string(v)
	default:
		return value.String()
	}
}

func mustMarshalNative(t *testing.T, value any) []byte {
	t.Helper()

	out, err := vmmsgpack.Marshal(value)
	if err != nil {
		t.Fatalf("msgpack marshal failed: %v", err)
	}

	return out
}

func mustEncodeRaw(t *testing.T, encode func(*vmmsgpack.Encoder) error) []byte {
	t.Helper()

	var buf bytes.Buffer
	enc := vmmsgpack.NewEncoder(&buf)

	if err := encode(enc); err != nil {
		t.Fatalf("raw encode failed: %v", err)
	}

	return buf.Bytes()
}

func assertValueEqual(t *testing.T, got, expected runtime.Value) {
	t.Helper()

	if runtime.CompareValues(got, expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, got)
	}
}

func TestMsgpackCodecEncode(t *testing.T) {
	codec := ferretmsgpack.Default
	ctx := context.Background()

	assertBytes := func(t *testing.T, value runtime.Value, expected []byte) {
		t.Helper()

		out, err := codec.Encode(value)
		if err != nil {
			t.Fatalf("encode failed: %v", err)
		}

		if !bytes.Equal(out, expected) {
			t.Fatalf("expected %x, got %x", expected, out)
		}
	}

	t.Run("none", func(t *testing.T) {
		assertBytes(t, runtime.None, mustMarshalNative(t, nil))
	})

	t.Run("boolean", func(t *testing.T) {
		assertBytes(t, runtime.True, mustMarshalNative(t, true))
		assertBytes(t, runtime.False, mustMarshalNative(t, false))
	})

	t.Run("string", func(t *testing.T) {
		assertBytes(t, runtime.NewString("hello"), mustMarshalNative(t, "hello"))
	})

	t.Run("int_matches_native", func(t *testing.T) {
		assertBytes(t, runtime.NewInt(42), mustMarshalNative(t, 42))
	})

	t.Run("float", func(t *testing.T) {
		assertBytes(t, runtime.NewFloat(42.5), mustMarshalNative(t, 42.5))
	})

	t.Run("binary", func(t *testing.T) {
		data := []byte("hi")
		assertBytes(t, runtime.NewBinary(data), mustMarshalNative(t, data))
	})

	t.Run("datetime", func(t *testing.T) {
		ts := time.Date(2024, time.January, 2, 3, 4, 5, 6, time.UTC)
		assertBytes(t, runtime.NewDateTime(ts), mustMarshalNative(t, ts))
	})

	t.Run("array", func(t *testing.T) {
		value := runtime.NewArrayWith(runtime.NewInt(1), runtime.NewString("a"), runtime.True)
		expected := mustMarshalNative(t, []any{1, "a", true})
		assertBytes(t, value, expected)
	})

	t.Run("deeply_nested_array", func(t *testing.T) {
		depth := 100_000
		value := nestedArray(depth)

		out, err := codec.Encode(value)
		if err != nil {
			t.Fatalf("encode failed: %v", err)
		}

		decoded, err := codec.Decode(out)
		if err != nil {
			t.Fatalf("decode failed: %v", err)
		}

		current := decoded
		for i := 0; i < depth; i++ {
			arr, ok := current.(*runtime.Array)
			if !ok {
				t.Fatalf("expected Array at depth %d, got %T", i, current)
			}

			current, err = arr.At(ctx, 0)
			if err != nil {
				t.Fatalf("array at depth %d failed: %v", i, err)
			}
		}

		assertValueEqual(t, current, runtime.NewInt(1))
	})

	t.Run("object_round_trip", func(t *testing.T) {
		obj := runtime.NewObject()
		if err := obj.Set(ctx, runtime.NewString("b"), runtime.NewInt(2)); err != nil {
			t.Fatalf("set b failed: %v", err)
		}
		if err := obj.Set(ctx, runtime.NewString("a"), runtime.NewInt(1)); err != nil {
			t.Fatalf("set a failed: %v", err)
		}

		out, err := codec.Encode(obj)
		if err != nil {
			t.Fatalf("encode failed: %v", err)
		}

		decoded, err := codec.Decode(out)
		if err != nil {
			t.Fatalf("decode failed: %v", err)
		}

		assertValueEqual(t, decoded, obj)
	})

	t.Run("deeply_nested_object", func(t *testing.T) {
		depth := 40_000
		value := nestedObject(depth)

		out, err := codec.Encode(value)
		if err != nil {
			t.Fatalf("encode failed: %v", err)
		}

		decoded, err := codec.Decode(out)
		if err != nil {
			t.Fatalf("decode failed: %v", err)
		}

		current := decoded
		for i := 0; i < depth; i++ {
			obj, ok := current.(*runtime.Object)
			if !ok {
				t.Fatalf("expected Object at depth %d, got %T", i, current)
			}

			current, err = obj.Get(ctx, runtime.NewString("x"))
			if err != nil {
				t.Fatalf("object get at depth %d failed: %v", i, err)
			}
		}

		assertValueEqual(t, current, runtime.NewInt(1))
	})

	t.Run("box", func(t *testing.T) {
		value := runtime.NewBox(7)
		out, err := codec.Encode(value)
		if err != nil {
			t.Fatalf("encode failed: %v", err)
		}

		decoded, err := codec.Decode(out)
		if err != nil {
			t.Fatalf("decode failed: %v", err)
		}

		assertValueEqual(t, decoded, runtime.NewInt(7))
	})

	t.Run("iterable_branch", func(t *testing.T) {
		iter := &iterOnly{items: []runtime.Value{runtime.NewInt(1), runtime.NewString("x")}}
		out, err := codec.Encode(iter)
		if err != nil {
			t.Fatalf("encode failed: %v", err)
		}

		decoded, err := codec.Decode(out)
		if err != nil {
			t.Fatalf("decode failed: %v", err)
		}

		assertValueEqual(t, decoded, runtime.NewArrayWith(runtime.NewInt(1), runtime.NewString("x")))
	})

	t.Run("unwrappable_branch", func(t *testing.T) {
		value := &unwrapValue{value: map[string]any{"a": 1}}
		out, err := codec.Encode(value)
		if err != nil {
			t.Fatalf("encode failed: %v", err)
		}

		decoded, err := codec.Decode(out)
		if err != nil {
			t.Fatalf("decode failed: %v", err)
		}

		expected := runtime.NewObjectWith(map[string]runtime.Value{
			"a": runtime.NewInt(1),
		})

		assertValueEqual(t, decoded, expected)
	})

	t.Run("msgpack_marshaler_branch", func(t *testing.T) {
		out, err := codec.Encode(&marshalerValue{})
		if err != nil {
			t.Fatalf("encode failed: %v", err)
		}

		decoded, err := codec.Decode(out)
		if err != nil {
			t.Fatalf("decode failed: %v", err)
		}

		expected := runtime.NewObjectWith(map[string]runtime.Value{
			"custom": runtime.True,
		})

		assertValueEqual(t, decoded, expected)
	})

	t.Run("object_msgpack_equivalence", func(t *testing.T) {
		obj := runtime.NewObjectWith(map[string]runtime.Value{
			"a": runtime.NewInt(1),
			"b": runtime.NewString("x"),
		})

		out, err := codec.Encode(obj)
		if err != nil {
			t.Fatalf("encode failed: %v", err)
		}

		var got map[string]any
		if err := vmmsgpack.Unmarshal(out, &got); err != nil {
			t.Fatalf("msgpack unmarshal failed: %v", err)
		}

		if got["b"] != "x" {
			t.Fatalf("expected string value x, got %#v", got["b"])
		}

		if got["a"] == nil || reflect.ValueOf(got["a"]).Kind() == reflect.Invalid ||
			got["a"] != int8(1) && got["a"] != int16(1) && got["a"] != int32(1) && got["a"] != int64(1) && got["a"] != int(1) &&
				got["a"] != uint8(1) && got["a"] != uint16(1) && got["a"] != uint32(1) && got["a"] != uint64(1) && got["a"] != uint(1) {
			t.Fatalf("expected integer value 1, got %#v", got["a"])
		}
	})
}

func TestMsgpackCodecEncodeHooks(t *testing.T) {
	t.Run("hooks_run_in_order", func(t *testing.T) {
		var order []string

		config := ferretmsgpack.Default.EncodeWith()
		config.PreHook(func(value runtime.Value) error {
			if value != runtime.NewInt(7) {
				t.Fatalf("expected hook value 7, got %v", value)
			}

			order = append(order, "pre1")
			return nil
		})
		config.PreHook(func(value runtime.Value) error {
			if value != runtime.NewInt(7) {
				t.Fatalf("expected hook value 7, got %v", value)
			}

			order = append(order, "pre2")
			return nil
		})
		config.PostHook(func(value runtime.Value, err error) error {
			if value != runtime.NewInt(7) {
				t.Fatalf("expected hook value 7, got %v", value)
			}

			if err != nil {
				t.Fatalf("expected successful encode, got %v", err)
			}

			order = append(order, "post")
			return nil
		})

		out, err := config.Encoder().Encode(runtime.NewInt(7))
		if err != nil {
			t.Fatalf("encode failed: %v", err)
		}

		decoded, err := ferretmsgpack.Default.Decode(out)
		if err != nil {
			t.Fatalf("decode failed: %v", err)
		}

		assertValueEqual(t, decoded, runtime.NewInt(7))

		expectedOrder := []string{"pre1", "pre2", "post"}
		if !reflect.DeepEqual(order, expectedOrder) {
			t.Fatalf("expected hook order %#v, got %#v", expectedOrder, order)
		}
	})

	t.Run("post_hooks_receive_encode_errors", func(t *testing.T) {
		var postErr error

		config := ferretmsgpack.Default.EncodeWith()
		config.PostHook(func(_ runtime.Value, err error) error {
			postErr = err
			return nil
		})

		_, err := config.Encoder().Encode(&badValue{Fn: func() {}})
		if err == nil {
			t.Fatal("expected encode error")
		}

		if postErr == nil {
			t.Fatal("expected post hook to receive encode error")
		}

		if postErr.Error() != err.Error() {
			t.Fatalf("expected post hook error %q, got %q", err, postErr)
		}
	})

	t.Run("nested_hooks_run_depth_first", func(t *testing.T) {
		var order []string

		config := ferretmsgpack.Default.EncodeWith()
		config.PreHook(func(value runtime.Value) error {
			order = append(order, "pre:"+hookValueLabel(value))
			return nil
		})
		config.PostHook(func(value runtime.Value, err error) error {
			if err != nil {
				t.Fatalf("expected successful encode, got %v", err)
			}

			order = append(order, "post:"+hookValueLabel(value))
			return nil
		})

		out, err := config.Encoder().Encode(runtime.NewObjectWith(map[string]runtime.Value{
			"items": runtime.NewArrayWith(runtime.NewString("leaf")),
		}))
		if err != nil {
			t.Fatalf("encode failed: %v", err)
		}

		decoded, err := ferretmsgpack.Default.Decode(out)
		if err != nil {
			t.Fatalf("decode failed: %v", err)
		}

		expectedValue := runtime.NewObjectWith(map[string]runtime.Value{
			"items": runtime.NewArrayWith(runtime.NewString("leaf")),
		})

		assertValueEqual(t, decoded, expectedValue)

		expectedOrder := []string{
			"pre:map",
			"pre:list",
			"pre:string:leaf",
			"post:string:leaf",
			"post:list",
			"post:map",
		}

		if !reflect.DeepEqual(order, expectedOrder) {
			t.Fatalf("expected hook order %#v, got %#v", expectedOrder, order)
		}
	})

	t.Run("ancestor_post_hooks_receive_descendant_errors", func(t *testing.T) {
		var postOrder []string
		var postErrs []string

		config := ferretmsgpack.Default.EncodeWith()
		config.PostHook(func(value runtime.Value, err error) error {
			postOrder = append(postOrder, hookValueLabel(value))

			if err == nil {
				postErrs = append(postErrs, "")
			} else {
				postErrs = append(postErrs, err.Error())
			}

			return nil
		})

		_, err := config.Encoder().Encode(runtime.NewObjectWith(map[string]runtime.Value{
			"items": runtime.NewArrayWith(&badValue{Fn: func() {}}),
		}))
		if err == nil {
			t.Fatal("expected encode error")
		}

		expectedOrder := []string{"bad", "list", "map"}
		if !reflect.DeepEqual(postOrder, expectedOrder) {
			t.Fatalf("expected post hook order %#v, got %#v", expectedOrder, postOrder)
		}

		if len(postErrs) != len(expectedOrder) {
			t.Fatalf("expected %d post hook errors, got %d", len(expectedOrder), len(postErrs))
		}

		for i, hookErr := range postErrs {
			if hookErr != err.Error() {
				t.Fatalf("expected hook %d to receive %q, got %q", i, err, hookErr)
			}
		}
	})

	t.Run("fallback_root_hooks_run_once", func(t *testing.T) {
		special := &iterOnly{items: []runtime.Value{runtime.NewString("leaf")}}
		value := nestedArrayWithLeaf(32_768, special)

		preCalls := 0
		postCalls := 0

		config := ferretmsgpack.Default.EncodeWith()
		config.PreHook(func(value runtime.Value) error {
			if value == special {
				preCalls++
			}

			return nil
		})
		config.PostHook(func(value runtime.Value, err error) error {
			if value == special {
				postCalls++
			}

			if err != nil {
				t.Fatalf("expected successful encode, got %v", err)
			}

			return nil
		})

		out, err := config.Encoder().Encode(value)
		if err != nil {
			t.Fatalf("encode failed: %v", err)
		}

		if len(out) == 0 {
			t.Fatal("expected encoded output")
		}

		if preCalls != 1 {
			t.Fatalf("expected fallback root pre-hook once, got %d", preCalls)
		}

		if postCalls != 1 {
			t.Fatalf("expected fallback root post-hook once, got %d", postCalls)
		}
	})

	t.Run("configured_encoder_does_not_mutate_default", func(t *testing.T) {
		calls := 0

		config := ferretmsgpack.Default.EncodeWith()
		config.PreHook(func(_ runtime.Value) error {
			calls++
			return nil
		})

		if _, err := config.Encoder().Encode(runtime.NewInt(1)); err != nil {
			t.Fatalf("configured encode failed: %v", err)
		}

		if calls != 1 {
			t.Fatalf("expected configured encoder to invoke hook once, got %d", calls)
		}

		if _, err := ferretmsgpack.Default.Encode(runtime.NewInt(2)); err != nil {
			t.Fatalf("default encode failed: %v", err)
		}

		if calls != 1 {
			t.Fatalf("expected default encoder to stay unmodified, got %d hook calls", calls)
		}
	})
}

func TestMsgpackCodecDecode(t *testing.T) {
	codec := ferretmsgpack.Default
	ctx := context.Background()

	t.Run("none", func(t *testing.T) {
		value, err := codec.Decode(mustMarshalNative(t, nil))
		if err != nil {
			t.Fatalf("decode failed: %v", err)
		}

		if value != runtime.None {
			t.Fatalf("expected None, got %T", value)
		}
	})

	t.Run("boolean", func(t *testing.T) {
		value, err := codec.Decode(mustMarshalNative(t, true))
		if err != nil {
			t.Fatalf("decode failed: %v", err)
		}

		if value != runtime.True {
			t.Fatalf("expected True, got %T", value)
		}
	})

	t.Run("string", func(t *testing.T) {
		value, err := codec.Decode(mustMarshalNative(t, "hello"))
		if err != nil {
			t.Fatalf("decode failed: %v", err)
		}

		if value != runtime.NewString("hello") {
			t.Fatalf("expected hello, got %v", value)
		}
	})

	t.Run("int", func(t *testing.T) {
		value, err := codec.Decode(mustMarshalNative(t, 1))
		if err != nil {
			t.Fatalf("decode failed: %v", err)
		}

		if _, ok := value.(runtime.Int); !ok {
			t.Fatalf("expected Int, got %T", value)
		}
	})

	t.Run("float", func(t *testing.T) {
		value, err := codec.Decode(mustMarshalNative(t, 1.5))
		if err != nil {
			t.Fatalf("decode failed: %v", err)
		}

		if _, ok := value.(runtime.Float); !ok {
			t.Fatalf("expected Float, got %T", value)
		}
	})

	t.Run("binary", func(t *testing.T) {
		data := []byte("hello")
		value, err := codec.Decode(mustMarshalNative(t, data))
		if err != nil {
			t.Fatalf("decode failed: %v", err)
		}

		if _, ok := value.(runtime.Binary); !ok {
			t.Fatalf("expected binary, got %T", value)
		}

		assertValueEqual(t, value, runtime.NewBinary(data))
	})

	t.Run("datetime", func(t *testing.T) {
		ts := time.Date(2024, time.January, 2, 3, 4, 5, 6, time.UTC)
		value, err := codec.Decode(mustMarshalNative(t, ts))
		if err != nil {
			t.Fatalf("decode failed: %v", err)
		}

		assertValueEqual(t, value, runtime.NewDateTime(ts))
	})

	t.Run("nested", func(t *testing.T) {
		value, err := codec.Decode(mustMarshalNative(t, map[string]any{
			"a": 1,
			"b": []any{true, nil, "x"},
		}))
		if err != nil {
			t.Fatalf("decode failed: %v", err)
		}

		obj, ok := value.(*runtime.Object)
		if !ok {
			t.Fatalf("expected Object, got %T", value)
		}

		a, err := obj.Get(ctx, runtime.NewString("a"))
		if err != nil {
			t.Fatalf("get a failed: %v", err)
		}

		if a != runtime.NewInt(1) {
			t.Fatalf("expected 1, got %v", a)
		}

		b, err := obj.Get(ctx, runtime.NewString("b"))
		if err != nil {
			t.Fatalf("get b failed: %v", err)
		}

		arr, ok := b.(*runtime.Array)
		if !ok {
			t.Fatalf("expected Array, got %T", b)
		}

		item, err := arr.At(ctx, 0)
		if err != nil {
			t.Fatalf("array at 0 failed: %v", err)
		}

		if item != runtime.True {
			t.Fatalf("expected true, got %v", item)
		}

		item, err = arr.At(ctx, 1)
		if err != nil {
			t.Fatalf("array at 1 failed: %v", err)
		}

		if item != runtime.None {
			t.Fatalf("expected none, got %v", item)
		}
	})

	t.Run("non_string_map_keys_normalize", func(t *testing.T) {
		data := mustEncodeRaw(t, func(enc *vmmsgpack.Encoder) error {
			if err := enc.EncodeMapLen(2); err != nil {
				return err
			}
			if err := enc.EncodeInt(1); err != nil {
				return err
			}
			if err := enc.EncodeString("one"); err != nil {
				return err
			}
			if err := enc.EncodeBool(true); err != nil {
				return err
			}
			return enc.EncodeString("yes")
		})

		value, err := codec.Decode(data)
		if err != nil {
			t.Fatalf("decode failed: %v", err)
		}

		obj, ok := value.(*runtime.Object)
		if !ok {
			t.Fatalf("expected Object, got %T", value)
		}

		one, err := obj.Get(ctx, runtime.NewString("1"))
		if err != nil {
			t.Fatalf("get 1 failed: %v", err)
		}
		assertValueEqual(t, one, runtime.NewString("one"))

		yes, err := obj.Get(ctx, runtime.NewString("true"))
		if err != nil {
			t.Fatalf("get true failed: %v", err)
		}
		assertValueEqual(t, yes, runtime.NewString("yes"))
	})

	t.Run("unsigned_overflow_error", func(t *testing.T) {
		data := mustEncodeRaw(t, func(enc *vmmsgpack.Encoder) error {
			return enc.EncodeUint64(math.MaxUint64)
		})

		_, err := codec.Decode(data)
		if err == nil {
			t.Fatal("expected overflow error")
		}
	})

	t.Run("int_overflow_float_on_32bit", func(t *testing.T) {
		if strconv.IntSize != 32 {
			t.Skip("only relevant on 32-bit platforms")
		}

		data := mustEncodeRaw(t, func(enc *vmmsgpack.Encoder) error {
			return enc.EncodeInt64(int64(math.MaxInt32) + 1)
		})

		value, err := codec.Decode(data)
		if err != nil {
			t.Fatalf("decode failed: %v", err)
		}

		if _, ok := value.(runtime.Float); !ok {
			t.Fatalf("expected Float, got %T", value)
		}
	})

	t.Run("empty_input_error", func(t *testing.T) {
		_, err := codec.Decode(nil)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("invalid_msgpack_error", func(t *testing.T) {
		_, err := codec.Decode([]byte{0xc1})
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("multiple_roots_error", func(t *testing.T) {
		data := append(mustMarshalNative(t, 1), mustMarshalNative(t, 2)...)

		_, err := codec.Decode(data)
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestMsgpackCodecDecodeHooks(t *testing.T) {
	t.Run("hooks_run_in_order", func(t *testing.T) {
		input := mustMarshalNative(t, map[string]any{"a": 1})
		var order []string

		config := ferretmsgpack.Default.DecodeWith()
		config.PreHook(func(data []byte) error {
			if !bytes.Equal(data, input) {
				t.Fatalf("expected input %x, got %x", input, data)
			}

			order = append(order, "pre1")
			return nil
		})
		config.PreHook(func(data []byte) error {
			if !bytes.Equal(data, input) {
				t.Fatalf("expected input %x, got %x", input, data)
			}

			order = append(order, "pre2")
			return nil
		})
		config.PostHook(func(data []byte, err error) error {
			if !bytes.Equal(data, input) {
				t.Fatalf("expected input %x, got %x", input, data)
			}

			if err != nil {
				t.Fatalf("expected successful decode, got %v", err)
			}

			order = append(order, "post")
			return nil
		})

		value, err := config.Decoder().Decode(input)
		if err != nil {
			t.Fatalf("decode failed: %v", err)
		}

		expectedValue := runtime.NewObjectWith(map[string]runtime.Value{
			"a": runtime.NewInt(1),
		})

		assertValueEqual(t, value, expectedValue)

		expectedOrder := []string{"pre1", "pre2", "post"}
		if !reflect.DeepEqual(order, expectedOrder) {
			t.Fatalf("expected hook order %#v, got %#v", expectedOrder, order)
		}
	})

	t.Run("post_hooks_receive_decode_errors", func(t *testing.T) {
		input := []byte{0xc1}
		var postErr error

		config := ferretmsgpack.Default.DecodeWith()
		config.PostHook(func(data []byte, err error) error {
			if !bytes.Equal(data, input) {
				t.Fatalf("expected input %x, got %x", input, data)
			}

			postErr = err
			return nil
		})

		_, err := config.Decoder().Decode(input)
		if err == nil {
			t.Fatal("expected decode error")
		}

		if postErr == nil {
			t.Fatal("expected post hook to receive decode error")
		}

		if postErr.Error() != err.Error() {
			t.Fatalf("expected post hook error %q, got %q", err, postErr)
		}
	})

	t.Run("configured_decoder_does_not_mutate_default", func(t *testing.T) {
		calls := 0

		config := ferretmsgpack.Default.DecodeWith()
		config.PreHook(func(_ []byte) error {
			calls++
			return nil
		})

		if _, err := config.Decoder().Decode(mustMarshalNative(t, 1)); err != nil {
			t.Fatalf("configured decode failed: %v", err)
		}

		if calls != 1 {
			t.Fatalf("expected configured decoder to invoke hook once, got %d", calls)
		}

		if _, err := ferretmsgpack.Default.Decode(mustMarshalNative(t, 1)); err != nil {
			t.Fatalf("default decode failed: %v", err)
		}

		if calls != 1 {
			t.Fatalf("expected default decoder to stay unmodified, got %d hook calls", calls)
		}
	})

	t.Run("pre_hook_error_short_circuits", func(t *testing.T) {
		expectedErr := errors.New("stop")

		config := ferretmsgpack.Default.DecodeWith()
		config.PreHook(func(_ []byte) error {
			return expectedErr
		})

		_, err := config.Decoder().Decode(mustMarshalNative(t, 1))
		if !errors.Is(err, expectedErr) {
			t.Fatalf("expected %v, got %v", expectedErr, err)
		}
	})
}
