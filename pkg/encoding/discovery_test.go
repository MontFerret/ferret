package encoding_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	ferretjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	ferretmsgpack "github.com/MontFerret/ferret/v2/pkg/encoding/msgpack"
	"github.com/MontFerret/ferret/v2/pkg/internal/encodingownership"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestCodecAcquisitionAdoptsBeforeFailure(t *testing.T) {
	for _, codec := range []encoding.Codec{ferretjson.Default, ferretmsgpack.Default} {
		for _, shape := range []string{"iterator", "list", "map"} {
			for _, failure := range []string{"writer", "hook", "cancel"} {
				t.Run(codec.ContentType()+"/"+shape+"/"+failure, func(t *testing.T) {
					leaf := runtime.NewArrayWith(runtime.String("stored"))
					key := runtime.NewArrayWith(runtime.String("key"))
					ctx, cancel := context.WithCancel(t.Context())
					defer cancel()
					acquisitions := 0
					acquire := func() {
						acquisitions++
						if failure == "cancel" && (shape == "map" || acquisitions == 2) {
							cancel()
						}
					}
					var value runtime.Value
					switch shape {
					case "iterator":
						value = &codecIterable{iter: &codecIterator{items: []runtime.Value{runtime.Int(0), leaf}, onNext: acquire}}
					case "list":
						value = &discoveryList{Array: runtime.NewArrayWith(runtime.Int(0), leaf), acquire: acquire}
					case "map":
						value = &discoveryMap{Object: runtime.NewObjectWith(map[string]runtime.Value{"key": leaf}), value: leaf, key: key, acquire: acquire}
					}
					foundLeaf, foundKey := false, false
					lookups, callbacks := 0, 0
					ctx = discoveryContext{Context: encodingownership.WithValueAdopter(ctx, func(item runtime.Value, stored bool) {
						callbacks++
						if stored && item == leaf {
							foundLeaf = true
						}
						if stored && item == key {
							foundKey = true
						}
					}), lookups: &lookups}
					sentinel := errors.New("encoding stopped")
					enc := encoding.Encoder(codec)
					if failure == "hook" {
						enc = codec.EncodeWith().PreHook(func(_ context.Context, item runtime.Value) error {
							if item == leaf {
								return sentinel
							}
							return nil
						}).Encoder()
					}
					dst := io.Writer(io.Discard)
					if failure == "writer" {
						limit := 2
						if shape == "map" {
							limit = 1
						}
						dst = &failingWriter{limit: limit, failure: sentinel}
					}
					err := enc.Encode(ctx, dst, value)
					cause := sentinel
					if failure == "cancel" {
						cause = context.Canceled
					}
					if !errors.Is(err, cause) || !foundLeaf || (shape == "map" && !foundKey) || lookups != 1 {
						t.Fatalf("discovery before %s: leaf=%v key=%v lookups=%d err=%v", failure, foundLeaf, foundKey, lookups, err)
					}
					firstCallbacks := callbacks
					if err := codec.Encode(t.Context(), io.Discard, runtime.Int(1)); err != nil || callbacks != firstCallbacks {
						t.Fatalf("callback retained: calls=%d->%d err=%v", firstCallbacks, callbacks, err)
					}
				})
			}
		}
	}
}
func TestCodecRootDiscoverySurvivesCancellation(t *testing.T) {
	for _, codec := range []encoding.Codec{ferretjson.Default, ferretmsgpack.Default} {
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		value := runtime.NewArrayWith(runtime.String("stored"))
		found := false
		ctx = encodingownership.WithValueAdopter(ctx, func(item runtime.Value, stored bool) {
			if item == value && stored {
				found = true
			}
		})
		if err := codec.Encode(ctx, io.Discard, value); !errors.Is(err, context.Canceled) || !found {
			t.Fatalf("%s canceled root: found=%v err=%v", codec.ContentType(), found, err)
		}
	}
}
