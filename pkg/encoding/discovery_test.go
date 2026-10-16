package encoding_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"
	"time"

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
		t.Run(codec.ContentType(), func(t *testing.T) {
			for _, mode := range []string{"default", "configured"} {
				t.Run(mode, func(t *testing.T) {
					for _, cancellation := range []string{"canceled", "expired_deadline", "during_discovery"} {
						t.Run(cancellation, func(t *testing.T) {
							ctx, cancel := context.WithCancel(t.Context())
							defer cancel()

							want := error(context.Canceled)
							switch cancellation {
							case "canceled":
								cancel()
							case "expired_deadline":
								var expire context.CancelFunc
								ctx, expire = context.WithDeadline(ctx, time.Unix(0, 0))
								defer expire()

								want = context.DeadlineExceeded
							}

							value := runtime.NewArrayWith(runtime.String("stored"))
							found := false
							callbacks, pre, post := 0, 0, 0
							ctx = encodingownership.WithValueAdopter(ctx, func(item runtime.Value, stored bool) {
								callbacks++

								if item == value && stored {
									found = true

									if cancellation == "during_discovery" {
										cancel()
									}
								}
							})

							enc := encoding.Encoder(codec)
							if mode == "configured" {
								enc = codec.EncodeWith().PreHook(func(context.Context, runtime.Value) error {
									pre++

									return nil
								}).PostHook(func(context.Context, runtime.Value, error) error {
									post++

									return nil
								}).Encoder()
							}

							writer := &failingWriter{limit: 100}
							err := enc.Encode(ctx, writer, value)
							if !errors.Is(err, want) || !found {
								t.Errorf("root discovery: found=%v err=%v, want %v", found, err, want)
							}

							if pre != 0 || post != 0 || writer.calls != 0 || writer.Len() != 0 || writer.closed {
								t.Errorf("canceled encode started work: pre=%d post=%d writes=%d bytes=%d closed=%v", pre, post, writer.calls, writer.Len(), writer.closed)
							}

							previousCallbacks, previousWrites := callbacks, writer.calls
							previousBytes := bytes.Clone(writer.Bytes())
							nextWriter := &failingWriter{limit: 100}
							if err := enc.Encode(t.Context(), nextWriter, runtime.Int(1)); err != nil {
								t.Fatalf("subsequent encode failed: %v", err)
							}

							if callbacks != previousCallbacks || writer.calls != previousWrites || !bytes.Equal(writer.Bytes(), previousBytes) || writer.closed {
								t.Fatal("subsequent encode reused prior callback or writer")
							}

							if nextWriter.calls == 0 || nextWriter.Len() == 0 || nextWriter.closed {
								t.Fatal("subsequent encode did not use the fresh borrowed writer")
							}
						})
					}
				})
			}
		})
	}
}

func TestCodecCanceledEncodeDoesNotEvaluateHostValues(t *testing.T) {
	for _, codec := range []encoding.Codec{ferretjson.Default, ferretmsgpack.Default} {
		t.Run(codec.ContentType(), func(t *testing.T) {
			for _, shape := range []string{"list", "map", "iterator", "marshaler"} {
				t.Run(shape, func(t *testing.T) {
					ctx, cancel := context.WithCancel(t.Context())
					cancel()

					calls := 0
					observe := func() { calls++ }
					iter := &codecIterator{onNext: observe}
					var value runtime.Value

					switch shape {
					case "list":
						value = &discoveryList{Array: runtime.NewArrayWith(runtime.Int(1)), acquire: observe, length: observe}
					case "map":
						value = &discoveryMap{Object: runtime.NewObject(), acquire: observe, length: observe}
					case "iterator":
						value = &codecIterable{iter: iter, onIterate: observe}
					case "marshaler":
						value = &discoveryMarshaler{Value: runtime.String("custom"), marshal: observe}
					}

					found := false
					ctx = encodingownership.WithValueAdopter(ctx, func(item runtime.Value, stored bool) {
						if item == value && stored {
							found = true
						}
					})
					writer := &failingWriter{limit: 100}
					err := codec.Encode(ctx, writer, value)
					if !errors.Is(err, context.Canceled) || !found {
						t.Errorf("canceled host root: found=%v err=%v", found, err)
					}

					if calls != 0 || iter.closes != 0 || writer.calls != 0 || writer.closed {
						t.Fatalf("canceled encode evaluated host value: calls=%d iterator closes=%d writes=%d writer closed=%v", calls, iter.closes, writer.calls, writer.closed)
					}
				})
			}
		})
	}
}
