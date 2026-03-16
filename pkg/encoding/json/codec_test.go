package json_test

import (
	"context"
	stdjson "encoding/json"
	"io"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/encoding/json"
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

func (m *marshalerValue) MarshalJSON() ([]byte, error) {
	return []byte(`{"custom":true}`), nil
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

func TestJSONCodecEncode(t *testing.T) {
	codec := json.Default

	assertJSON := func(t *testing.T, value runtime.Value, expected string) {
		t.Helper()

		out, err := codec.Encode(value)
		if err != nil {
			t.Fatalf("encode failed: %v", err)
		}

		if string(out) != expected {
			t.Fatalf("expected %s, got %s", expected, string(out))
		}
	}

	t.Run("none", func(t *testing.T) {
		assertJSON(t, runtime.None, "null")
	})

	t.Run("boolean", func(t *testing.T) {
		assertJSON(t, runtime.True, "true")
		assertJSON(t, runtime.False, "false")
	})

	t.Run("string_no_escape", func(t *testing.T) {
		assertJSON(t, runtime.NewString("<tag>"), "\"<tag>\"")
	})

	t.Run("int_matches_stdlib", func(t *testing.T) {
		expected, err := stdjson.Marshal(42)
		if err != nil {
			t.Fatalf("std json marshal failed: %v", err)
		}

		out, err := codec.Encode(runtime.NewInt(42))
		if err != nil {
			t.Fatalf("encode failed: %v", err)
		}

		if string(out) != string(expected) {
			t.Fatalf("expected %s, got %s", string(expected), string(out))
		}
	})

	t.Run("float", func(t *testing.T) {
		assertJSON(t, runtime.NewFloat(42.5), "42.5")
	})

	t.Run("binary", func(t *testing.T) {
		data := []byte("hi")
		expected, err := stdjson.Marshal(data)
		if err != nil {
			t.Fatalf("std json marshal failed: %v", err)
		}

		out, err := codec.Encode(runtime.NewBinary(data))
		if err != nil {
			t.Fatalf("encode failed: %v", err)
		}

		if string(out) != string(expected) {
			t.Fatalf("expected %s, got %s", string(expected), string(out))
		}
	})

	t.Run("datetime", func(t *testing.T) {
		ts := time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC)
		expected, err := stdjson.Marshal(ts)
		if err != nil {
			t.Fatalf("std json marshal failed: %v", err)
		}

		out, err := codec.Encode(runtime.NewDateTime(ts))
		if err != nil {
			t.Fatalf("encode failed: %v", err)
		}

		if string(out) != string(expected) {
			t.Fatalf("expected %s, got %s", string(expected), string(out))
		}
	})

	t.Run("array", func(t *testing.T) {
		arr := runtime.NewArrayWith(runtime.NewInt(1), runtime.NewString("a"), runtime.True)
		assertJSON(t, arr, "[1,\"a\",true]")
	})

	t.Run("object_canonical_order", func(t *testing.T) {
		obj := runtime.NewObjectWith(map[string]runtime.Value{
			"b": runtime.NewInt(2),
			"a": runtime.NewInt(1),
		})

		assertJSON(t, obj, "{\"a\":1,\"b\":2}")
	})

	t.Run("box", func(t *testing.T) {
		assertJSON(t, runtime.NewBox(7), "7")
	})

	t.Run("iterable_branch", func(t *testing.T) {
		iter := &iterOnly{items: []runtime.Value{runtime.NewInt(1), runtime.NewString("x")}}
		assertJSON(t, iter, "[1,\"x\"]")
	})

	t.Run("unwrappable_branch", func(t *testing.T) {
		value := &unwrapValue{value: map[string]any{"a": 1}}
		expected, err := stdjson.Marshal(map[string]any{"a": 1})
		if err != nil {
			t.Fatalf("std json marshal failed: %v", err)
		}

		out, err := codec.Encode(value)
		if err != nil {
			t.Fatalf("encode failed: %v", err)
		}

		if string(out) != string(expected) {
			t.Fatalf("expected %s, got %s", string(expected), string(out))
		}
	})

	t.Run("json_marshaler_branch", func(t *testing.T) {
		assertJSON(t, &marshalerValue{}, "{\"custom\":true}")
	})

	t.Run("object_json_equivalence", func(t *testing.T) {
		obj := runtime.NewObjectWith(map[string]runtime.Value{
			"a": runtime.NewInt(1),
			"b": runtime.NewString("x"),
		})

		out, err := codec.Encode(obj)
		if err != nil {
			t.Fatalf("encode failed: %v", err)
		}

		var got map[string]any
		if err := stdjson.Unmarshal(out, &got); err != nil {
			t.Fatalf("std json unmarshal failed: %v", err)
		}

		expected := map[string]any{
			"a": float64(1),
			"b": "x",
		}

		if !reflect.DeepEqual(got, expected) {
			t.Fatalf("expected %#v, got %#v", expected, got)
		}
	})
}

func TestJSONCodecEncodeHooks(t *testing.T) {
	t.Run("hooks_run_in_order", func(t *testing.T) {
		var order []string

		config := json.Default.EncodeWith()
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

		if string(out) != "7" {
			t.Fatalf("expected 7, got %s", out)
		}

		expectedOrder := []string{"pre1", "pre2", "post"}
		if !reflect.DeepEqual(order, expectedOrder) {
			t.Fatalf("expected hook order %#v, got %#v", expectedOrder, order)
		}
	})

	t.Run("post_hooks_receive_encode_errors", func(t *testing.T) {
		var postErr error

		config := json.Default.EncodeWith()
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

	t.Run("configured_encoder_does_not_mutate_default", func(t *testing.T) {
		calls := 0

		config := json.Default.EncodeWith()
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

		if _, err := json.Default.Encode(runtime.NewInt(2)); err != nil {
			t.Fatalf("default encode failed: %v", err)
		}

		if calls != 1 {
			t.Fatalf("expected default encoder to stay unmodified, got %d hook calls", calls)
		}
	})
}

func TestJSONCodecDecode(t *testing.T) {
	codec := json.Default
	ctx := context.Background()

	t.Run("none", func(t *testing.T) {
		value, err := codec.Decode([]byte("null"))
		if err != nil {
			t.Fatalf("decode failed: %v", err)
		}

		if value != runtime.None {
			t.Fatalf("expected None, got %T", value)
		}
	})

	t.Run("boolean", func(t *testing.T) {
		value, err := codec.Decode([]byte("true"))
		if err != nil {
			t.Fatalf("decode failed: %v", err)
		}

		if value != runtime.True {
			t.Fatalf("expected True, got %T", value)
		}
	})

	t.Run("string", func(t *testing.T) {
		value, err := codec.Decode([]byte(`"hello"`))
		if err != nil {
			t.Fatalf("decode failed: %v", err)
		}

		if value != runtime.NewString("hello") {
			t.Fatalf("expected hello, got %v", value)
		}
	})

	t.Run("int", func(t *testing.T) {
		value, err := codec.Decode([]byte("1"))
		if err != nil {
			t.Fatalf("decode failed: %v", err)
		}

		if _, ok := value.(runtime.Int); !ok {
			t.Fatalf("expected Int, got %T", value)
		}
	})

	t.Run("float", func(t *testing.T) {
		value, err := codec.Decode([]byte("1.5"))
		if err != nil {
			t.Fatalf("decode failed: %v", err)
		}

		if _, ok := value.(runtime.Float); !ok {
			t.Fatalf("expected Float, got %T", value)
		}
	})

	t.Run("exponent", func(t *testing.T) {
		value, err := codec.Decode([]byte("1e2"))
		if err != nil {
			t.Fatalf("decode failed: %v", err)
		}

		f, ok := value.(runtime.Float)
		if !ok {
			t.Fatalf("expected Float, got %T", value)
		}

		if f != runtime.NewFloat(100) {
			t.Fatalf("expected 100, got %v", f)
		}
	})

	t.Run("overflow_int_fallback", func(t *testing.T) {
		maxInt := int64(^uint(0) >> 1)
		overflow := new(big.Int).SetInt64(maxInt)
		overflow.Add(overflow, big.NewInt(1))

		value, err := codec.Decode([]byte(overflow.String()))
		if err != nil {
			t.Fatalf("decode failed: %v", err)
		}

		if _, ok := value.(runtime.Float); !ok {
			t.Fatalf("expected Float, got %T", value)
		}
	})

	t.Run("nested", func(t *testing.T) {
		value, err := codec.Decode([]byte(`{"a":1,"b":[true,null,"x"]}`))
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
			t.Fatalf("expected null, got %v", item)
		}
	})

	t.Run("empty_input_error", func(t *testing.T) {
		_, err := codec.Decode([]byte(""))
		if err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("invalid_json_error", func(t *testing.T) {
		_, err := codec.Decode([]byte("{"))
		if err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("multiple_roots_error", func(t *testing.T) {
		_, err := codec.Decode([]byte("1 2"))
		if err == nil {
			t.Fatalf("expected error")
		}
	})
}

func TestJSONCodecDecodeHooks(t *testing.T) {
	t.Run("hooks_run_in_order", func(t *testing.T) {
		input := []byte(`{"a":1}`)
		var order []string

		config := json.Default.DecodeWith()
		config.PreHook(func(data []byte) error {
			if string(data) != string(input) {
				t.Fatalf("expected input %s, got %s", input, data)
			}

			order = append(order, "pre1")
			return nil
		})
		config.PreHook(func(data []byte) error {
			if string(data) != string(input) {
				t.Fatalf("expected input %s, got %s", input, data)
			}

			order = append(order, "pre2")
			return nil
		})
		config.PostHook(func(data []byte, err error) error {
			if string(data) != string(input) {
				t.Fatalf("expected input %s, got %s", input, data)
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

		if runtime.CompareValues(value, expectedValue) != 0 {
			t.Fatalf("expected decoded value %s, got %s", expectedValue, value)
		}

		expectedOrder := []string{"pre1", "pre2", "post"}
		if !reflect.DeepEqual(order, expectedOrder) {
			t.Fatalf("expected hook order %#v, got %#v", expectedOrder, order)
		}
	})

	t.Run("post_hooks_receive_decode_errors", func(t *testing.T) {
		input := []byte("{")
		var postErr error

		config := json.Default.DecodeWith()
		config.PostHook(func(data []byte, err error) error {
			if string(data) != string(input) {
				t.Fatalf("expected input %s, got %s", input, data)
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

		config := json.Default.DecodeWith()
		config.PreHook(func(_ []byte) error {
			calls++
			return nil
		})

		if _, err := config.Decoder().Decode([]byte("1")); err != nil {
			t.Fatalf("configured decode failed: %v", err)
		}

		if calls != 1 {
			t.Fatalf("expected configured decoder to invoke hook once, got %d", calls)
		}

		if _, err := json.Default.Decode([]byte("1")); err != nil {
			t.Fatalf("default decode failed: %v", err)
		}

		if calls != 1 {
			t.Fatalf("expected default decoder to stay unmodified, got %d hook calls", calls)
		}
	})
}
