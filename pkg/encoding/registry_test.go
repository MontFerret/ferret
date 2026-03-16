package encoding_test

import (
	"context"
	"errors"
	"testing"

	ferretencoding "github.com/MontFerret/ferret/v2/pkg/encoding"
	ferretjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type contentTypeCodec struct {
	contentType string
	ferretjson.Codec
}

func (c contentTypeCodec) ContentType() string {
	return c.contentType
}

func TestNewRegistryHasExplicitJSONCodec(t *testing.T) {
	registry := ferretencoding.NewRegistry(ferretjson.Default)

	codec, err := registry.Codec(ferretjson.ContentType)
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

func TestNewRegistryIsEmptyByDefault(t *testing.T) {
	registry := ferretencoding.NewRegistry()

	if _, err := registry.Codec(ferretjson.ContentType); !errors.Is(err, ferretencoding.ErrCodecNotFound) {
		t.Fatalf("expected ErrCodecNotFound, got %v", err)
	}
}

func TestRegistryContextResolvers(t *testing.T) {
	registry := ferretencoding.NewRegistry(ferretjson.Default)
	ctx := ferretencoding.WithRegistry(context.Background(), registry)

	if _, err := ferretencoding.CodecFrom(ctx, ferretjson.ContentType); err != nil {
		t.Fatalf("codec from context failed: %v", err)
	}

	if _, err := ferretencoding.EncoderFrom(ctx, ferretjson.ContentType); err != nil {
		t.Fatalf("encoder from context failed: %v", err)
	}

	if _, err := ferretencoding.DecoderFrom(ctx, ferretjson.ContentType); err != nil {
		t.Fatalf("decoder from context failed: %v", err)
	}
}

func TestRegistryContextResolversError(t *testing.T) {
	if _, err := ferretencoding.CodecFrom(context.Background(), ferretjson.ContentType); !errors.Is(err, ferretencoding.ErrRegistryNotFound) {
		t.Fatalf("expected ErrRegistryNotFound, got %v", err)
	}

	ctx := ferretencoding.WithRegistry(context.Background(), ferretencoding.NewEmptyRegistry())
	if _, err := ferretencoding.CodecFrom(ctx, ferretjson.ContentType); !errors.Is(err, ferretencoding.ErrCodecNotFound) {
		t.Fatalf("expected ErrCodecNotFound, got %v", err)
	}
}

func TestRegistryNormalizesContentType(t *testing.T) {
	registry := ferretencoding.NewEmptyRegistry()
	if err := registry.Register(contentTypeCodec{
		Codec:       ferretjson.Codec{},
		contentType: "Application/X-Test; Charset=UTF-8",
	}); err != nil {
		t.Fatalf("register failed: %v", err)
	}

	if _, err := registry.Codec("application/x-test"); err != nil {
		t.Fatalf("expected normalized lookup to work: %v", err)
	}
}
