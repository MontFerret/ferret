package sdk_test

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/sdk"
)

type controlledDecodeOptions struct {
	Raw   runtime.Value `ferret:"raw"`
	Name  string        `ferret:"name"`
	Count int8          `ferret:"count"`
}

type controlledNestedOptions struct {
	Name   string                  `ferret:"name"`
	Nested controlledDecodeOptions `ferret:"nested"`
}

func TestDecodeValue(t *testing.T) {
	source := runtime.NewObjectWith(map[string]runtime.Value{
		"name":  runtime.NewString("Ferret"),
		"count": runtime.NewInt(2),
	})

	decoded, err := sdk.DecodeValue[controlledDecodeOptions](
		t.Context(),
		source,
		sdk.RequireType(runtime.TypeMap),
		sdk.DisallowUnknownFields(),
	)
	if err != nil {
		t.Fatalf("decode value: %v", err)
	}
	if decoded.Name != "Ferret" || decoded.Count != 2 {
		t.Fatalf("unexpected value: %#v", decoded)
	}
}

func TestRequireType(t *testing.T) {
	t.Run("accepts any listed runtime type", func(t *testing.T) {
		decoded, err := sdk.DecodeValue[string](
			t.Context(),
			runtime.NewString("Ferret"),
			sdk.RequireType(runtime.TypeMap, runtime.TypeString),
		)
		if err != nil || decoded != "Ferret" {
			t.Fatalf("got %q, %v", decoded, err)
		}
	})

	t.Run("rejects mismatched root type", func(t *testing.T) {
		_, err := sdk.DecodeValue[string](
			t.Context(),
			runtime.NewString("Ferret"),
			sdk.RequireType(runtime.TypeMap),
		)

		assertDecodeError(t, err, "$", sdk.DecodeErrorKindType, true)
		if !errors.Is(err, runtime.ErrInvalidType) {
			t.Fatalf("expected invalid type, got %v", err)
		}
	})

	t.Run("requires at least one expected type", func(t *testing.T) {
		_, err := sdk.DecodeValue[string](
			t.Context(),
			runtime.NewString("Ferret"),
			sdk.RequireType(),
		)

		assertDecodeError(t, err, "$", sdk.DecodeErrorKindType, true)
		if !errors.Is(err, runtime.ErrInvalidArgument) {
			t.Fatalf("expected invalid argument, got %v", err)
		}
	})

	t.Run("rejects nil expected types", func(t *testing.T) {
		_, err := sdk.DecodeValue[string](
			t.Context(),
			runtime.NewString("Ferret"),
			sdk.RequireType(runtime.TypeString, nil),
		)

		assertDecodeError(t, err, "$", sdk.DecodeErrorKindType, true)
		if !errors.Is(err, runtime.ErrInvalidArgument) {
			t.Fatalf("expected invalid argument, got %v", err)
		}
	})
}

func TestOnlyFields(t *testing.T) {
	t.Run("matches root fields case insensitively", func(t *testing.T) {
		source := runtime.NewObjectWith(map[string]runtime.Value{
			"NAME": runtime.NewString("Ferret"),
		})

		decoded, err := sdk.DecodeValue[controlledDecodeOptions](
			t.Context(),
			source,
			sdk.OnlyFields("name"),
			sdk.DisallowUnknownFields(),
		)
		if err != nil || decoded.Name != "Ferret" {
			t.Fatalf("got %#v, %v", decoded, err)
		}
	})

	t.Run("rejects a tagged but disallowed root field", func(t *testing.T) {
		source := runtime.NewObjectWith(map[string]runtime.Value{
			"name":  runtime.NewString("Ferret"),
			"count": runtime.NewInt(2),
		})

		_, err := sdk.DecodeValue[controlledDecodeOptions](
			t.Context(),
			source,
			sdk.OnlyFields("name"),
		)

		assertDecodeError(t, err, "$", sdk.DecodeErrorKindUnknownField, true)
		if !errors.Is(err, runtime.ErrInvalidArgument) {
			t.Fatalf("expected invalid argument, got %v", err)
		}
	})

	t.Run("does not restrict nested fields", func(t *testing.T) {
		source := runtime.NewObjectWith(map[string]runtime.Value{
			"nested": runtime.NewObjectWith(map[string]runtime.Value{
				"name":  runtime.NewString("Ferret"),
				"count": runtime.NewInt(2),
			}),
		})

		decoded, err := sdk.DecodeValue[controlledNestedOptions](
			t.Context(),
			source,
			sdk.OnlyFields("nested"),
		)
		if err != nil {
			t.Fatalf("decode nested value: %v", err)
		}
		if decoded.Nested.Name != "Ferret" || decoded.Nested.Count != 2 {
			t.Fatalf("unexpected nested value: %#v", decoded)
		}
	})

	t.Run("combines with strict nested decoding", func(t *testing.T) {
		source := runtime.NewObjectWith(map[string]runtime.Value{
			"nested": runtime.NewObjectWith(map[string]runtime.Value{
				"name":  runtime.NewString("Ferret"),
				"extra": runtime.True,
			}),
		})

		_, err := sdk.DecodeValue[controlledNestedOptions](
			t.Context(),
			source,
			sdk.OnlyFields("nested"),
			sdk.DisallowUnknownFields(),
		)

		assertDecodeError(t, err, "$.nested", sdk.DecodeErrorKindUnknownField, true)
	})

	t.Run("requires a struct target", func(t *testing.T) {
		_, err := sdk.DecodeValue[map[string]string](
			t.Context(),
			runtime.NewObject(),
			sdk.OnlyFields("name"),
		)

		assertDecodeError(t, err, "$", sdk.DecodeErrorKindType, true)
		if !errors.Is(err, runtime.ErrInvalidArgument) {
			t.Fatalf("expected invalid argument, got %v", err)
		}
	})

	t.Run("rejects an empty configured field name as an authoring error", func(t *testing.T) {
		_, err := sdk.DecodeValue[controlledDecodeOptions](
			t.Context(),
			runtime.NewObject(),
			sdk.OnlyFields(""),
		)

		assertDecodeError(t, err, "$", sdk.DecodeErrorKindType, true)
		if !errors.Is(err, runtime.ErrInvalidArgument) {
			t.Fatalf("expected invalid argument, got %v", err)
		}
	})
}

func TestDisallowNoneValues(t *testing.T) {
	t.Run("takes precedence over a required root type", func(t *testing.T) {
		_, err := sdk.DecodeValue[controlledDecodeOptions](
			t.Context(),
			runtime.None,
			sdk.RequireType(runtime.TypeMap),
			sdk.DisallowNoneValues(),
		)

		assertDecodeError(t, err, "$", sdk.DecodeErrorKindNone, true)
		if !errors.Is(err, runtime.ErrInvalidArgument) {
			t.Fatalf("expected invalid argument, got %v", err)
		}
	})

	t.Run("rejects None for a native field", func(t *testing.T) {
		source := runtime.NewObjectWith(map[string]runtime.Value{
			"name": runtime.None,
		})

		_, err := sdk.DecodeValue[controlledDecodeOptions](
			t.Context(),
			source,
			sdk.DisallowNoneValues(),
		)

		assertDecodeError(t, err, "$.name", sdk.DecodeErrorKindNone, true)
		if !errors.Is(err, runtime.ErrInvalidArgument) {
			t.Fatalf("expected invalid argument, got %v", err)
		}
	})

	t.Run("preserves None for an exact runtime Value field", func(t *testing.T) {
		source := runtime.NewObjectWith(map[string]runtime.Value{
			"raw": runtime.None,
		})

		decoded, err := sdk.DecodeValue[controlledDecodeOptions](
			t.Context(),
			source,
			sdk.DisallowNoneValues(),
		)
		if err != nil {
			t.Fatalf("decode runtime value: %v", err)
		}
		if decoded.Raw != runtime.None {
			t.Fatalf("got %#v, want runtime.None", decoded.Raw)
		}
	})

	t.Run("preserves root None for an exact runtime Value", func(t *testing.T) {
		decoded, err := sdk.DecodeValue[runtime.Value](
			t.Context(),
			runtime.None,
			sdk.DisallowNoneValues(),
		)
		if err != nil {
			t.Fatalf("decode runtime value: %v", err)
		}
		if decoded != runtime.None {
			t.Fatalf("got %#v, want runtime.None", decoded)
		}
	})

	t.Run("preserves None in an exact runtime Value slice", func(t *testing.T) {
		decoded, err := sdk.DecodeValue[[]runtime.Value](
			t.Context(),
			runtime.NewArrayWith(runtime.None),
			sdk.DisallowNoneValues(),
		)
		if err != nil {
			t.Fatalf("decode runtime values: %v", err)
		}
		if len(decoded) != 1 || decoded[0] != runtime.None {
			t.Fatalf("got %#v, want [runtime.None]", decoded)
		}
	})

	t.Run("retains default None zeroing", func(t *testing.T) {
		decoded, err := sdk.DecodeValue[string](t.Context(), runtime.None)
		if err != nil || decoded != "" {
			t.Fatalf("got %q, %v", decoded, err)
		}
	})
}

func TestDecodeErrorMetadata(t *testing.T) {
	t.Run("range failure is safe and keeps its sentinel", func(t *testing.T) {
		source := runtime.NewObjectWith(map[string]runtime.Value{
			"count": runtime.NewInt(128),
		})

		_, err := sdk.DecodeValue[controlledDecodeOptions](t.Context(), source)
		assertDecodeError(t, err, "$.count", sdk.DecodeErrorKindRange, true)
		if !errors.Is(err, runtime.ErrInvalidArgumentType) {
			t.Fatalf("expected invalid argument type, got %v", err)
		}
	})

	t.Run("iterator failure is an unsafe source error", func(t *testing.T) {
		sentinel := runtime.Error(runtime.ErrInvalidArgument, "private source detail")
		source := &conversionIteratorValue{terminalErr: sentinel}

		_, err := sdk.DecodeValue[[]string](t.Context(), source)
		assertDecodeError(t, err, "$[0]", sdk.DecodeErrorKindSource, false)
		if !errors.Is(err, sentinel) {
			t.Fatalf("expected source sentinel, got %v", err)
		}
		if !errors.Is(err, runtime.ErrInvalidArgument) {
			t.Fatalf("expected wrapped runtime sentinel, got %v", err)
		}
	})

	t.Run("caller cancellation is an unsafe source error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		_, err := sdk.DecodeValue[string](ctx, runtime.NewString("Ferret"))
		assertDecodeError(t, err, "$", sdk.DecodeErrorKindSource, false)
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected cancellation, got %v", err)
		}
	})

	t.Run("DecodeArg preserves typed metadata", func(t *testing.T) {
		_, err := sdk.DecodeArg[controlledDecodeOptions](
			t.Context(),
			[]runtime.Value{runtime.NewString("wrong")},
			0,
			sdk.RequireType(runtime.TypeMap),
		)

		assertDecodeError(t, err, "$", sdk.DecodeErrorKindType, true)
		if !errors.Is(err, runtime.ErrInvalidArgument) || !errors.Is(err, runtime.ErrInvalidType) {
			t.Fatalf("expected argument and type sentinels, got %v", err)
		}
	})
}

func assertDecodeError(
	t *testing.T,
	err error,
	path string,
	kind sdk.DecodeErrorKind,
	safe bool,
) {
	t.Helper()

	var decodeErr *sdk.DecodeError
	if !errors.As(err, &decodeErr) {
		t.Fatalf("expected DecodeError, got %T: %v", err, err)
	}
	if decodeErr.Path() != path {
		t.Fatalf("path = %q, want %q", decodeErr.Path(), path)
	}
	if decodeErr.Kind() != kind {
		t.Fatalf("kind = %q, want %q", decodeErr.Kind(), kind)
	}
	if decodeErr.SafeToExpose() != safe {
		t.Fatalf("safe = %t, want %t", decodeErr.SafeToExpose(), safe)
	}
}
