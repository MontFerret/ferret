package ferret

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestEngineInjectsRegistryIntoContext(t *testing.T) {
	expected := encoding.NewRegistry()
	calledWithExpected := false

	eng, err := New(
		WithEncodingRegistry(expected),
		WithFunction("CHECK_REGISTRY", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			registry, err := encoding.RegistryFrom(ctx)
			if err != nil {
				return runtime.False, nil
			}

			if registry == expected {
				calledWithExpected = true
				return runtime.True, nil
			}

			return runtime.False, nil
		}),
	)
	if err != nil {
		t.Fatalf("new engine failed: %v", err)
	}

	result, err := eng.Run(context.Background(), file.NewAnonymousSource("RETURN CHECK_REGISTRY()"))
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	defer result.Close()

	value, err := One(context.Background(), result)
	if err != nil {
		t.Fatalf("read result failed: %v", err)
	}

	if value != runtime.True {
		t.Fatalf("expected true, got %s", value.String())
	}

	if !calledWithExpected {
		t.Fatalf("expected function to receive the engine registry")
	}
}

func TestEngineDefaultRegistryContainsJSONEncoder(t *testing.T) {
	eng, err := New(
		WithFunction("HAS_JSON_CODEC", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			_, err := encoding.EncoderFrom(ctx, encoding.ContentTypeJSON)
			if err != nil {
				return runtime.False, nil
			}

			return runtime.True, nil
		}),
	)
	if err != nil {
		t.Fatalf("new engine failed: %v", err)
	}

	result, err := eng.Run(context.Background(), file.NewAnonymousSource("RETURN HAS_JSON_CODEC()"))
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	defer result.Close()

	value, err := One(context.Background(), result)
	if err != nil {
		t.Fatalf("read result failed: %v", err)
	}

	if value != runtime.True {
		t.Fatalf("expected true, got %s", value.String())
	}
}
