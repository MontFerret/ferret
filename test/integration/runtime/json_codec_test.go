package runtime_test

import (
	stdjson "encoding/json"
	"reflect"
	"testing"
	"time"

	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestJSONCodecBaseTypes(t *testing.T) {
	codec := encodingjson.Default

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

	t.Run("int", func(t *testing.T) {
		assertJSON(t, runtime.NewInt(42), "42")
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

	t.Run("object", func(t *testing.T) {
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
