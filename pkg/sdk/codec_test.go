package sdk_test

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/sdk"
)

type encoderFunc[T any] func(context.Context, T) (runtime.Value, error)

type codecContextKey struct{}

func (fn encoderFunc[T]) Encode(ctx context.Context, value T) (runtime.Value, error) {
	return fn(ctx, value)
}

func TestDefaultCodecRoundTrip(t *testing.T) {
	codec := sdk.DefaultCodec[int]()

	encoded, err := codec.Encode(t.Context(), 42)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if encoded != runtime.NewInt(42) {
		t.Fatalf("unexpected encoded value: %v", encoded)
	}

	decoded, err := codec.Decode(t.Context(), encoded)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if decoded != 42 {
		t.Fatalf("unexpected decoded value: %d", decoded)
	}
}

func TestNewCodecFallsBackPerCallback(t *testing.T) {
	tests := []struct {
		codec       sdk.Codec[int]
		wantEncoded runtime.Value
		name        string
		wantDecoded int
	}{
		{
			name:        "both defaults",
			codec:       sdk.NewCodec[int](nil, nil),
			wantEncoded: runtime.NewInt(3),
			wantDecoded: 7,
		},
		{
			name: "custom encoder",
			codec: sdk.NewCodec[int](func(context.Context, int) (runtime.Value, error) {
				return runtime.NewString("encoded"), nil
			}, nil),
			wantEncoded: runtime.NewString("encoded"),
			wantDecoded: 7,
		},
		{
			name: "custom decoder",
			codec: sdk.NewCodec[int](nil, func(context.Context, runtime.Value) (int, error) {
				return 11, nil
			}),
			wantEncoded: runtime.NewInt(3),
			wantDecoded: 11,
		},
		{
			name: "both custom",
			codec: sdk.NewCodec[int](
				func(context.Context, int) (runtime.Value, error) {
					return runtime.NewString("encoded"), nil
				},
				func(context.Context, runtime.Value) (int, error) {
					return 11, nil
				},
			),
			wantEncoded: runtime.NewString("encoded"),
			wantDecoded: 11,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encoded, err := test.codec.Encode(t.Context(), 3)
			if err != nil {
				t.Fatalf("encode: %v", err)
			}
			if encoded != test.wantEncoded {
				t.Fatalf("encoded value: got %v, want %v", encoded, test.wantEncoded)
			}

			decoded, err := test.codec.Decode(t.Context(), runtime.NewInt(7))
			if err != nil {
				t.Fatalf("decode: %v", err)
			}
			if decoded != test.wantDecoded {
				t.Fatalf("decoded value: got %d, want %d", decoded, test.wantDecoded)
			}
		})
	}
}

func TestCodecPropagatesContextAndCallbackErrors(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	codec := sdk.DefaultCodec[int]()
	if _, err := codec.Encode(ctx, 1); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected encode cancellation, got %v", err)
	}
	if _, err := codec.Decode(ctx, runtime.NewInt(1)); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected decode cancellation, got %v", err)
	}

	callbackErr := errors.New("codec callback failed")
	callbackCtx := context.WithValue(t.Context(), codecContextKey{}, "expected")
	encodeReceivedContext := false
	decodeReceivedContext := false
	codec = sdk.NewCodec[int](
		func(ctx context.Context, _ int) (runtime.Value, error) {
			encodeReceivedContext = ctx == callbackCtx
			return runtime.None, callbackErr
		},
		func(ctx context.Context, _ runtime.Value) (int, error) {
			decodeReceivedContext = ctx == callbackCtx
			return 0, callbackErr
		},
	)
	if _, err := codec.Encode(callbackCtx, 1); !errors.Is(err, callbackErr) {
		t.Fatalf("expected encode callback error, got %v", err)
	}
	if _, err := codec.Decode(callbackCtx, runtime.NewInt(1)); !errors.Is(err, callbackErr) {
		t.Fatalf("expected decode callback error, got %v", err)
	}
	if !encodeReceivedContext || !decodeReceivedContext {
		t.Fatalf("expected callbacks to receive the caller context: encode=%t decode=%t", encodeReceivedContext, decodeReceivedContext)
	}
}

func TestNativeIteratorsAcceptEncoderOnly(t *testing.T) {
	intEncoder := encoderFunc[int](func(_ context.Context, value int) (runtime.Value, error) {
		return runtime.NewInt(value), nil
	})
	stringEncoder := encoderFunc[string](func(_ context.Context, value string) (runtime.Value, error) {
		return runtime.NewString(value), nil
	})

	sliceIterator := sdk.NewSliceIteratorWithEncoding([]int{2}, intEncoder)
	value, key, err := sliceIterator.Next(t.Context())
	if err != nil || value != runtime.NewInt(2) || key != runtime.ZeroInt {
		t.Fatalf("slice iteration: value=%v key=%v err=%v", value, key, err)
	}

	mapIterator := sdk.NewMapIteratorWithEncoding(map[string]int{"two": 2}, stringEncoder, intEncoder)
	value, key, err = mapIterator.Next(t.Context())
	if err != nil || value != runtime.NewInt(2) || key != runtime.NewString("two") {
		t.Fatalf("map iteration: value=%v key=%v err=%v", value, key, err)
	}
}
