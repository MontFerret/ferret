package sdk_test

import (
	"context"
	"errors"
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/sdk"
)

type strictDecodeOptions struct {
	Name string `ferret:"name"`
}

type nestedDecodeItem struct {
	Count int8 `ferret:"count"`
}

type nestedDecodeOptions struct {
	Items []nestedDecodeItem `ferret:"items"`
}

type decodeArgOptions struct {
	Name  string `ferret:"name"`
	Count int8   `ferret:"count"`
}

func TestEncodeReportsConversionFailures(t *testing.T) {
	t.Run("nested path", func(t *testing.T) {
		_, err := sdk.Encode(t.Context(), []any{1, make(chan int)})
		assertErrorContains(t, err, "$[1]")
		assertErrorContains(t, err, "cannot encode chan int")
	})

	t.Run("unsigned overflow", func(t *testing.T) {
		_, err := sdk.Encode(t.Context(), uint64(math.MaxUint64))
		assertErrorContains(t, err, "overflows Ferret Int")
	})

	t.Run("slice cycle", func(t *testing.T) {
		input := make([]any, 1)
		input[0] = input

		_, err := sdk.Encode(t.Context(), input)
		assertErrorContains(t, err, "$[0]")
		assertErrorContains(t, err, "cycle detected")
	})

	t.Run("map cycle", func(t *testing.T) {
		input := make(map[string]any)
		input["self"] = input

		_, err := sdk.Encode(t.Context(), input)
		assertErrorContains(t, err, `$["self"]`)
		assertErrorContains(t, err, "cycle detected")
	})

	t.Run("caller cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		_, err := sdk.Encode(ctx, []int{1, 2})
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected cancellation, got %v", err)
		}
	})
}

func TestDecodeOptionsAndPaths(t *testing.T) {
	source := runtime.NewObjectWith(map[string]runtime.Value{
		"NAME":  runtime.NewString("Ferret"),
		"extra": runtime.True,
	})

	t.Run("permissive and case insensitive by default", func(t *testing.T) {
		var output strictDecodeOptions
		if err := sdk.Decode(t.Context(), source, &output); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if output.Name != "Ferret" {
			t.Fatalf("got name %q", output.Name)
		}
	})

	t.Run("strict unknown fields", func(t *testing.T) {
		var output strictDecodeOptions
		err := sdk.Decode(t.Context(), source, &output, sdk.DisallowUnknownFields())
		assertErrorContains(t, err, `unknown field "extra"`)
	})

	t.Run("nested overflow path", func(t *testing.T) {
		nested := runtime.NewObjectWith(map[string]runtime.Value{
			"items": runtime.NewArrayWith(runtime.NewObjectWith(map[string]runtime.Value{
				"count": runtime.NewInt(128),
			})),
		})

		var output nestedDecodeOptions
		err := sdk.Decode(t.Context(), nested, &output)
		assertErrorContains(t, err, "$.items[0].count")
		assertErrorContains(t, err, "integer overflow")
	})

	t.Run("host value exact type", func(t *testing.T) {
		input := strictDecodeOptions{Name: "host"}
		var output strictDecodeOptions
		if err := sdk.Decode(t.Context(), sdk.NewHostValue(input), &output); err != nil {
			t.Fatalf("decode host value: %v", err)
		}
		if output != input {
			t.Fatalf("got %#v, want %#v", output, input)
		}
	})
}

func TestDecodeIteratorFailuresAndCancellation(t *testing.T) {
	t.Run("caller cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		var output []string
		err := sdk.Decode(ctx, &conversionIteratorValue{}, &output)
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected cancellation, got %v", err)
		}
	})

	t.Run("deadline is not exhaustion", func(t *testing.T) {
		source := &conversionIteratorValue{terminalErr: context.DeadlineExceeded}
		var output []string

		err := sdk.Decode(t.Context(), source, &output)
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("expected deadline error, got %v", err)
		}
		if !source.closed {
			t.Fatal("iterator was not closed")
		}
	})

	t.Run("close failure", func(t *testing.T) {
		closeErr := errors.New("close failed")
		source := &conversionIteratorValue{closeErr: closeErr}
		var output []string

		err := sdk.Decode(t.Context(), source, &output)
		if !errors.Is(err, closeErr) {
			t.Fatalf("expected close error, got %v", err)
		}
		assertErrorContains(t, err, "$: close iterator")
	})

	t.Run("collection length overflow", func(t *testing.T) {
		if strconv.IntSize != 32 {
			t.Skip("native int can represent every runtime.Int value")
		}

		source := &conversionIteratorValue{length: runtime.Int(math.MaxInt64)}
		var output []string
		err := sdk.Decode(t.Context(), source, &output)
		assertErrorContains(t, err, "collection length overflows int")
	})
}

func TestDecodeArguments(t *testing.T) {
	args := []runtime.Value{runtime.NewObjectWith(map[string]runtime.Value{
		"name":  runtime.NewString("Ferret"),
		"count": runtime.NewInt(2),
	})}

	decoded, err := sdk.DecodeArg[decodeArgOptions](t.Context(), args, 0, sdk.DisallowUnknownFields())
	if err != nil {
		t.Fatalf("decode argument: %v", err)
	}
	if decoded != (decodeArgOptions{Name: "Ferret", Count: 2}) {
		t.Fatalf("unexpected argument: %#v", decoded)
	}

	_, err = sdk.DecodeArg[decodeArgOptions](t.Context(), []runtime.Value{runtime.NewString("wrong")}, 0)
	position, ok, _ := runtime.InvalidArgumentDetails(err)
	if !ok || position != 0 {
		t.Fatalf("expected position 0, got position=%d ok=%v err=%v", position, ok, err)
	}

	fallback := decodeArgOptions{Name: "fallback", Count: 7}
	missing, err := sdk.DecodeArgOr(t.Context(), nil, 1, fallback)
	if err != nil || missing != fallback {
		t.Fatalf("expected fallback, got %#v, %v", missing, err)
	}

	invalid := runtime.NewObjectWith(map[string]runtime.Value{
		"name":  runtime.NewString("changed"),
		"count": runtime.NewInt(128),
	})
	got, err := sdk.DecodeArgOr(t.Context(), []runtime.Value{invalid}, 0, fallback)
	if err == nil || got != fallback {
		t.Fatalf("expected unchanged fallback and error, got %#v, %v", got, err)
	}
}

func assertErrorContains(t *testing.T, err error, substring string) {
	t.Helper()

	if err == nil {
		t.Fatalf("expected error containing %q", substring)
	}
	if !strings.Contains(err.Error(), substring) {
		t.Fatalf("error %q does not contain %q", err, substring)
	}
}
