package session

import (
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	encodingmsgpack "github.com/MontFerret/ferret/v2/pkg/encoding/msgpack"
	"github.com/MontFerret/ferret/v2/pkg/internal/encodingownership"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestCodecDiscoversStoredMutationBeforeKeyWriteFailure(t *testing.T) {
	for _, codec := range []encoding.Codec{encodingjson.Default, encodingmsgpack.Default} {
		t.Run(codec.ContentType(), func(t *testing.T) {
			first := &abandonmentResource{Value: runtime.String("first")}
			added := &abandonmentResource{Value: runtime.String("added")}
			array := runtime.NewArrayWith(first)
			key := &abandonmentKey{Value: runtime.String("key"), mutate: func() {
				if err := array.Append(t.Context(), added); err != nil {
					t.Fatal(err)
				}
			}}
			value := &abandonmentMap{Object: runtime.NewObjectWith(map[string]runtime.Value{"key": array}), value: array, key: key}
			var result vm.Result
			adopter := outputValueAdopter{result: &result}
			failure := errors.New("key write failed")
			ctx := encodingownership.WithValueAdopter(t.Context(), adopter.Adopt)
			if err := codec.Encode(ctx, &abandonmentWriter{failure: failure}, value); !errors.Is(err, failure) {
				t.Fatalf("lost writer failure: %v", err)
			}
			if first.closes != 0 || added.closes != 0 {
				t.Fatal("codec closed children")
			}
			for range 2 {
				if err := result.Close(); err != nil {
					t.Fatal(err)
				}
			}
			if first.closes != 1 || added.closes != 1 {
				t.Fatalf("mutation cleanup: first=%d added=%d", first.closes, added.closes)
			}
		})
	}
}
