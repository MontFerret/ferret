package encoding

import (
	"context"
	"errors"
	"testing"

	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestNewRegistryHasJSONCodec(t *testing.T) {
	registry := NewRegistry()

	codec, err := registry.Codec(ContentTypeJSON)
	if err != nil {
		t.Fatalf("expected json codec, got error: %v", err)
	}

	input := runtime.NewObjectWith(map[string]runtime.Value{
		"foo": runtime.NewString("bar"),
	})

	data, err := codec.Encode(input)
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	decoded, err := codec.Decode(data)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}

	if runtime.CompareValues(decoded, input) != 0 {
		t.Fatalf("decoded value mismatch: %s", decoded.String())
	}
}

func TestRegistryContextResolvers(t *testing.T) {
	registry := NewRegistry()
	ctx := WithRegistry(context.Background(), registry)

	if _, err := CodecFrom(ctx, ContentTypeJSON); err != nil {
		t.Fatalf("codec from context failed: %v", err)
	}

	if _, err := EncoderFrom(ctx, ContentTypeJSON); err != nil {
		t.Fatalf("encoder from context failed: %v", err)
	}

	if _, err := DecoderFrom(ctx, ContentTypeJSON); err != nil {
		t.Fatalf("decoder from context failed: %v", err)
	}
}

func TestRegistryContextResolversError(t *testing.T) {
	if _, err := CodecFrom(context.Background(), ContentTypeJSON); !errors.Is(err, ErrRegistryNotFound) {
		t.Fatalf("expected ErrRegistryNotFound, got %v", err)
	}

	ctx := WithRegistry(context.Background(), NewEmptyRegistry())
	if _, err := CodecFrom(ctx, ContentTypeJSON); !errors.Is(err, ErrCodecNotFound) {
		t.Fatalf("expected ErrCodecNotFound, got %v", err)
	}
}

func TestRegistryNormalizesContentType(t *testing.T) {
	registry := NewEmptyRegistry()
	if err := registry.Register("Application/X-Test; Charset=UTF-8", encodingjson.Codec{}); err != nil {
		t.Fatalf("register failed: %v", err)
	}

	if _, err := registry.Codec("application/x-test"); err != nil {
		t.Fatalf("expected normalized lookup to work: %v", err)
	}
}
