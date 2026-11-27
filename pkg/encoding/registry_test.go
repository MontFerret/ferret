package encoding_test

import (
	"context"
	"errors"
	"testing"

	ferretencoding "github.com/MontFerret/ferret/v2/pkg/encoding"
	ferretjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	ferretmsgpack "github.com/MontFerret/ferret/v2/pkg/encoding/msgpack"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type contentTypeCodec struct {
	ferretjson.Codec
	contentType string
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

	data, err := ferretencoding.EncodeBytes(t.Context(), codec, input)
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	decoded, err := ferretencoding.DecodeBytes(t.Context(), codec, data)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}

	equal, err := runtime.EqualValues(t.Context(), decoded, input)
	if err != nil || !equal {
		t.Fatalf("decoded value mismatch: %s", decoded.String())
	}
}

func TestNewRegistryHasExplicitMsgpackCodec(t *testing.T) {
	registry := ferretencoding.NewRegistry(ferretmsgpack.Default)

	codec, err := registry.Codec(ferretmsgpack.ContentType)
	if err != nil {
		t.Fatalf("expected msgpack codec, got error: %v", err)
	}

	input := runtime.NewObjectWith(map[string]runtime.Value{
		"foo": runtime.NewString("bar"),
	})

	data, err := ferretencoding.EncodeBytes(t.Context(), codec, input)
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	decoded, err := ferretencoding.DecodeBytes(t.Context(), codec, data)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}

	equal, err := runtime.EqualValues(t.Context(), decoded, input)
	if err != nil || !equal {
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

func TestRegistryContextBindings(t *testing.T) {
	callerRegistry := ferretencoding.NewRegistry(ferretmsgpack.Default)
	callerContext := ferretencoding.WithRegistry(nil, callerRegistry)
	jsonCodec := &contentTypeCodec{contentType: ferretjson.ContentType}
	finalized := ferretencoding.NewRegistry(jsonCodec)
	finalizedContext := ferretencoding.WithReadOnlyRegistry(callerContext, finalized)

	override := ferretencoding.WithRegistry(finalizedContext, callerRegistry)
	if registry, err := ferretencoding.RegistryFrom(override); err != nil || registry != callerRegistry {
		t.Fatalf("caller binding returned %p, %v, want %p", registry, err, callerRegistry)
	}
	if _, err := ferretencoding.CodecFrom(override, ferretmsgpack.ContentType); err != nil {
		t.Fatalf("caller binding did not shadow the finalized binding: %v", err)
	}
	if _, err := ferretencoding.CodecFrom(finalizedContext, ferretmsgpack.ContentType); !errors.Is(err, ferretencoding.ErrCodecNotFound) {
		t.Fatalf("finalized binding did not shadow the caller binding: %v", err)
	}

	snapshot, err := ferretencoding.RegistryFrom(finalizedContext)
	if err != nil || snapshot == finalized {
		t.Fatalf("expected a detached snapshot, got %p, %v", snapshot, err)
	}
	if got, err := snapshot.Codec(ferretjson.ContentType); err != nil || got != jsonCodec {
		t.Fatalf("snapshot codec = %T, %v, want the original instance", got, err)
	}
	if err := snapshot.Register(ferretmsgpack.Default); err != nil {
		t.Fatal(err)
	}
	another, err := ferretencoding.RegistryFrom(finalizedContext)
	if err != nil || another == snapshot {
		t.Fatalf("expected a fresh snapshot, got %p, %v", another, err)
	}
	if _, err := another.Codec(ferretmsgpack.ContentType); !errors.Is(err, ferretencoding.ErrCodecNotFound) {
		t.Fatalf("editing one snapshot affected a later snapshot: %v", err)
	}
}

func TestRegistryContextMissingBindings(t *testing.T) {
	parent := ferretencoding.WithRegistry(nil, ferretencoding.NewRegistry(ferretjson.Default))
	for _, tc := range []struct {
		ctx  context.Context
		name string
	}{
		{name: "nil context"},
		{name: "unbound context", ctx: t.Context()},
		{name: "nil registry shadows parent", ctx: ferretencoding.WithRegistry(parent, nil)},
		{name: "nil read-only registry shadows parent", ctx: ferretencoding.WithReadOnlyRegistry(parent, nil)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := ferretencoding.RegistryFrom(tc.ctx); !errors.Is(err, ferretencoding.ErrRegistryNotFound) {
				t.Errorf("registry error = %v, want ErrRegistryNotFound", err)
			}
			if _, err := ferretencoding.CodecFrom(tc.ctx, ferretjson.ContentType); !errors.Is(err, ferretencoding.ErrRegistryNotFound) {
				t.Errorf("codec error = %v, want ErrRegistryNotFound", err)
			}
			if _, err := ferretencoding.EncoderFrom(tc.ctx, ferretjson.ContentType); !errors.Is(err, ferretencoding.ErrRegistryNotFound) {
				t.Errorf("encoder error = %v, want ErrRegistryNotFound", err)
			}
			if _, err := ferretencoding.DecoderFrom(tc.ctx, ferretjson.ContentType); !errors.Is(err, ferretencoding.ErrRegistryNotFound) {
				t.Errorf("decoder error = %v, want ErrRegistryNotFound", err)
			}
		})
	}

	ctx := ferretencoding.WithReadOnlyRegistry(nil, ferretencoding.NewEmptyRegistry())
	if _, err := ferretencoding.CodecFrom(ctx, ""); !errors.Is(err, ferretencoding.ErrEmptyContentType) {
		t.Fatalf("invalid content type error = %v, want ErrEmptyContentType", err)
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
